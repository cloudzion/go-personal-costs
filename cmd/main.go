package main

import (
	"fmt"
	"go-personal-costs/internal/handlers"
	"go-personal-costs/internal/repository"
	"net/http"
)

func main() {
	repo, err := repository.NewExpenseRepository()
	if err != nil {
		panic(err)
	}
	handlers.Repo = repo

	mux := http.NewServeMux()

	// Публичные маршруты (доступны всем)
	mux.HandleFunc("/register", handlers.RegisterHandler)
	mux.HandleFunc("/login", handlers.LoginHandler)
	mux.HandleFunc("/logout", handlers.LogoutHandler)

	// Защищенные маршруты (проходят через AuthMiddleware)
	mux.HandleFunc("/", handlers.AuthMiddleware(handlers.ExpensesHandler))
	mux.HandleFunc("/edit", handlers.AuthMiddleware(handlers.EditHandler))
	mux.HandleFunc("/delete", handlers.AuthMiddleware(handlers.DeleteHandler))

	fmt.Println("Сервер Этапа 5 (с авторизацией) запущен на http://localhost:8080")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	}
}
