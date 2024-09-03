package data

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	conn *sql.DB
}

func New(dbPath string) *Database {
	var err error
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	migrate(db)

	return &Database{conn: db}
}

func migrate(db *sql.DB) {
	chatUserTable := `CREATE TABLE IF NOT EXISTS chat_users (
		id VARCHAR PRIMARY KEY,
		secret VARCHAR NOT NULL
	);`

	chatTable := `CREATE TABLE IF NOT EXISTS chats (
		id VARCHAR PRIMARY KEY,
		chat_user_id VARCHAR,
		role VARCHAR,
		text VARCHAR,
		audio VARCHAR,
		FOREIGN KEY(chat_user_id) REFERENCES chat_users(id)
	);`

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(chatUserTable)
	if err != nil {
		log.Fatal(err)
	}

	_, err = tx.Exec(chatTable)
	if err != nil {
		log.Fatal(err)
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}
}

func (d *Database) BeginTx() (*sql.Tx, error) {
	return d.conn.Begin()
}

func (d *Database) RollbackTx(tx *sql.Tx) error {
	return tx.Rollback()
}

func (d *Database) CommitTx(tx *sql.Tx) error {
	return tx.Commit()
}
