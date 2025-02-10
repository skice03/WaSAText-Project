package database

// Verifying user existence
func (db *appdbimpl) UserExists(username string) (bool, error) {
	var exists bool
	err := db.c.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)`, username).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// Creating a new user
func (db *appdbimpl) CreateUser(username string, securityKey string) (int, error) {
	res, err := db.c.Exec(`
		INSERT INTO users (username, security_key) VALUES (?, ?)`, username, securityKey)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

// Retrieving the username by user id
func (db *appdbimpl) GetUserId(username string) (int, error) {
	var userId int
	err := db.c.QueryRow(`
		SELECT id FROM users WHERE username = ?`, username).Scan(&userId)
	if err != nil {
		return 0, err
	}
	return userId, nil
}

// Retrieving the security key by user id
func (db *appdbimpl) GetUserKey(userId int) (string, error) {
	var securityKey string
	err := db.c.QueryRow(`
		SELECT security_key FROM users WHERE ID = ?`, userId).Scan(&securityKey)
	if err != nil {
		return "", err
	}
	return securityKey, nil
}

// Get the username
func (db *appdbimpl) GetUsername(userID int) (string, error) {
	var username string
	err := db.c.QueryRow(`
		SELECT username FROM users WHERE ID = ?`, userID).Scan(&username)
	if err != nil {
		return "", err
	}
	return username, nil
}

// Update the username
func (db *appdbimpl) UpdateUsername(userID int, newUsername string) error {
	_, err := db.c.Exec("UPDATE users SET username = ? WHERE id = ?", newUsername, userID)
	return err
}

// Create a new conversation (either private / group chat)
func (db *appdbimpl) NewChat(name string, group bool) (int, error) {
	res, err := db.c.Exec(`
		INSERT INTO chats (name, group_chat) VALUES (?, ?)`, name, group)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

// Add an user to the newly created chat
func (db *appdbimpl) AddChatMembers(userId int, chatId int) error {
	_, err := db.c.Exec(`
		INSERT INTO chat_members (chat_id, user_id) VALUES (?, ?)`,
		chatId, userId)
	if err != nil {
		return err
	}
	return nil
}

// Checking if the conversations is a group
func (db *appdbimpl) GroupChat(chatId int) (bool, error) {
	var groupChat bool
	err := db.c.QueryRow(`
		SELECT group_chat FROM chats WHERE id = ?`, chatId).Scan(&groupChat)
	if err != nil {
		return false, err
	}
	return groupChat, nil
}

// Checking if the conversation is private (between 2 users)
func (db *appdbimpl) PrivateChat(userId int, chatId int) (bool, error) {
	var memberCount int
	err := db.c.QueryRow(`
		SELECT COUNT(1) FROM chat_members WHERE user_id = ? AND chat_id = ?`, userId, chatId).Scan(&memberCount)
	if err != nil {
		return false, err
	}
	return memberCount > 0, nil
}

// Updating the conversation name
func (db *appdbimpl) SetChatName(chatId int, newName string) error {
	_, err := db.c.Exec(`
		UPDATE chats SET name = ? WHERE id = ?`, newName, chatId)
	if err != nil {
		return err
	}
	return nil
}

// Fetching the conversations of an user
func (db *appdbimpl) GetUserChats(userId int) ([]int, error) {
	var chatList []int

	query := `SELECT chat_id FROM chat_members WHERE user_id = ?`
	rows, err := db.c.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var chatId int
		if err := rows.Scan(&chatId); err != nil {
			return nil, err
		}
		chatList = append(chatList, chatId)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return chatList, nil
}
