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
	//here
	//doLogin(string) (User, error)

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

	// checking whether the table exists or not ; if the db is empty we'll create the structure
	var tableName string

	// users table
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
	/*	usersStmt := `CREATE TABLE IF NOT EXISTS Users (
				username TEXT NOT NULL,
				gif_image BLOB,
				security_key TEXT NOT NULL,
				id INTEGER NOT NULL PRIMARY KEY
				);`
		if _, err := db.Exec(usersStmt); err != nil {
			return nil, fmt.Errorf("error creating Users table: %w", err)
		}
	*/

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='chats';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE chats (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			gif_photo BLOB NULL
		);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure (chats): %w", err)
		}
	}
	/*	chatsStmt := `CREATE TABLE IF NOT EXISTS Chats (
				name TEXT NOT NULL,
				gif_image BLOB,
				id INTEGER NOT NULL PRIMARY KEY
				);`
		if _, err := db.Exec(chatsStmt); err != nil {
			return nil, fmt.Errorf("error creating Chats table: %w", err)
		}
	*/

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
	/*	chatMembersStmt := `CREATE TABLE IF NOT EXISTS ChatMembers (
					chat_id INTEGER NOT NULL,
					user_id INTEGER NOT NULL,
					PRIMARY KEY (chat_id, user_id),
					FOREIGN KEY (chat_id) REFERENCES Chats(id),
					FOREIGN KEY (user_id) REFERENCES Users(id)
					);`
		if _, err := db.Exec(chatMembersStmt); err != nil {
			return nil, fmt.Errorf("error creating ChatMembers table: %w", err)
		}
	*/

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='messages';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE messages (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			chat_id INTEGER NOT NULL,
			sender_id INTEGER NOT NULL,
			text_message TEXT NOT NULL,
			gif_photo BLOB NOT NULL,
			timestamp DATETIME NOT NULL,
			status TEXT NOT NULL,
			FOREIGN KEY (chat_id) REFERENCES chats(id),
			FOREIGN KEY (sender_id) REFERENCES users(id)
		);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure (messages): %w", err)
		}
	}

	/*	messagesStmt := `CREATE TABLE IF NOT EXISTS Messages (
				text_content TEXT,
				gif_image BLOB,
				status TEXT NOT NULL,
				timestamp DATETIME NOT NULL,
				sender_id INTEGER NOT NULL,
				id INTEGER NOT NULL PRIMARY KEY,
				FOREIGN KEY (chat_id) REFERENCES Chats(id),
				FOREIGN KEY (sender_id) REFERENCES Users(id)
				);`
		if _, err := db.Exec(messagesStmt); err != nil {
			return nil, fmt.Errorf("error creating Messages table: %w", err)
		}
	*/
	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
