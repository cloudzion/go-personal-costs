package handlers

import (
	"html/template"
	"net/http"
	"strconv"
	"time"

	"go-personal-costs/internal/models"
)

// Хранилище трат в оперативной памяти для Этапа 2
var Expenses []models.Expense

// ExpensesHandler управляет выводом списка трат (GET) и их добавлением (POST)
func ExpensesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Читаем файлы напрямую из web/templates/
		tmpl, err := template.ParseFiles("web/templates/layout.html", "web/templates/index.html")
		if err != nil {
			http.Error(w, "Ошибка загрузки шаблонов: "+err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.ExecuteTemplate(w, "layout", Expenses)
		if err != nil {
			http.Error(w, "Ошибка отображения страницы: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Ошибка обработки данных формы", http.StatusBadRequest)
			return
		}

		amountStr := r.FormValue("amount")
		description := r.FormValue("description")
		dateStr := r.FormValue("date")

		amount, _ := strconv.ParseFloat(amountStr, 64)
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			date = time.Now()
		}

		newExpense := models.Expense{
			Amount:      amount,
			Description: description,
			Date:        date,
		}

		Expenses = append(Expenses, newExpense)

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
}
