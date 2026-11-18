package repository

import (
	"github.com/jmoiron/sqlx"

	"biblioteca-api/internal/models"
)

type LaptopRepository struct {
	db *sqlx.DB
}

func NewLaptopRepository(db *sqlx.DB) *LaptopRepository {
	return &LaptopRepository{db: db}
}

func (r *LaptopRepository) FindAll() ([]models.Laptop, error) {
	laptops := make([]models.Laptop, 0)
	if err := r.db.Select(&laptops, `SELECT * FROM laptops ORDER BY brand, model, id`); err != nil {
		return nil, err
	}
	return laptops, nil
}

func (r *LaptopRepository) FindByID(id int) (*models.Laptop, error) {
	var laptop models.Laptop
	if err := r.db.Get(&laptop, `SELECT * FROM laptops WHERE id = ?`, id); err != nil {
		return nil, err
	}
	return &laptop, nil
}

func (r *LaptopRepository) Create(laptop *models.Laptop) error {
	result, err := r.db.Exec(
		`INSERT INTO laptops (brand, model, serial_number, status) VALUES (?, ?, ?, 'available')`,
		laptop.Brand, laptop.Model, laptop.SerialNumber,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	laptop.ID = int(id)
	created, err := r.FindByID(laptop.ID)
	if err != nil {
		return err
	}
	*laptop = *created
	return nil
}
