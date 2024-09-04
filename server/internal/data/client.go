package data

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

type Database struct {
	conn *sql.DB
}

func New(dbPath string) *Database {
	var err error
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	migrate(db)

	return &Database{conn: db}
}

func migrate(db *sql.DB) {
	chatUserTable := `CREATE TABLE IF NOT EXISTS chat_users (
		id VARCHAR PRIMARY KEY,
		secret VARCHAR NOT NULL,
		language VARCHAR DEFAULT 'en'
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
