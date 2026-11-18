package models

import "time"

// Laptop representa un equipo físico identificado por un serial único.
type Laptop struct {
	ID           int       `db:"id" json:"id"`
	Brand        string    `db:"brand" json:"brand"`
	Model        string    `db:"model" json:"model"`
	SerialNumber string    `db:"serial_number" json:"serial_number"`
	Status       string    `db:"status" json:"status"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

// LaptopRequest contiene los datos necesarios para registrar un portátil.
type LaptopRequest struct {
	Brand        string `json:"brand" binding:"required" example:"Lenovo"`
	Model        string `json:"model" binding:"required" example:"ThinkPad E14"`
	SerialNumber string `json:"serial_number" binding:"required" example:"UTB-LAP-001"`
}

// LaptopLoan representa el préstamo de un portátil.
type LaptopLoan struct {
	ID         int        `db:"id" json:"id"`
	UserID     int        `db:"user_id" json:"user_id"`
	LaptopID   int        `db:"laptop_id" json:"laptop_id"`
	LoanedAt   time.Time  `db:"loaned_at" json:"loaned_at"`
	DueAt      time.Time  `db:"due_at" json:"due_at"`
	ReturnedAt *time.Time `db:"returned_at" json:"returned_at,omitempty"`
	Status     string     `db:"status" json:"status"`
}

// LaptopLoanRequest es el payload para solicitar un portátil.
type LaptopLoanRequest struct {
	LaptopID int `json:"laptop_id" binding:"required,min=1" example:"1"`
}

// LaptopLoanDetail incluye los datos básicos del equipo prestado.
type LaptopLoanDetail struct {
	LaptopLoan
	LaptopBrand        string `db:"laptop_brand" json:"laptop_brand"`
	LaptopModel        string `db:"laptop_model" json:"laptop_model"`
	LaptopSerialNumber string `db:"laptop_serial_number" json:"laptop_serial_number"`
}
