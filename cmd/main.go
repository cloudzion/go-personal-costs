package main

import (
	"fmt"
	"net/http"

	"go-personal-costs/internal/handlers"
)

func main() {
	mux := http.NewServeMux()

	// Подключаем наш обработчик интерфейса
	mux.HandleFunc("/", handlers.ExpensesHandler)

	fmt.Println("Сервер Этапа 2 успешно запущен на http://localhost:8080")
	
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Printf("Ошибка запуска сервера: %v\n", err)
	}
}
