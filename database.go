package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func addSession(sessionHash string, userID int) {
	tx, err := db.Begin()
	checkErr(err)
	stmt, err := tx.Prepare("INSERT INTO SESSIONS (user_id, session_hash) VALUES (?, ?)")
	checkErr(err)
	defer stmt.Close()
	_, err = stmt.Exec(userID, sessionHash)
	checkErr(err)
	tx.Commit()
}

func isValidPassword(username string, password string) (int, bool) {
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
	var passwordHash string
	var userID int
	if rows.Next() {
		err = rows.Scan(&userID, &passwordHash)
		if err != nil {
			log.Print(err)
			return 0, false
		}
		err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
		if err == nil {
			return userID, true
		}
		log.Printf("Invalid password for: %s", username)
	}
	return 0, false
}

func getUserIDFromSessionHash(sessionHash string) (int, error) {
	// Then it checks if the session is in the SESSIONS table
	stmt, err := db.Prepare("SELECT user_id FROM SESSIONS WHERE session_hash=?")
	if err != nil {
		return -1, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(sessionHash)
	if err != nil {
		return -1, err
	}
	defer rows.Close()
	var userID int
	if rows.Next() {
		err = rows.Scan(&userID)
		if err != nil {
			return -1, err
		}
		return userID, nil
	}
	return -1, fmt.Errorf("Session not found")
}

func createDatabase() {
	log.Print("Creating new Database")
	_, err := db.Exec(`
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
		session_hash TEXT NOT NULL,
		create_date DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES USERS (user_id) 
	);`)
	checkErr(err)
}

func addUser(firstName string, lastName string, email string, username string, password string) {
	tx, err := db.Begin()
	stmt, err := tx.Prepare(`
	INSERT INTO USERS (
		first_name,
		last_name,
		email,
		username,
		password
	) VALUES (?, ?, ?, ?, ?);`)
	checkErr(err)
	defer stmt.Close()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	checkErr(err)
	_, err = stmt.Exec(firstName, lastName, email, username, string(hash))
	checkErr(err)
	tx.Commit()
}

func listUsers() {
	// SELECT type, name FROM sqlite_master where type="table"
	rows, err := db.Query("SELECT * FROM USERS")
	checkErr(err)
	defer rows.Close()
	log.Println("List of current users")
	var uid int
	var name string
	var lastName, email, username, password string
	for rows.Next() {
		err = rows.Scan(&uid, &name, &lastName, &email, &username, &password)
		checkErr(err)
		log.Printf("%s %s, %s\n", name, lastName, email)
	}
}

func initDatabase() {
	os.Remove("database.sqlite")
	var err error
	db, err = sql.Open("sqlite3", "./database.sqlite")
	checkErr(err)
	createDatabase()
	addUser("John", "Smith", "jonh.smith@example.com", "jsmith", "123")
	listUsers()
}
