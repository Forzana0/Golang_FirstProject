package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"thirdProgram/models"

	"github.com/go-chi/chi/v5"
)

func AddExpense(w http.ResponseWriter, r *http.Request) {
	userIDstr := chi.URLParam(r, "user_id")

	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	expanse := &models.Expense{}
	if err = json.NewDecoder(r.Body).Decode(expanse); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	var currentBalance float64
	err = Conn.QueryRow(ctx, "SELECT balance FROM users WHERE id = $1", userID).Scan(&currentBalance)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if expanse.Price > currentBalance {
		http.Error(w, "Not enough money", http.StatusBadRequest)
		return
	}

	_, err = Conn.Exec(ctx, "INSERT INTO expenses(name, description, price, user_id, category_id) VALUES($1, $2, $3, $4, $5)", expanse.Name, expanse.Description, expanse.Price, userID, expanse.CategoryID)
	if err != nil {
		http.Error(w, "Error while adding new expense", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
