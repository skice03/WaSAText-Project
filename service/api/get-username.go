package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"wasatext/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	authHeader := r.Header.Get("Authorization")

	const bearerPrefix = "Bearer "
	if len(authHeader) <= len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Invalid authorization format"})
		return
	}

	// Extract the token from the Bearer prefix
	token := authHeader[len(bearerPrefix):]

	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Empty token"})
	}
	// Extract the user id from the url path
	getUid := ps.ByName("id")

	// Convert user id to integer
	reqUid, err := strconv.Atoi(getUid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Invalid or unauthorized user ID"})
		return
	}

	// Retrieve the user's name from the database
	username, err := rt.db.GetUsername(reqUid)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "User not found"})
		ctx.Logger.WithError(err).Error("Database fail")
		return
	}

	// Respond with the username
	response := struct {
		Username string `json:"username"`
	}{
		Username: username,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Error encoding response"})
	}
}
