package main

import (
	"fmt"
	"net/http"
	"go-personal-costs/internal/handlers"
	"go-personal-costs/internal/repository"
)

func main() {
	// Создаем репозиторий (он сам откроет expenses.db)
	repo, err := repository.NewExpenseRepository()
	if err != nil {
		panic(err)
	}
	handlers.Repo = repo

	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.ExpensesHandler)
	mux.HandleFunc("/edit", handlers.EditHandler)
	mux.HandleFunc("/delete", handlers.DeleteHandler)

	fmt.Println("Сервер Этапа 3 (SQLite + Repository) запущен на http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
