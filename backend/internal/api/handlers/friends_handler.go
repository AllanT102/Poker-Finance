package handlers

import (
	"net/http"
	"time"

	"backend/internal/config"
	"backend/internal/models"
	"backend/internal/constants"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateFriendRequestInput struct {
	UserID      uuid.UUID `json:"user_id" binding:"required"`
	FriendEmail string    `json:"friend_email" binding:"required"`
}

type UpdateFriendRequestInput struct {
	Status string `json:"status" binding:"required"` // e.g., "pending", "accepted", "declined"
}

// CreateFriendRequest godoc
// @Summary Create a new friend request
// @Description Create a new friend request
// @Tags friend-request
// @Accept json
// @Produce json
// @Param FriendRequest body CreateFriendRequestInput true "Friend Request Input"
// @Success 201 {object} models.FriendRequest
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /friend-request [post]
func CreateFriendRequest(c *gin.Context) {
	var input CreateFriendRequestInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	var user models.User
	result := config.DB.First(&user, "email = ?", input.FriendEmail)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "user not found"})
		} else {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
		}
		return
	}

	friendRequest := models.FriendRequest{
		ID:        uuid.New(),
		UserID:    input.UserID,
		FriendID:  user.ID,
		Status:    constants.FRIEND_REQUEST_PENDING,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result = config.DB.Create(&friendRequest)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, friendRequest)
}

// UpdateFriendRequest godoc
// @Summary Update a friend request by ID
// @Description Update a FriendRequest's status by ID
// @Tags friend-request
// @Accept json
// @Produce json
// @Param id path string true "friend request ID"
// @Param FriendRequest body UpdateFriendRequestInput true "FriendRequest Update Input"
// @Success 200 {object} models.FriendRequest
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /FriendRequests/{id} [put]
func UpdateFriendRequest(c *gin.Context) {
	id := c.Param("id")
	FriendRequestID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid Friend Request ID format"})
		return
	}

	var input UpdateFriendRequestInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	var friendRequest models.FriendRequest
	result := config.DB.First(&friendRequest, "id = ?", FriendRequestID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Friend Request not found"})
		} else {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
		}
		return
	}

	friendRequest.Status = input.Status
	friendRequest.UpdatedAt = time.Now()

	result = config.DB.Save(&friendRequest)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, friendRequest)
}
