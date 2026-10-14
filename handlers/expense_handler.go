package handlers

import (
	"encoding/json"
	"net/http"
	"thirdProgram/models"
)

func AddExpense(w http.ResponseWriter, r *http.Request) {

	userID, ok := parseUserID(w, r)
	if !ok {
		return
	}

	expanse := &models.Expense{}
	if err := json.NewDecoder(r.Body).Decode(expanse); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	var currentBalance float64
	err := Conn.QueryRow(ctx, "SELECT balance FROM users WHERE id = $1", userID).Scan(&currentBalance)
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

func GetExpenses(w http.ResponseWriter, r *http.Request) {
	userID, ok := parseUserID(w, r)
	if !ok {
		return
	}

	var expenses []models.Expense
	ctx := r.Context()

	rows, err := Conn.Query(ctx, "SELECT id, name, description, price, created_at FROM expenses WHERE user_id = $1", userID)
	if err != nil {
		http.Error(w, "Error request to DB", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var e models.Expense
		if err := rows.Scan(&e.ID, &e.Name, &e.Description, &e.Price, &e.CreatedAt); err != nil {
			http.Error(w, "Error while reading from DB", http.StatusInternalServerError)
			return
		}
		expenses = append(expenses, e)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(expenses); err != nil {
		http.Error(w, " Error encoding json", http.StatusInternalServerError)
		return
	}
}
