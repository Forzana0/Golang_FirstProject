package handlers

import (
	"encoding/json"
	"net/http"
	"thirdProgram/models"
)

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var users []models.User

	rows, err := h.pool.Query(ctx, "SELECT * FROM users ORDER BY id")
	if err != nil {
		http.Error(w, "Error request to DB", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Balance); err != nil {
			http.Error(w, "Error while reading from DB", http.StatusInternalServerError)
			return
		}
		users = append(users, u)
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Error encoding json", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := &models.User{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	_, err := h.pool.Exec(ctx, "INSERT INTO users (username, balance) VALUES($1, $2)", user.Username, user.Balance)
	if err != nil {
		http.Error(w, "Error while adding new user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := &models.User{}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if user.ID == 0 {
		http.Error(w, "UserId can not be 0", http.StatusBadRequest)

	}

	_, err := h.pool.Exec(ctx, "DELETE FROM users WHERE id =$1", user.ID)
	if err != nil {
		http.Error(w, "Error while deleting user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
