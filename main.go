package main

import (
	"context"
	"net/http"
	"thirdProgram/database"
	"thirdProgram/handlers"

	"github.com/go-chi/chi/v5"
)

func main() {
	ctx := context.Background()
	var err error
	handlers.Conn, err = database.Connect(ctx, "postgres://postgres:@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer handlers.Conn.Close(ctx)

	r := chi.NewRouter()

	//User_Handlers
	r.Get("/users", handlers.GetUsers)
	r.Post("/users", handlers.CreateUser)
	r.Delete("/users", handlers.DeleteUser)
	//Category_handlers
	r.Get("/users/{user_id}/categories", handlers.GetCategories)
	r.Post("/users/{user_id}/categories", handlers.CreateCategory)
	//Expense_handler
	r.Post("/users/{user_id}/expanses", handlers.AddExpense)
	http.ListenAndServe(":9091", r)
}
