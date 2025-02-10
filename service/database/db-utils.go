package database

// Verifying user existence
func (db *appdbimpl) UserExists(username string) (bool, error) {
	var exists bool
	err := db.c.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)`, username).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// Creating a new user
func (db *appdbimpl) CreateUser(username string, securityKey string) (int, error) {
	res, err := db.c.Exec(`INSERT INTO users (username, security_key) VALUES (?, ?)`, username, securityKey)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

// Retrieving the security key via user id
func (db *appdbimpl) GetUserKey(userId int) (string, error) {
	var securityKey string
	err := db.c.QueryRow(`SELECT security_key FROM users WHERE ID = ?`, userId).Scan(&securityKey)
	if err != nil {
		return "", err
	}
	return securityKey, nil
}

// Retrieving the user id via security key
func (db *appdbimpl) GetUserIdByKey(securityKey string) (int, error) {
	var userId int
	err := db.c.QueryRow(`SELECT id FROM users WHERE security_key = ?`, securityKey).Scan(&userId)
	if err != nil {
		return 0, err
	}
	return userId, nil
}

// Getting the username of an user
func (db *appdbimpl) GetUsername(userId int) (string, error) {
	var username string
	err := db.c.QueryRow(`SELECT username FROM users WHERE ID = ?`, userId).Scan(&username)
	if err != nil {
		return "", err
	}
	return username, nil
}

// Retrieving the user id via username
func (db *appdbimpl) GetUserIdByUsername(username string) (int, error) {
	var userId int
	err := db.c.QueryRow(`SELECT id FROM users WHERE username = ?`, username).Scan(&userId)
	if err != nil {
		return 0, err
	}
	return userId, nil
}

// Update the username
func (db *appdbimpl) UpdateUsername(userId int, newUsername string) error {
	_, err := db.c.Exec("UPDATE users SET username = ? WHERE id = ?", newUsername, userId)
	return err
}

// Fetching the conversations of an user
func (db *appdbimpl) GetUserChats(userId int) ([]int, error) {
	var chatList []int

	rows, err := db.c.Query(`SELECT chat_id FROM chat_members WHERE user_id = ?`, userId)
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

// Create a new conversation (either private / group chat)
func (db *appdbimpl) NewChat(chatName string, groupChat bool) (int, error) {
	res, err := db.c.Exec(`INSERT INTO chats (name, group_chat) VALUES (?, ?)`, chatName, groupChat)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

// Add an user to the newly created chat
func (db *appdbimpl) AddChatMembers(userId int, chatId int) error {
	_, err := db.c.Exec(`INSERT INTO chat_members (user_id, chat_id) VALUES (?, ?)`, userId, chatId)
	if err != nil {
		return err
	}
	return nil
}

// Checking if the user belongs to the conversation
func (db *appdbimpl) ChatMember(userId int, chatId int) (bool, error) {
	var exists bool
	err := db.c.QueryRow(`SELECT EXISTS(SELECT 1 FROM chat_members WHERE user_id = ? AND chat_id = ?)`, userId, chatId).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// Checking if the conversation is a group
func (db *appdbimpl) GroupChat(chatId int) (bool, error) {
	var groupChat bool
	err := db.c.QueryRow(`SELECT group_chat FROM chats WHERE id = ?`, chatId).Scan(&groupChat)
	if err != nil {
		return false, err
	}
	return groupChat, nil
}

// Updating the conversation name
func (db *appdbimpl) SetChatName(chatId int, newName string) error {
	_, err := db.c.Exec(`UPDATE chats SET name = ? WHERE id = ?`, newName, chatId)
	if err != nil {
		return err
	}
	return nil
}

// Retrieving the conversation name
func (db *appdbimpl) GetChatName(chatId int) (string, error) {
	var chatName string

	err := db.c.QueryRow(`SELECT name FROM chats WHERE id = ?`, chatId).Scan(&chatName)
	if err != nil {
		return "", err
	}

	return chatName, nil
}

// Retrieve the members from a conversation
func (db *appdbimpl) GetChatMembers(chatId int) ([]int, error) {
	var userList []int

	rows, err := db.c.Query(`SELECT user_id FROM chat_members WHERE chat_id = ?`, chatId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var userId int
		if err := rows.Scan(&userId); err != nil {
			return nil, err
		}
		userList = append(userList, userId)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userList, nil
}

// Removing a member from a chat
func (db *appdbimpl) RemoveChatMember(userId int, chatId int) error {
	_, err := db.c.Exec(`
		DELETE FROM chat_members WHERE user_id = ? AND chat_id = ?`, userId, chatId)
	if err != nil {
		return err
	}

	return nil
}
