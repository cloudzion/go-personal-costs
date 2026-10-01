package models

import "time"

// Expense описывает структуру одной траты
type Expense struct {
	Amount      float64   // Сумма
	Description string    // Описание
	Date        time.Time // Дата
}
