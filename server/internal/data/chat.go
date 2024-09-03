package data

import (
	"database/sql"
	"strings"

	"github.com/google/uuid"
)

type Chat struct {
	Text  string `json:"text"`
	Audio string `json:"audio"`
}

type Entry struct {
	ID         string `json:"id"`
	ChatUserID string `json:"chat_user_id"`
	Role       string `json:"role"`
	Text       string `json:"text"`
	Audio      string `json:"audio"`
}

func (d *Database) CreateChat(tx *sql.Tx, chatUserID, role, text, audio string) (*Entry, error) {
	id := uuid.New().String()
	_, err := tx.Exec("INSERT INTO chats (id, chat_user_id, role, text, audio) VALUES (?, ?, ?, ?, ?)",
		id, chatUserID, role, text, audio)
	if err != nil {
		return nil, err
	}
	return &Entry{ID: id, ChatUserID: chatUserID, Role: role, Text: text, Audio: audio}, nil
}

func (d *Database) CreateChats(tx *sql.Tx, chatUserID string, chats []Entry) ([]Entry, error) {
	query := "INSERT INTO chats (id, chat_user_id, role, text, audio) VALUES "
	var values []interface{}
	placeholders := make([]string, len(chats))

	for i, chat := range chats {
		chat.ID = uuid.New().String()
		chat.ChatUserID = chatUserID

		placeholders[i] = "(?, ?, ?, ?, ?)"

		values = append(values, chat.ID, chat.ChatUserID, chat.Role, chat.Text, chat.Audio)
	}

	query += strings.Join(placeholders, ",")

	_, err := tx.Exec(query, values...)
	if err != nil {
		return nil, err
	}

	return chats, nil
}

func (d *Database) GetChatsByChatUserID(chatUserID string) ([]Entry, error) {
	rows, err := d.conn.Query("SELECT id, chat_user_id, role, text, audio FROM chats WHERE chat_user_id = ?", chatUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []Entry
	for rows.Next() {
		var chat Entry
		err := rows.Scan(&chat.ID, &chat.ChatUserID, &chat.Role, &chat.Text, &chat.Audio)
		if err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}
	return chats, nil
}
