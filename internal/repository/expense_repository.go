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

	initSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS sessions (
		session_id TEXT PRIMARY KEY,
		user_id INTEGER NOT NULL,
		expires_at DATETIME NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);
	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE
	);
	INSERT OR IGNORE INTO categories (name) VALUES ('Еда'), ('Транспорт'), ('Жильё'), ('Развлечения'), ('Прочее');
	CREATE TABLE IF NOT EXISTS expenses (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		amount REAL NOT NULL,
		description TEXT NOT NULL,
		date TEXT NOT NULL,
		category_id INTEGER,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (category_id) REFERENCES categories(id)
	);`

	if _, err = db.Exec(initSQL); err != nil {
		return nil, err
	}

	return &ExpenseRepository{db: db}, nil
}

// --- Методы для пользователей и сессий ---

func (r *ExpenseRepository) CreateUser(email, passwordHash string) error {
	_, err := r.db.Exec("INSERT INTO users (email, password_hash) VALUES (?, ?)", email, passwordHash)
	return err
}

func (r *ExpenseRepository) GetUserByEmail(email string) (*models.User, error) {
	row := r.db.QueryRow("SELECT id, email, password_hash FROM users WHERE email = ?", email)
	var u models.User
	err := row.Scan(&u.ID, &u.Email, &u.PasswordHash)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *ExpenseRepository) CreateSession(sessionID string, userID int, expiresAt time.Time) error {
	_, err := r.db.Exec("INSERT INTO sessions (session_id, user_id, expires_at) VALUES (?, ?, ?)", sessionID, userID, expiresAt)
	return err
}

func (r *ExpenseRepository) GetUserIDBySession(sessionID string) (int, error) {
	row := r.db.QueryRow("SELECT user_id, expires_at FROM sessions WHERE session_id = ?", sessionID)
	var userID int
	var expiresAt time.Time
	if err := row.Scan(&userID, &expiresAt); err != nil {
		return 0, err
	}
	if time.Now().After(expiresAt) {
		r.DeleteSession(sessionID)
		return 0, sql.ErrNoRows
	}
	return userID, nil
}

func (r *ExpenseRepository) DeleteSession(sessionID string) error {
	_, err := r.db.Exec("DELETE FROM sessions WHERE session_id = ?", sessionID)
	return err
}

// --- Методы для трат (строго с привязкой к userID) ---

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

func (r *ExpenseRepository) GetFiltered(userID int, categoryID, dateFrom, dateTo string) ([]models.Expense, error) {
	query := `
		SELECT e.id, e.amount, e.description, e.date, COALESCE(e.category_id, 0), COALESCE(c.name, 'Без категории')
		FROM expenses e
		LEFT JOIN categories c ON e.category_id = c.id
		WHERE e.user_id = ?`

	args := []interface{}{userID}

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

func (r *ExpenseRepository) Create(userID int, amount float64, desc string, date time.Time, categoryID int) error {
	query := "INSERT INTO expenses (user_id, amount, description, date, category_id) VALUES (?, ?, ?, ?, ?)"
	_, err := r.db.Exec(query, userID, amount, desc, date.Format("2006-01-02"), categoryID)
	return err
}

func (r *ExpenseRepository) GetByID(id, userID int) (*models.Expense, error) {
	query := "SELECT id, amount, description, date, COALESCE(category_id, 0) FROM expenses WHERE id = ? AND user_id = ?"
	row := r.db.QueryRow(query, id, userID)

	var exp models.Expense
	var dateStr string
	err := row.Scan(&exp.ID, &exp.Amount, &exp.Description, &dateStr, &exp.CategoryID)
	if err != nil {
		return nil, err
	}
	exp.Date, _ = time.Parse("2006-01-02", dateStr)
	return &exp, nil
}

func (r *ExpenseRepository) Update(id, userID int, amount float64, desc string, date time.Time, categoryID int) error {
	query := "UPDATE expenses SET amount = ?, description = ?, date = ?, category_id = ? WHERE id = ? AND user_id = ?"
	_, err := r.db.Exec(query, amount, desc, date.Format("2006-01-02"), categoryID, id, userID)
	return err
}

func (r *ExpenseRepository) Delete(id, userID int) error {
	query := "DELETE FROM expenses WHERE id = ? AND user_id = ?"
	_, err := r.db.Exec(query, id, userID)
	return err
}
