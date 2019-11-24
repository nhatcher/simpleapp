package main

import (
	"database/sql"
	"fmt"
	"log"

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

func getUserIDFromSessionHash(sessionHash string) (int, int, error) {
	// Then it checks if the session is in the SESSIONS table
	stmt, err := db.Prepare("SELECT SESSIONS.user_id AS id, usertype_id FROM SESSIONS, USERS WHERE session_hash=? AND USERS.user_id=id;")
	if err != nil {
		return -1, -1, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(sessionHash)
	if err != nil {
		return -1, -1, err
	}
	defer rows.Close()
	var userID int
	var usertypeID int
	if rows.Next() {
		err = rows.Scan(&userID, &usertypeID)
		if err != nil {
			return -1, -1, err
		}
		return usertypeID, userID, nil
	}
	return -1, -1, fmt.Errorf("Session not found")
}

func addUser(firstName string, lastName string, email string, username string, password string, usertypeID int) {
	tx, err := db.Begin()
	stmt, err := tx.Prepare(`
	INSERT INTO USERS (
		first_name,
		last_name,
		email,
		username,
		password,
		usertype_id
	) VALUES (?, ?, ?, ?, ?, ?);`)
	checkErr(err)
	defer stmt.Close()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	checkErr(err)
	_, err = stmt.Exec(firstName, lastName, email, username, string(hash), usertypeID)
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
	var isRoot bool
	var name, lastName, email, username, password string
	for rows.Next() {
		err = rows.Scan(&uid, &name, &lastName, &email, &username, &password, &isRoot)
		checkErr(err)
		log.Printf("%s %s, %s\n", name, lastName, email)
	}
}

func initDatabase() {
	var err error
	db, err = sql.Open("sqlite3", "./database.sqlite")
	checkErr(err)
	addUser("John", "Smith", "jonh.smith@example.com", "jsmith", "123", 1)
	addUser("Penelope", "Glamour", "penelope.glamour@example.com", "gpenelope", "123", 2)
	listUsers()
}
