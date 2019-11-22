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

type user struct {
	Username string
	Password string
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
	_, err = getUserIDFromSessionHash(sessionHash)
	if err != nil {
		log.Print(err)
		return false
	}
	return true
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

func rpcHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("RPC: %s", r.URL)
	if r.Method != "POST" {
		panic("Invalid method")
	}
	path := r.URL.Path[5:]
	w.Header().Add("Content-Type", "application/json")
	if path == "login/" {
		decoder := json.NewDecoder(r.Body)
		var t user
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
		panic("Invalid RPC")
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
	log.Printf("Serving Admin file: %s", path)
	appPath := fmt.Sprintf("%s", path)
	http.ServeFile(w, r, appPath)
	log.Printf("Serving Admin file: %s", appPath)
}

func main() {
	initDatabase()
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/rpc/", rpcHandler)
	http.HandleFunc("/admin/", adminHandler)
	log.Fatal(http.ListenAndServe(":1312", nil))
}
