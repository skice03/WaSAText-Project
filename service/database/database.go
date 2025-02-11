package database

import (
	"database/sql"
	"errors"
	"fmt"
	"image/gif"
	"time"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	UserExists(username string) (bool, error)
	CreateUser(username string, securityKey string) (int, error)
	GetUserKey(userId int) (string, error)
	GetUserIdByKey(securityKey string) (int, error)
	GetUsername(userId int) (string, error)
	GetUserIdByUsername(username string) (int, error)
	UpdateUsername(userId int, newUsername string) error
	GetUserChats(userId int) ([]int, error)
	NewChat(chatName string, groupChat bool) (int, error)
	AddChatMember(userId int, chatId int) error
	ChatMember(userId int, chatId int) (bool, error)
	GroupChat(chatId int) (bool, error)
	SetChatName(chatId int, newName string) error
	GetChatName(chatId int) (string, error)
	GetChatMembers(chatId int) ([]int, error)
	GetUserCount() (int, error)
	RemoveChatMember(userId int, chatId int) error
	AddComment(textContent string, senderId int, messageId int) error
	RemoveComment(senderId int, messageId int) error
	SendMessage(chatId int, senderId int, textContent string, forwarded bool, timestamp time.Time) error
	DeleteMessage(messageId int) error
	ViewMessage(userId int, messageId int) error
	ReceiveMessage(userId int, messageId int) error
	GetChatMessages(chatId int) ([]int, error)
	GetMessageComments(messageId int) ([]int, []string, error)
	SeenMessage(messageId int) ([]int, error)
	GetMessage(messageId int) (int, string, bool, time.Time, error)
	Ping() error
}

type User struct {
	ID          uint64
	Username    string
	GifImage    *gif.GIF
	SecurityKey string
}

type Chat struct {
	ID       uint64
	Name     string
	Members  []uint64
	GifImage *gif.GIF
}

type Message struct {
	ID          uint64
	TextContent string
	GifContent  *gif.GIF
	Status      string
	Timestamp   *time.Time
	SenderId    uint64
}

type appdbimpl struct {
	c *sql.DB
}

func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Checking whether the table exists or not ; if the db is empty we'll create the structure
	var tableName string

	// Users table
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='users';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE users (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			security_key TEXT NOT NULL UNIQUE,
			gif_photo BLOB NULL
		);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure (users): %w", err)
		}
	}

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='chats';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE chats (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			gif_photo BLOB NULL,
			group_chat BOOL
		);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure (chats): %w", err)
		}
	}

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='chat_members';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE chat_members (
			chat_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			PRIMARY KEY (chat_id, user_id),
			FOREIGN KEY (chat_id) REFERENCES chats(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure (chat_members): %w", err)
		}
	}

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='messages';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE messages (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			chat_id INTEGER NOT NULL,
			sender_id INTEGER NOT NULL,
			text_message TEXT,
			gif_photo BLOB,
			timestamp DATETIME NOT NULL,
			forwarded BOOL,
			FOREIGN KEY (chat_id) REFERENCES chats(id),
			FOREIGN KEY (sender_id) REFERENCES users(id)
		);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure (messages): %w", err)
		}
	}

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='message_status';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE message_status (
			message_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			sent BOOL,
			seen BOOL,
			comment TEXT NOT NULL,
			PRIMARY KEY (user_id, message_id),
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (message_id) REFERENCES chats(id)
		);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure (message_status): %w", err)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
