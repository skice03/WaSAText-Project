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
		INSERT INTO users (username, security_key) 
		VALUES (?, ?)`, username, securityKey)
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
