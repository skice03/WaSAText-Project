package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"wasatext/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

// Helper function to return JSON errors
func returnErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// Helper function for token extraction
func AuthToken(r *http.Request) (string, bool) {
	authHeader := r.Header.Get("Authorization")
	const bearerPrefix = "Bearer "

	if len(authHeader) <= len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		return "", false
	}

	token := authHeader[len(bearerPrefix):]
	return token, token != ""
}

func (rt *_router) addToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Checking the token
	token, valid := AuthToken(r)
	if !valid {
		returnErrorResponse(w, http.StatusUnauthorized, "Invalid authorization format")
		return
	}

	// Parsing the chat id
	chatId, err := strconv.Atoi(ps.ByName("chatId"))
	if err != nil {
		returnErrorResponse(w, http.StatusBadRequest, "Invalid conversation id")
		return
	}

	// Validate user authentication
	userId, err := rt.db.GetUserIdByKey(token)
	if err != nil {
		returnErrorResponse(w, http.StatusUnauthorized, "Auth error")
		return
	}

	// Check if user is a member of the chat
	isMember, err := rt.db.ChatMember(userId, chatId)
	if err != nil || !isMember {
		returnErrorResponse(w, http.StatusUnauthorized, "Auth error")
		return
	}

	// Parse JSON request body
	var reqBody struct {
		Members []int `json:"members"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		returnErrorResponse(w, http.StatusBadRequest, "Invalid JSON provided")
		return
	}

	if len(reqBody.Members) == 0 {
		returnErrorResponse(w, http.StatusBadRequest, "Missing required field: userIds")
		return
	}

	if len(reqBody.Members) > 2000 {
		returnErrorResponse(w, http.StatusBadRequest, "Too many user IDs provided")
		return
	}

	// Check if the chat is a group
	isGroup, err := rt.db.GroupChat(chatId)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to check if chat is a group")
		returnErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	if !isGroup {
		returnErrorResponse(w, http.StatusForbidden, "Not a group")
		return
	}

	// Add users to group
	for _, userId := range reqBody.Members {
		err := rt.db.AddChatMember(userId, chatId)
		if err != nil {
			ctx.Logger.WithError(err).Error("Failed to add user to group")
			returnErrorResponse(w, http.StatusNotFound, "User not found")
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
