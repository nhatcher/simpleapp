package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	_ "github.com/mattn/go-sqlite3"
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

func addSession(password string, userID int) {
	db, err := sql.Open("sqlite3", "./database.sqlite")
	checkErr(err)
	defer db.Close()
	tx, err := db.Begin()
	checkErr(err)
	stmt, err := tx.Prepare("INSERT INTO SESSIONS (user_id, password) VALUES (?, ?)")
	checkErr(err)
	defer stmt.Close()
	_, err = stmt.Exec(userID, password)
	checkErr(err)
	tx.Commit()
	checkErr(err)
}

func isValidUser(username string, password string) (int, bool) {
	db, err := sql.Open("sqlite3", "./database.sqlite")
	if err != nil {
		return 0, false
	}
	defer db.Close()
	stmt, err := db.Prepare("SELECT user_id, password FROM USERS where username=?")
	if err != nil {
		log.Print(err)
		return 0, false
	}
	defer stmt.Close()
	rows, err := stmt.Query(username)
	if err != nil {
		log.Print(err)
		return 0, false
	}
	defer rows.Close()
	var pword string
	var userID int
	if rows.Next() {
		err = rows.Scan(&userID, &pword)
		if err != nil {
			log.Print(err)
			return 0, false
		}
		if password == pword {
			return userID, true
		}
		log.Printf("Invalid password for: %s", username)
	}
	return 0, false
}

func isLoggedIn(r *http.Request) bool {
	sessionCookie, err := r.Cookie("session")
	if err != nil {
		log.Print(err)
		return false
	}
	sessionID, err := url.QueryUnescape(sessionCookie.Value)
	if err != nil {
		log.Print(err)
		return false
	}
	db, err := sql.Open("sqlite3", "./database.sqlite")
	if err != nil {
		log.Print(err)
		return false
	}
	defer db.Close()
	stmt, err := db.Prepare("SELECT user_id FROM SESSIONS WHERE password=?")
	if err != nil {
		log.Print(err)
		return false
	}
	defer stmt.Close()
	rows, err := stmt.Query(sessionID)
	if err != nil {
		log.Print(err)
		return false
	}
	defer rows.Close()
	var usrID int
	if rows.Next() {
		err = rows.Scan(&usrID)
		if err != nil {
			log.Print(err)
			return false
		}
		return true
	}
	log.Printf("Session not found")
	return false
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]
	log.Printf("Serving loging: %s", path)
	http.ServeFile(w, r, path)
}

func rpcHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("RPC: %s", r.URL)
	if r.Method != "POST" {
		panic("Invalid method")
	}
	path := r.URL.Path[5:]
	if path == "login/" {
		decoder := json.NewDecoder(r.Body)
		var t user
		err := decoder.Decode(&t)
		checkErr(err)
		userID, isValid := isValidUser(t.Username, t.Password)
		if isValid {
			sessionPassword := generateSessionPassword()
			addCookie(w, "session", sessionPassword, true)
			addCookie(w, "username", t.Username, false)
			addSession(sessionPassword, userID)
		}
		fmt.Fprintf(w, "%t", isValid)
	} else if path == "logout/" {
		// TODO: Remove Session
	} else {
		panic("Invalid RPC")
	}
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]
	log.Printf("Seriving file: %s", path)
	if !isLoggedIn(r) {
		http.ServeFile(w, r, "login/index.html")
	} else {
		appPath := fmt.Sprintf("app/%s", path)
		http.ServeFile(w, r, appPath)
	}
}

func testDB() {
	db, err := sql.Open("sqlite3", "./database.sqlite")
	checkErr(err)
	defer db.Close()
	// SELECT type, name FROM sqlite_master where type="table"
	rows, err := db.Query("SELECT * FROM USERS")
	if err != nil {
		log.Print("Creating new Database with a user")
		_, err = db.Exec(`
		CREATE TABLE USERS (
			user_id INTEGER PRIMARY KEY,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			username TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL
		);`)
		checkErr(err)
		_, err = db.Exec(`
		CREATE TABLE SESSIONS (
			session_id INTEGER PRIMARY KEY,
			user_id INTEGER,
			password TEXT NOT NULL,
			create_date DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES USERS (user_id) 
		);`)
		checkErr(err)
		_, err = db.Exec(`
		INSERT INTO USERS (
			first_name,
			last_name,
			email,
			username,
			password
		) VALUES (
			"John",
			"Smith",
			"john.smith@example.com",
			"nsmith",
			"abracadabra"
		);`)
		checkErr(err)
	} else {
		log.Print("List of current users")
		var uid int
		var name string
		var lastName, email, username, password string
		for rows.Next() {
			err = rows.Scan(&uid, &name, &lastName, &email, &username, &password)
			checkErr(err)
			log.Printf("%s %s, %s", name, lastName, email)
		}
		rows.Close()
	}
}

func main() {
	testDB()
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/login/", loginHandler)
	http.HandleFunc("/rpc/", rpcHandler)
	log.Fatal(http.ListenAndServe(":1312", nil))
}
