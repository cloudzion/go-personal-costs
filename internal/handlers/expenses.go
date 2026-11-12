package handlers

import (
	"go-personal-costs/internal/models"
	"go-personal-costs/internal/repository"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

var Repo *repository.ExpenseRepository

type PageData struct {
	IsAuth         bool
	Expenses       []models.Expense
	Categories     []models.Category
	FilterCategory string
	FilterFrom     string
	FilterTo       string
	EditExpense    *models.Expense
}

func ExpensesHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекаем userID из контекста (его туда положил AuthMiddleware)
	userID := r.Context().Value(UserIDKey).(int)

	if r.Method == http.MethodGet {
		catFilter := r.URL.Query().Get("category")
		fromFilter := r.URL.Query().Get("from")
		toFilter := r.URL.Query().Get("to")

		expensesList, err := Repo.GetFiltered(userID, catFilter, fromFilter, toFilter)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		categoriesList, err := Repo.GetCategories()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := PageData{
			IsAuth:         true,
			Expenses:       expensesList,
			Categories:     categoriesList,
			FilterCategory: catFilter,
			FilterFrom:     fromFilter,
			FilterTo:       toFilter,
		}

		tmpl, err := template.ParseFiles("web/templates/layout.html", "web/templates/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.ExecuteTemplate(w, "layout", data)
		return
	}

	if r.Method == http.MethodPost {
		r.ParseForm()
		amount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)
		desc := r.FormValue("description")
		catID, _ := strconv.Atoi(r.FormValue("category_id"))
		date, err := time.Parse("2006-01-02", r.FormValue("date"))
		if err != nil {
			date = time.Now()
		}

		Repo.Create(userID, amount, desc, date, catID)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}

func EditHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(int)
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))

	if r.Method == http.MethodGet {
		exp, err := Repo.GetByID(id, userID)
		if err != nil {
			http.Error(w, "Запись не найдена или у вас нет прав", http.StatusNotFound)
			return
		}

		categoriesList, err := Repo.GetCategories()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := PageData{
			IsAuth:      true,
			EditExpense: exp,
			Categories:  categoriesList,
		}

		tmpl, err := template.ParseFiles("web/templates/layout.html", "web/templates/edit.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.ExecuteTemplate(w, "layout", data)
		return
	}

	if r.Method == http.MethodPost {
		r.ParseForm()
		amount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)
		desc := r.FormValue("description")
		catID, _ := strconv.Atoi(r.FormValue("category_id"))
		date, _ := time.Parse("2006-01-02", r.FormValue("date"))

		Repo.Update(id, userID, amount, desc, date, catID)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(int)
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	Repo.Delete(id, userID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
