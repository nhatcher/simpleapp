package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

type userData struct {
	Username string
	Password string
}

type user struct {
	Name string `json:"name"`
	LastName string `json:"lastName"`
	Username string `json:"username"`
	Email string `json:"email"`
	UserID int `json:"userID"`
}

type register struct {
	Name string `json:"name"`
	LastName string `json:"lastName"`
	Email string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	UserType int `json:"userType"`
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func getSessionHash(r *http.Request) (string, error) {
	// First gets the session Cookie
	sessionCookie, err := r.Cookie("session")
	if err != nil {
		log.Print(err)
		return "", err
	}
	return url.QueryUnescape(sessionCookie.Value)
}

func isLoggedIn(r *http.Request) bool {
	sessionHash, err := getSessionHash(r)
	if err != nil {
		log.Print(err)
		return false
	}
	_, _, err = getUserIDFromSessionHash(sessionHash)
	if err != nil {
		log.Print(err)
		return false
	}
	return true
}

func isAdminLoggedIn(r *http.Request) bool {
	sessionHash, err := getSessionHash(r)
	if err != nil {
		log.Print(err)
		return false
	}
	userID, usertypeID, err := getUserIDFromSessionHash(sessionHash)
	if err != nil {
		log.Print(err)
		return false
	}
	if usertypeID == 2 {
		return true
	}
	log.Printf("Not root user: %d", userID)
	return false
}

func generateSessionPassword() string {
	b := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, b)
	checkErr(err)
	return base64.URLEncoding.EncodeToString(b)
}

func addCookie(w http.ResponseWriter, name string, value string, httpOnly bool) {
	expire := time.Now().AddDate(0, 0, 1)
	cookie := http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		Path:     "/",
		Expires:  expire,
		HttpOnly: httpOnly,
	}
	http.SetCookie(w, &cookie)
}

func adminRPCHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Admin RPC: %s", r.URL)
	if !isAdminLoggedIn(r) {
		fmt.Fprint(w, "{\"success\": false}")
		return
	}
	path := r.URL.Path[11:]
	w.Header().Add("Content-Type", "application/json")
	if path == "list-users" {
		userList := listUsers()
		users, err := json.Marshal(userList)
		checkErr(err)
		w.Write(users)
	} else if path == "delete-users" {
		if r.Method != "POST" {
			// panic("Invalid method")
			fmt.Fprint(w, "{\"success\": false}")
			return
		}
		body, err2 := json.Marshal(r.Body)
		checkErr(err2)
		// log.Printf("%v", r.Body)
		log.Printf("%v", body)
		decoder := json.NewDecoder(r.Body)
		var t user
		err := decoder.Decode(&t)
		checkErr(err)
		log.Printf("%v", t)
		deleteUser(t.UserID)
		log.Printf("%v", t.UserID)
	} else if path == "" {
		fmt.Fprint(w, "{\"success\": true}")
	} else if path == "add-user" {
		if r.Method != "POST" {
			// panic("Invalid method")
			fmt.Fprint(w, "{\"success\": false}")
			return
		}
		decoder := json.NewDecoder(r.Body)
		var t register
		err := decoder.Decode(&t)
		checkErr(err)
		addUser(t.Name, t.LastName, t.Email, t.Username, t.Password, t.UserType)
		fmt.Fprint(w, "{\"success\": true}")
	} else {
		fmt.Fprint(w, "{\"success\": false}")
		// panic("Invalid RPC")
	}
}

func rpcHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("RPC: %s", r.URL)
	if r.Method != "POST" {
	  // panic("Invalid method")
	  fmt.Fprint(w, "{\"success\": false}")
	  return
	}
	path := r.URL.Path[5:]
	w.Header().Add("Content-Type", "application/json")
	if path == "login/" {
	  decoder := json.NewDecoder(r.Body)
	  var t userData
	  err := decoder.Decode(&t)
	  checkErr(err)
	  userID, isValid := isValidPassword(t.Username, t.Password)
	  if isValid {
		sessionPassword := generateSessionPassword()
		addCookie(w, "session", sessionPassword, true)
		addCookie(w, "username", t.Username, false)
		addSession(sessionPassword, userID)
	  }
	  fmt.Fprintf(w, "{\"success\":%t}", isValid)
	} else if path == "logout/" {
	  // Remove Session and username cookies
	  addCookie(w, "session", "", true)
	  addCookie(w, "username", "", false)
	  fmt.Fprintf(w, "{\"success\":%t}", true)
	} else {
	  // panic("Invalid RPC")
	  fmt.Fprint(w, "{\"success\": false}")
	}
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]
	log.Printf("Serving file: %s", path)
	if !isLoggedIn(r) {
		appPath := fmt.Sprintf("login/%s", path)
		http.ServeFile(w, r, appPath)
	} else {
		appPath := fmt.Sprintf("app/%s", path)
		http.ServeFile(w, r, appPath)
	}
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]
	if isAdminLoggedIn(r) {
		log.Printf("Serving Admin file: %s", path)
		http.ServeFile(w, r, path)
	} else {
		appPath := fmt.Sprintf("login/%s", path)
		http.ServeFile(w, r, appPath)
	}
}

func main() {
	initDatabase()
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/rpc/", rpcHandler)
	http.HandleFunc("/admin/", adminHandler)
	http.HandleFunc("/admin/rpc/", adminRPCHandler)
	log.Fatal(http.ListenAndServe(":1312", nil))
}
