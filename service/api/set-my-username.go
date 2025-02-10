package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"wasatext/service/api/reqcontext"
	"wasatext/service/utils"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) setMyUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	const bearerPrefix = "Bearer "

	// Check authorization header
	token := r.Header.Get("Authorization")
	if len(token) <= len(bearerPrefix) || token[:len(bearerPrefix)] != bearerPrefix {
		writeErrorResponse(w, http.StatusUnauthorized, "You are not logged in. Please log in to continue.")
		return
	}
	token = token[len(bearerPrefix):]

	// Extract and validate user id from url
	userID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid user ID.")
		return
	}

	// Verify if the token matches the user's actual key
	securityKey, err := rt.db.GetUserKey(userID)
	if err != nil || securityKey != token {
		writeErrorResponse(w, http.StatusUnauthorized, "Invalid session. Please log in again.")
		return
	}

	// Parse and validate new username from request body
	var requestBody struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body.")
		return
	}
	if !utils.ValidUsername(requestBody.Username) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Username must be alphanumeric and between 3 and 16 characters"})
		return
	}

	// Check if the username is already taken by another user
	exists, err := rt.db.UserExists(requestBody.Username)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to check username availability.")
		return
	}
	if exists {
		writeErrorResponse(w, http.StatusConflict, "Username already taken. Please choose a different one.")
		return
	}

	// Update the username in the database
	if err := rt.db.UpdateUsername(userID, requestBody.Username); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to update username.")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeErrorResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}
