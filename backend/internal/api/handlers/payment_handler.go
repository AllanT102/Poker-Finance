package handlers

import (
	"backend/internal/config"
	"backend/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreatePaymentDetailsInput struct {
	PayerID uuid.UUID `json:"payer_id" binding:"required"`
	PayeeID uuid.UUID `json:"payee_id" binding:"required"`
	Amount  float64   `json:"amount" binding:"required"`
	Details string    `json:"details"`
	Status  string    `json:"status" binding:"required"`
}

type UpdatePaymentDetailsInput struct {
	PayerID       uuid.UUID `json:"payer_id"`
	PayeeID       uuid.UUID `json:"payee_id"`
	Amount        float64   `json:"amount"`
	Details       string    `json:"details"`
	Status        string    `json:"status"`
	TimeCompleted time.Time `json:"time_completed"`
}

// GetPaymentDetailsByID godoc
// @Summary Get payment details by ID
// @Description Get payment details by ID
// @Tags payment-details
// @Param id path string true "Payment Details ID"
// @Success 200 {object} models.PaymentDetails
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /payment-details/{id} [get]
func GetPaymentDetailsByID(c *gin.Context) {
	id := c.Param("id")
	var paymentDetails models.PaymentDetails
	result := config.DB.First(&paymentDetails, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Payment details not found"})
		} else {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, paymentDetails)
}

// CreatePaymentDetails godoc
// @Summary Create new payment details
// @Description Create new payment details
// @Tags payment-details
// @Accept json
// @Produce json
// @Param paymentDetails body CreatePaymentDetailsInput true "Payment Details Input"
// @Success 201 {object} models.PaymentDetails
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /payment-details [post]
func CreatePaymentDetails(c *gin.Context) {
	var input CreatePaymentDetailsInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	paymentDetails := models.PaymentDetails{
		ID:            uuid.New(),
		PayerID:       input.PayerID,
		PayeeID:       input.PayeeID,
		Amount:        input.Amount,
		Details:       input.Details,
		TimeSubmitted: time.Now(),
		Status:        input.Status,
	}

	result := config.DB.Create(&paymentDetails)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, paymentDetails)
}

// UpdatePaymentDetails godoc
// @Summary Update payment details by ID
// @Description Update payment details by ID
// @Tags payment-details
// @Accept json
// @Produce json
// @Param id path string true "Payment Details ID"
// @Param paymentDetails body UpdatePaymentDetailsInput true "Payment Details Input"
// @Success 200 {object} models.PaymentDetails
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /payment-details/{id} [put]
func UpdatePaymentDetails(c *gin.Context) {
	id := c.Param("id")
	var paymentDetails models.PaymentDetails
	result := config.DB.First(&paymentDetails, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Payment details not found"})
		} else {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
		}
		return
	}

	var input UpdatePaymentDetailsInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	result = config.DB.Model(&paymentDetails).Updates(models.PaymentDetails{
		PayerID:       input.PayerID,
		PayeeID:       input.PayeeID,
		Amount:        input.Amount,
		Details:       input.Details,
		Status:        input.Status,
		TimeCompleted: input.TimeCompleted,
	})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, paymentDetails)
}
