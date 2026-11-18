package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"

	"biblioteca-api/internal/models"
	"biblioteca-api/internal/service"
)

type LaptopHandler struct {
	laptopService *service.LaptopService
}

func NewLaptopHandler(laptopService *service.LaptopService) *LaptopHandler {
	return &LaptopHandler{laptopService: laptopService}
}

// List godoc
// @Summary      Listar portátiles
// @Tags         laptops
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} models.Laptop
// @Router       /laptops [get]
func (h *LaptopHandler) List(c *gin.Context) {
	laptops, err := h.laptopService.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo obtener el catálogo de portátiles"})
		return
	}
	c.JSON(http.StatusOK, laptops)
}

// Create godoc
// @Summary      Registrar un portátil (solo admin)
// @Tags         laptops
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body models.LaptopRequest true "Datos del portátil"
// @Success      201 {object} models.Laptop
// @Router       /laptops [post]
func (h *LaptopHandler) Create(c *gin.Context) {
	var req models.LaptopRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Brand = strings.TrimSpace(req.Brand)
	req.Model = strings.TrimSpace(req.Model)
	req.SerialNumber = strings.TrimSpace(req.SerialNumber)
	if req.Brand == "" || req.Model == "" || req.SerialNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "marca, modelo y serial son obligatorios"})
		return
	}

	laptop, err := h.laptopService.Create(req)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			c.JSON(http.StatusConflict, gin.H{"error": "ya existe un portátil con ese serial"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo registrar el portátil"})
		return
	}
	c.JSON(http.StatusCreated, laptop)
}

// Borrow godoc
// @Summary      Solicitar un portátil
// @Tags         laptop-loans
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body models.LaptopLoanRequest true "ID del portátil"
// @Success      201 {object} models.LaptopLoan
// @Router       /laptop-loans [post]
func (h *LaptopHandler) Borrow(c *gin.Context) {
	var req models.LaptopLoanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	loan, err := h.laptopService.Borrow(c.GetInt("user_id"), req)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "portátil no encontrado"})
		return
	}
	if errors.Is(err, service.ErrLaptopUnavailable) {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo registrar el préstamo"})
		return
	}
	c.JSON(http.StatusCreated, loan)
}

// MyLoans godoc
// @Summary      Mis préstamos de portátiles
// @Tags         laptop-loans
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} models.LaptopLoanDetail
// @Router       /laptop-loans/me [get]
func (h *LaptopHandler) MyLoans(c *gin.Context) {
	loans, err := h.laptopService.ListLoansByUser(c.GetInt("user_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudieron obtener los préstamos"})
		return
	}
	c.JSON(http.StatusOK, loans)
}

// Return godoc
// @Summary      Devolver un portátil
// @Tags         laptop-loans
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID del préstamo"
// @Success      200 {object} map[string]string
// @Router       /laptop-loans/{id}/return [put]
func (h *LaptopHandler) Return(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	err = h.laptopService.Return(id, c.GetInt("user_id"))
	switch {
	case errors.Is(err, service.ErrLaptopLoanNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, service.ErrLaptopLoanNotOwner):
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	case errors.Is(err, service.ErrLaptopLoanReturned):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo devolver el portátil"})
	default:
		c.JSON(http.StatusOK, gin.H{"message": "portátil devuelto correctamente"})
	}
}
