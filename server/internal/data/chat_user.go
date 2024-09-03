package data

import (
	"database/sql"

	"github.com/google/uuid"
)

type ChatUser struct {
	ID     string `json:"id"`
	Secret string `json:"secret"`
}

func (d *Database) CreateChatUser(tx *sql.Tx, secret string) (*ChatUser, error) {
	id := uuid.New().String()
	_, err := tx.Exec("INSERT INTO chat_users (id, secret) VALUES (?, ?)", id, secret)
	if err != nil {
		return nil, err
	}

	return &ChatUser{ID: id, Secret: secret}, nil
}

func (d *Database) GetChatUser(id string) (*ChatUser, error) {
	var user ChatUser
	err := d.conn.QueryRow("SELECT id, secret FROM chat_users WHERE id = ?", id).Scan(&user.ID, &user.Secret)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
