package api

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
	"net/http"
	"wasatext/service/api/reqcontext"
	"wasatext/service/utils"

	"github.com/julienschmidt/httprouter"
)

// generating the api key
func generateApiKey() (string, error) {
	var apiKey []byte
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 16

	for i := 0; i < keyLength; i++ {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return "error generating the key", err
		}
		apiKey = append(apiKey, alphabet[index.Int64()])
	}

	return string(apiKey), nil
}

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var requestBody struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	username := requestBody.Name
	if !utils.ValidUsername(username) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Username must be alphanumeric and between 3 and 16 characters"})
		return
	}

	exists, err := rt.db.UserExistence(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Database error"})
		return
	}

	var userId int
	var apiKey string

	if !exists {
		apiKey, err = generateApiKey()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "Error generating API key"})
			return
		}

		userId, err = rt.db.CreateUser(username, apiKey)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "Error adding new user"})
			return
		}
	} else {
		userId, err = rt.db.GetUserId(username)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "Error retrieving user ID"})
			return
		}

		apiKey, err = rt.db.GetUserKey(userId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "Error retrieving API key"})
			return
		}
	}

	response := struct {
		Username string `json:"username"`
		UserId   int    `json:"userId"`
		APIKey   string `json:"apiKey"`
	}{
		Username: username,
		UserId:   userId,
		APIKey:   apiKey,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
}
