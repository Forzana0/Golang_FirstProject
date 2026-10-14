package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func parseUserID(w http.ResponseWriter, r *http.Request) (int, bool) {
	userIDstr := chi.URLParam(r, "user_id")
	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return 0, false
	}
	return userID, true
}
