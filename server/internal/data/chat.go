package data

import (
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

func (d *Database) CreateChat(chatUserID, role, text, audio string) (*Entry, error) {
	id := uuid.New().String()
	_, err := d.conn.Exec("INSERT INTO chats (id, chat_user_id, role, text, audio) VALUES (?, ?, ?, ?, ?)",
		id, chatUserID, role, text, audio)
	if err != nil {
		return nil, err
	}
	return &Entry{ID: id, ChatUserID: chatUserID, Role: role, Text: text, Audio: audio}, nil
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
