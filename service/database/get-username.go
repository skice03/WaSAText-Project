package database

func (db *appdbimpl) GetUsername(userID int) (string, error) {
	var username string
	err := db.c.QueryRow(`
		SELECT username FROM users WHERE ID = ?`, userID).Scan(&username)
	if err != nil {
		return "", err
	}
	return username, nil
}
