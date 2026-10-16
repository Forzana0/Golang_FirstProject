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
	err := Pool.QueryRow(ctx, "SELECT balance FROM users WHERE id = $1", userID).Scan(&currentBalance)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if expanse.Price > currentBalance {
		http.Error(w, "Not enough money", http.StatusBadRequest)
		return
	}

	tx, err := Pool.Begin(ctx)
	if err != nil {
		http.Error(w, "Error starting transaction for adding expense", http.StatusInternalServerError)
		return
	}

	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, "UPDATE users SET balance = balance - $1 WHERE id = $2", expanse.Price, userID)
	if err != nil {
		http.Error(w, "Error while updating the balance", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(ctx, "INSERT INTO expenses(name, description, price, user_id, category_id) VALUES($1, $2, $3, $4, $5)", expanse.Name, expanse.Description, expanse.Price, userID, expanse.CategoryID)
	if err != nil {
		http.Error(w, "Error while adding new expense", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(ctx); err != nil {
		http.Error(w, "Error commiting transaction", http.StatusInternalServerError)
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

	rows, err := Pool.Query(ctx, "SELECT id, name, description, price, created_at, user_id, category_id FROM expenses WHERE user_id = $1", userID)
	if err != nil {
		http.Error(w, "Error request to DB", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var e models.Expense
		if err := rows.Scan(&e.ID, &e.Name, &e.Description, &e.Price, &e.CreatedAt, &e.UserID, &e.CategoryID); err != nil {
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

func DeleteExpense(w http.ResponseWriter, r *http.Request) {
	userID, ok := parseUserID(w, r)
	if !ok {
		return
	}

	expenseID, ok := parseExpenseID(w, r)
	if !ok {
		return
	}

	ctx := r.Context()

	var price float64
	err := Pool.QueryRow(ctx, "SELECT price FROM expenses WHERE id = $1", expenseID).Scan(&price)
	if err != nil {
		http.Error(w, "Expense ID not found", http.StatusNotFound)
		return
	}

	tx, err := Pool.Begin(ctx)
	if err != nil {
		http.Error(w, "Error starting transaction for deleting expense", http.StatusInternalServerError)
		return
	}

	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, "UPDATE users SET balance = balance + $1 WHERE id = $2", price, userID)
	if err != nil {
		http.Error(w, "Error while updating the balance", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(ctx, "DELETE FROM expenses WHERE id = $1 AND user_id = $2", expenseID, userID)
	if err != nil {
		http.Error(w, "Error while deleting expense", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(ctx); err != nil {
		http.Error(w, "Error while commitin transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
