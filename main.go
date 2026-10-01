package main

import (
	"fmt"
	"net/http"
)

// Обработчик главной страницы
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Добро пожаловать в веб-приложение «Учёт личных трат»!")
}

// Обработчик страницы "О проекте"
func aboutHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Этот проект поможет вам контролировать свои расходы, распределять их по категориям и смотреть статистику.")
}

// Обработчик страницы проверки связи
func pingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "405 Method Not Allowed. Используйте метод GET.")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "pong")
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/about", aboutHandler)
	mux.HandleFunc("/ping", pingHandler)

	fmt.Println("Сервер успешно запущен на http://localhost:8080")
	
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Printf("Ошибка запуска сервера: %v\n", err)
	}
}
