package database

// checking whether an user exists or not
func (db *appdbimpl) UserExistence(username string) (bool, error) {
	var count int
	err := db.c.QueryRow(`
		SELECT COUNT(1) 
		FROM users 
		WHERE username = ?`, username).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// creating a new user
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

// retrieving the username by user id
func (db *appdbimpl) GetUserId(username string) (int, error) {
	var userId int
	err := db.c.QueryRow(`
		SELECT id FROM users WHERE username = ?`, username).Scan(&userId)
	if err != nil {
		return 0, err
	}
	return userId, nil
}

// retrieving the api key by user id
func (db *appdbimpl) GetUserKey(userId int) (string, error) {
	var key string
	err := db.c.QueryRow(`
		SELECT security_key FROM users WHERE ID = ?`, userId).Scan(&key)
	if err != nil {
		return "", err
	}
	return key, nil
}
