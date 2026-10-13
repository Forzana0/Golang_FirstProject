package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"thirdProgram/models"

	"github.com/go-chi/chi/v5"
)

func GetCategories(w http.ResponseWriter, r *http.Request) {

	userIDstr := chi.URLParam(r, "user_id")

	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	rows, err := Conn.Query(ctx, "SELECT name, id  FROM categories WHERE user_id IS NULL OR user_id = $1", userID)
	if err != nil {
		http.Error(w, "Error request to DB", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var categories []models.Category

	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.Name, &c.ID); err != nil {
			http.Error(w, "Error while reading from DB", http.StatusInternalServerError)
			return
		}
		categories = append(categories, c)
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(categories); err != nil {
		http.Error(w, "Error while encoding json", http.StatusInternalServerError)
		return
	}
}

func CreateCategory(w http.ResponseWriter, r *http.Request) {
	userIDstr := chi.URLParam(r, "user_id")

	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	category := &models.Category{}

	if err := json.NewDecoder(r.Body).Decode(category); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	_, err = Conn.Exec(ctx, "INSERT INTO categories (name, user_id) VALUES($1, $2)", category.Name, userID)
	if err != nil {
		http.Error(w, "Error while adding new category", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

}
