package api

import (
	"encoding/json"
	"net/http"
	"wasatext/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) newChat(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Checking for auth (bearer token)
	authHeader := r.Header.Get("Authorization")
	const bearerPrefix = "Bearer "
	if len(authHeader) <= len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Invalid authorization format"})
		return
	}

	token := authHeader[len(bearerPrefix):]
	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Empty token"})
		return
	}

	// Decoding the JSON req body into a struct
	var reqBody struct {
		Members []int `json:"members"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON provided"})
		return
	}

	// Checking for any members
	if len(reqBody.Members) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Missing required field: Members"})
		return
	}

	// Verifying if the member count isn't > than the max size declared in the API
	if len(reqBody.Members) > 2000 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Too many user IDs provided"})
		return
	}

	var chatId int
	var err error

	// Verifying the number of users (Treating the case for a private or a group chat)
	if len(reqBody.Members) == 2 {
		user1, err := rt.db.GetUsername(reqBody.Members[0])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "User not found"})
			ctx.Logger.WithError(err).Error("Database fail")
			return
		}

		user2, err := rt.db.GetUsername(reqBody.Members[1])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "User not found"})
			ctx.Logger.WithError(err).Error("Database fail")
			return
		}

		chatId, err = rt.db.NewChat("Chat between "+user1+" and "+user2, false)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create conversation"})
			ctx.Logger.WithError(err).Error("Database fail")
			return
		}
	} else {
		chatId, err = rt.db.NewChat("Group chat", true)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create conversation"})
			ctx.Logger.WithError(err).Error("Database fail")
			return
		}
	}

	// Add members to the chat
	for _, userId := range reqBody.Members {
		err := rt.db.AddChatMembers(userId, chatId)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "User not found"})
			ctx.Logger.WithError(err).Error("Database fail")
			return
		}
	}

	// The newly created chat
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]int{"chatId": chatId})
}
