package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"

	"biblioteca-api/internal/models"
)

var (
	ErrLaptopUnavailable  = errors.New("el portátil no está disponible")
	ErrLaptopLoanNotFound = errors.New("préstamo de portátil no encontrado")
	ErrLaptopLoanReturned = errors.New("el portátil ya fue devuelto")
	ErrLaptopLoanNotOwner = errors.New("el préstamo pertenece a otro usuario")
)

type LaptopLoanRepository struct {
	db *sqlx.DB
}

func NewLaptopLoanRepository(db *sqlx.DB) *LaptopLoanRepository {
	return &LaptopLoanRepository{db: db}
}

// Borrow bloquea el equipo durante la operación para evitar dos préstamos simultáneos.
func (r *LaptopLoanRepository) Borrow(userID, laptopID int, dueAt time.Time) (*models.LaptopLoan, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var status string
	if err := tx.Get(&status, `SELECT status FROM laptops WHERE id = ? FOR UPDATE`, laptopID); err != nil {
		return nil, err
	}
	if status != "available" {
		return nil, ErrLaptopUnavailable
	}

	result, err := tx.Exec(
		`INSERT INTO laptop_loans (user_id, laptop_id, due_at, status) VALUES (?, ?, ?, 'active')`,
		userID, laptopID, dueAt,
	)
	if err != nil {
		return nil, err
	}
	if _, err := tx.Exec(`UPDATE laptops SET status = 'loaned' WHERE id = ?`, laptopID); err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return r.FindByID(int(id))
}

func (r *LaptopLoanRepository) FindByID(id int) (*models.LaptopLoan, error) {
	var loan models.LaptopLoan
	if err := r.db.Get(&loan, `SELECT * FROM laptop_loans WHERE id = ?`, id); err != nil {
		return nil, err
	}
	return &loan, nil
}

func (r *LaptopLoanRepository) FindByUser(userID int) ([]models.LaptopLoanDetail, error) {
	query := `
		SELECT ll.*, l.brand AS laptop_brand, l.model AS laptop_model,
		       l.serial_number AS laptop_serial_number
		FROM laptop_loans ll
		JOIN laptops l ON l.id = ll.laptop_id
		WHERE ll.user_id = ?
		ORDER BY ll.loaned_at DESC`

	loans := make([]models.LaptopLoanDetail, 0)
	if err := r.db.Select(&loans, query, userID); err != nil {
		return nil, err
	}
	return loans, nil
}

// Return valida propietario y estado dentro de la misma transacción que libera el equipo.
func (r *LaptopLoanRepository) Return(loanID, userID int) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var loan models.LaptopLoan
	if err := tx.Get(&loan, `SELECT * FROM laptop_loans WHERE id = ? FOR UPDATE`, loanID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrLaptopLoanNotFound
		}
		return err
	}
	if loan.UserID != userID {
		return ErrLaptopLoanNotOwner
	}
	if loan.Status == "returned" {
		return ErrLaptopLoanReturned
	}

	if _, err := tx.Exec(
		`UPDATE laptop_loans SET status = 'returned', returned_at = NOW() WHERE id = ?`, loanID,
	); err != nil {
		return err
	}
	if _, err := tx.Exec(`UPDATE laptops SET status = 'available' WHERE id = ?`, loan.LaptopID); err != nil {
		return err
	}
	return tx.Commit()
}
