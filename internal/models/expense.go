package models

import "time"

type Category struct {
	ID   int
	Name string
}

type Expense struct {
	ID           int
	Amount       float64
	Description  string
	Date         time.Time
	CategoryID   int
	CategoryName string
}
