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

func parseExpenseID(w http.ResponseWriter, r *http.Request) (int, bool) {
	expenseIDstr := chi.URLParam(r, "expense_id")
	expenseID, err := strconv.Atoi(expenseIDstr)
	if err != nil {
		http.Error(w, "Invalid expense ID format", http.StatusBadRequest)
		return 0, false
	}
	return expenseID, true
}
