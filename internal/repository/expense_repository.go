package repository

import (
	"database/sql"
	"time"
	"go-personal-costs/internal/models"
	_ "modernc.org/sqlite" // Драйвер SQLite
)

type ExpenseRepository struct {
	db *sql.DB
}

// Создаем новый репозиторий и подключаем базу данных в файл expenses.db
func NewExpenseRepository() (*ExpenseRepository, error) {
	db, err := sql.Open("sqlite", "expenses.db")
	if err != nil {
		return nil, err
	}

	// Сразу создаем таблицу, если её нет
	statement := `
	CREATE TABLE IF NOT EXISTS expenses (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		amount REAL NOT NULL,
		description TEXT NOT NULL,
		date TEXT NOT NULL
	);`
	_, err = db.Exec(statement)
	if err != nil {
		return nil, err
	}

	return &ExpenseRepository{db: db}, nil
}

// 1. Получить все траты (GetAll)
func (r *ExpenseRepository) GetAll() ([]models.Expense, error) {
	rows, err := r.db.Query("SELECT id, amount, description, date FROM expenses ORDER BY date DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []models.Expense
	for rows.Next() {
		var exp models.Expense
		var dateStr string
		err := rows.Scan(&exp.ID, &exp.Amount, &exp.Description, &dateStr)
		if err != nil {
			return nil, err
		}
		exp.Date, _ = time.Parse("2006-01-02", dateStr)
		expenses = append(expenses, exp)
	}
	return expenses, nil
}

// 2. Создать трату (Create) — Используем подготовленные запросы с "?" от SQL-инъекций
func (r *ExpenseRepository) Create(amount float64, desc string, date time.Time) error {
	query := "INSERT INTO expenses (amount, description, date) VALUES (?, ?, ?)"
	_, err := r.db.Exec(query, amount, desc, date.Format("2006-01-02"))
	return err
}

// 3. Получить по ID для редактирования (GetByID)
func (r *ExpenseRepository) GetByID(id int) (*models.Expense, error) {
	query := "SELECT id, amount, description, date FROM expenses WHERE id = ?"
	row := r.db.QueryRow(query, id)

	var exp models.Expense
	var dateStr string
	err := row.Scan(&exp.ID, &exp.Amount, &exp.Description, &dateStr)
	if err != nil {
		return nil, err
	}
	exp.Date, _ = time.Parse("2006-01-02", dateStr)
	return &exp, nil
}

// 4. Обновить трату (Update)
func (r *ExpenseRepository) Update(id int, amount float64, desc string, date time.Time) error {
	query := "UPDATE expenses SET amount = ?, description = ?, date = ? WHERE id = ?"
	_, err := r.db.Exec(query, amount, desc, date.Format("2006-01-02"), id)
	return err
}

// 5. Удалить трату (Delete)
func (r *ExpenseRepository) Delete(id int) error {
	query := "DELETE FROM expenses WHERE id = ?"
	_, err := r.db.Exec(query, id)
	return err
}
