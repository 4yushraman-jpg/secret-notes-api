package database

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite", "backend.db")
	if err != nil {
		log.Fatal(err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	query := `CREATE TABLE IF NOT EXISTS users (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    email TEXT UNIQUE NOT NULL,
	    password TEXT NOT NULL
	);`

	_, err = DB.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	query = `CREATE TABLE IF NOT EXISTS notes (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    content TEXT NOT NULL,
	    user_id INTEGER NOT NULL,
	    FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	_, err = DB.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database initialized!")
}
