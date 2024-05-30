package handlers

import (
	"net/http"

	"backend/internal/config"
	"backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateTransactionInput struct {
	UserID           uuid.UUID `json:"user_id"`
	PaymentDetailsID uuid.UUID `json:"payment_details_id"`
	Status           string    `json:"status"`
}

type UpdateTransactionStatusInput struct {
	Status string `json:"status"`
}

// CreateTransaction godoc
// @Summary Create a new transaction
// @Description Create a new transaction
// @Tags transactions
// @Accept json
// @Produce json
// @Param transaction body CreateTransactionInput true "Transaction Input"
// @Success 201 {object} models.Transaction
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /transactions [post]
func CreateTransaction(c *gin.Context) {
	var input CreateTransactionInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction := models.Transaction{
		UserID:           input.UserID,
		PaymentDetailsID: input.PaymentDetailsID,
		Status:           input.Status,
	}
	result := config.DB.Create(&transaction)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, transaction)
}

// UpdateTransactionStatus godoc
// @Summary Update transaction status by ID
// @Description Update transaction status by ID
// @Tags transactions
// @Accept json
// @Produce json
// @Param id path string true "Transaction ID"
//
//	@Param status body UpdateTransactionStatusInput true "Status Update Input"
//
// @Success 200 {object} models.Transaction
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /transactions/{id} [put]
func UpdateTransactionStatus(c *gin.Context) {
	id := c.Param("id")
	transactionID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction ID format"})
		return
	}

	var input UpdateTransactionStatusInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var transaction models.Transaction
	result := config.DB.First(&transaction, "id = ?", transactionID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		}
		return
	}

	result = config.DB.Model(&transaction).Update("status", input.Status)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, transaction)
}
