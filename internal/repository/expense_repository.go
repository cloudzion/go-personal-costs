package repository

import (
	"database/sql"
	"go-personal-costs/internal/models"
	"time"

	_ "modernc.org/sqlite"
)

type ExpenseRepository struct {
	db *sql.DB
}

func NewExpenseRepository() (*ExpenseRepository, error) {
	db, err := sql.Open("sqlite", "expenses.db")
	if err != nil {
		return nil, err
	}

	// Создаем таблицу категорий и заполняем базовыми значениями
	catQuery := `
	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE
	);
	INSERT OR IGNORE INTO categories (name) VALUES ('Еда'), ('Транспорт'), ('Жильё'), ('Развлечения'), ('Прочее');`
	
	if _, err = db.Exec(catQuery); err != nil {
		return nil, err
	}

	// Создаем таблицу трат со связью (FOREIGN KEY)
	expQuery := `
	CREATE TABLE IF NOT EXISTS expenses (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		amount REAL NOT NULL,
		description TEXT NOT NULL,
		date TEXT NOT NULL,
		category_id INTEGER,
		FOREIGN KEY (category_id) REFERENCES categories(id)
	);`

	if _, err = db.Exec(expQuery); err != nil {
		return nil, err
	}

	return &ExpenseRepository{db: db}, nil
}

// Получить все категории для выпадающих списков (<select>)
func (r *ExpenseRepository) GetCategories() ([]models.Category, error) {
	rows, err := r.db.Query("SELECT id, name FROM categories ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

// Получить траты с использованием JOIN и динамических фильтров (WHERE)
func (r *ExpenseRepository) GetFiltered(categoryID string, dateFrom string, dateTo string) ([]models.Expense, error) {
	query := `
		SELECT e.id, e.amount, e.description, e.date, COALESCE(e.category_id, 0), COALESCE(c.name, 'Без категории')
		FROM expenses e
		LEFT JOIN categories c ON e.category_id = c.id
		WHERE 1=1`
	
	var args []interface{}

	// Динамически добавляем условия, если фильтры переданы
	if categoryID != "" {
		query += " AND e.category_id = ?"
		args = append(args, categoryID)
	}
	if dateFrom != "" {
		query += " AND e.date >= ?"
		args = append(args, dateFrom)
	}
	if dateTo != "" {
		query += " AND e.date <= ?"
		args = append(args, dateTo)
	}

	query += " ORDER BY e.date DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []models.Expense
	for rows.Next() {
		var exp models.Expense
		var dateStr string
		err := rows.Scan(&exp.ID, &exp.Amount, &exp.Description, &dateStr, &exp.CategoryID, &exp.CategoryName)
		if err != nil {
			return nil, err
		}
		exp.Date, _ = time.Parse("2006-01-02", dateStr)
		expenses = append(expenses, exp)
	}
	return expenses, nil
}

// Создать трату с указанием категории
func (r *ExpenseRepository) Create(amount float64, desc string, date time.Time, categoryID int) error {
	query := "INSERT INTO expenses (amount, description, date, category_id) VALUES (?, ?, ?, ?)"
	_, err := r.db.Exec(query, amount, desc, date.Format("2006-01-02"), categoryID)
	return err
}

// Получить трату по ID для редактирования
func (r *ExpenseRepository) GetByID(id int) (*models.Expense, error) {
	query := "SELECT id, amount, description, date, COALESCE(category_id, 0) FROM expenses WHERE id = ?"
	row := r.db.QueryRow(query, id)

	var exp models.Expense
	var dateStr string
	err := row.Scan(&exp.ID, &exp.Amount, &exp.Description, &dateStr, &exp.CategoryID)
	if err != nil {
		return nil, err
	}
	exp.Date, _ = time.Parse("2006-01-02", dateStr)
	return &exp, nil
}

// Обновить трату вместе с категорией
func (r *ExpenseRepository) Update(id int, amount float64, desc string, date time.Time, categoryID int) error {
	query := "UPDATE expenses SET amount = ?, description = ?, date = ?, category_id = ? WHERE id = ?"
	_, err := r.db.Exec(query, amount, desc, date.Format("2006-01-02"), categoryID, id)
	return err
}

// Удалить трату
func (r *ExpenseRepository) Delete(id int) error {
	query := "DELETE FROM expenses WHERE id = ?"
	_, err := r.db.Exec(query, id)
	return err
}
