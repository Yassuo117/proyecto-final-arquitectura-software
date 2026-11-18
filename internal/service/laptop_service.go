package service

import (
	"time"

	"biblioteca-api/internal/models"
	"biblioteca-api/internal/repository"
)

const laptopLoanDurationDays = 7

// Errores de negocio expuestos por la capa de servicio.
var (
	ErrLaptopUnavailable  = repository.ErrLaptopUnavailable
	ErrLaptopLoanNotFound = repository.ErrLaptopLoanNotFound
	ErrLaptopLoanReturned = repository.ErrLaptopLoanReturned
	ErrLaptopLoanNotOwner = repository.ErrLaptopLoanNotOwner
)

type LaptopService struct {
	laptopRepo     *repository.LaptopRepository
	laptopLoanRepo *repository.LaptopLoanRepository
}

func NewLaptopService(
	laptopRepo *repository.LaptopRepository,
	laptopLoanRepo *repository.LaptopLoanRepository,
) *LaptopService {
	return &LaptopService{laptopRepo: laptopRepo, laptopLoanRepo: laptopLoanRepo}
}

func (s *LaptopService) List() ([]models.Laptop, error) {
	return s.laptopRepo.FindAll()
}

func (s *LaptopService) Create(req models.LaptopRequest) (*models.Laptop, error) {
	laptop := &models.Laptop{
		Brand:        req.Brand,
		Model:        req.Model,
		SerialNumber: req.SerialNumber,
	}
	if err := s.laptopRepo.Create(laptop); err != nil {
		return nil, err
	}
	return laptop, nil
}

func (s *LaptopService) Borrow(userID int, req models.LaptopLoanRequest) (*models.LaptopLoan, error) {
	dueAt := time.Now().AddDate(0, 0, laptopLoanDurationDays)
	return s.laptopLoanRepo.Borrow(userID, req.LaptopID, dueAt)
}

func (s *LaptopService) Return(loanID, userID int) error {
	return s.laptopLoanRepo.Return(loanID, userID)
}

func (s *LaptopService) ListLoansByUser(userID int) ([]models.LaptopLoanDetail, error) {
	return s.laptopLoanRepo.FindByUser(userID)
}
