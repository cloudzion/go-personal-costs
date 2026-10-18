package handlers

import (
	"html/template"
	"net/http"
	"strconv"
	"time"
	"go-personal-costs/internal/repository"
)

var Repo *repository.ExpenseRepository

func ExpensesHandler(w http.ResponseWriter, r *http.Request) {
	// Вывод списка всех трат (GET на "/")
	if r.Method == http.MethodGet {
		expensesList, err := Repo.GetAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles("web/templates/layout.html", "web/templates/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.ExecuteTemplate(w, "layout", expensesList)
		return
	}

	// Создание новой траты (POST на "/")
	if r.Method == http.MethodPost {
		r.ParseForm()
		amount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)
		desc := r.FormValue("description")
		date, err := time.Parse("2006-01-02", r.FormValue("date"))
		if err != nil {
			date = time.Now()
		}

		Repo.Create(amount, desc, date)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}

// Обработчик страницы редактирования (/edit)
func EditHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))

	if r.Method == http.MethodGet {
		exp, err := Repo.GetByID(id)
		if err != nil {
			http.Error(w, "Запись не найдена", http.StatusNotFound)
			return
		}

		tmpl, err := template.ParseFiles("web/templates/layout.html", "web/templates/edit.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.ExecuteTemplate(w, "layout", exp)
		return
	}

	if r.Method == http.MethodPost {
		r.ParseForm()
		amount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)
		desc := r.FormValue("description")
		date, _ := time.Parse("2006-01-02", r.FormValue("date"))

		Repo.Update(id, amount, desc, date)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// Обработчик удаления (/delete)
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	Repo.Delete(id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
