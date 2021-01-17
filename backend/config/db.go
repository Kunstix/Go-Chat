package config

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/kunstix/gochat/auth"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./chatdb.db")
	if err != nil {
		log.Fatal(err)
	}

	sqlStmt := `    
    CREATE TABLE IF NOT EXISTS room (
        id VARCHAR(255) NOT NULL PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        private TINYINT NULL
    );
    `
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal("%q: %s\n", err, sqlStmt)
	}

	sqlStmt = ` 
    CREATE TABLE IF NOT EXISTS user (
        id VARCHAR(255) NOT NULL PRIMARY KEY,
        name VARCHAR(255) NOT NULL UNIQUE,
        username VARCHAR(255) NULL UNIQUE,
        password VARCHARR(255) NULL
    );
    `
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal("%q: %s\n", err, sqlStmt)
	}

	password, _ := auth.GeneratePassword("pass1234")

	sqlStmt = `INSERT into user (id, name, username, password) VALUES
                    ('` + uuid.New().String() + `', 'kunstix', 'kunstix','` + password + `')`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal("%q: %s\n", err, sqlStmt)
	}

	return db
}
