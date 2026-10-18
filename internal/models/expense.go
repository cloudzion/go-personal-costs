package models

import "time"

type Expense struct {
	ID          int       // Добавили ID для CRUD операций
	Amount      float64   
	Description string    
	Date        time.Time 
}
