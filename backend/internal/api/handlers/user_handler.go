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

type CreateUserInput struct {
	Email       string `json:"email" binding:"required"`
	Name        string `json:"name" binding:"required"`
	DisplayName string `json:"display_name"`
	Balance     int    `json:"balance" binding:"required"`
}

type UpdateUserInput struct {
	Email       string `json:"email"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Balance     int    `json:"balance"`
}

type UpdateUserPlayedGameInput struct {
	BuyIn     int `json:"buy_in"`
	EndAmount int `json:"end_amount"`
}

type CreateUserPlayedGameInput struct {
	GameID uuid.UUID `json:"game_id" binding:"required"`
	UserID uuid.UUID `json:"user_id" binding:"required"`
	BuyIn  int       `json:"buy_in" binding:"required"`
}

// GetUserByID godoc
// @Summary Get a user by ID
// @Description Get a user by ID
// @Tags users
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users/{id} [get]
func GetUserByID(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	result := config.DB.First(&user, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, user)
}

// GetUserByEmail godoc
// @Summary Get a user by email
// @Description Get a user by email
// @Tags users
// @Param id path string true "User email"
// @Success 200 {object} models.User
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users/{id} [get]
func GetUserByEmail(c *gin.Context) {
	email := c.Param("email")
	var user models.User
	result := config.DB.First(&user, "email = ?", email)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, user)
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateUserInput true "User Input"
// @Success 201 {object} models.User
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users [post]
func CreateUser(c *gin.Context) {
	var input CreateUserInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		ID:          uuid.New(),
		Email:       input.Email,
		Name:        input.Name,
		DisplayName: input.DisplayName,
		Balance:     input.Balance,
		CreatedAt:   time.Now(),
	}

	result := config.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// UpdateUser godoc
// @Summary Update user information
// @Description Update user information
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body UpdateUserInput true "User Update Input"
// @Success 200 {object} models.User
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users/{id} [put]
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	result := config.DB.First(&user, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		}
		return
	}

	var input UpdateUserInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedData := models.User{
		Email:       input.Email,
		Name:        input.Name,
		DisplayName: input.DisplayName,
		Balance:     input.Balance,
	}

	result = config.DB.Model(&user).Updates(updatedData)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetUserPlayedGames godoc
// @Summary Get played games by user ID
// @Description Get played games by user ID
// @Tags users
// @Param id path string true "User ID"
// @Success 200 {array} models.PlayedGames
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users/{id}/played-games [get]
func GetUserPlayedGames(c *gin.Context) {
	id := c.Param("id")
	var playedGames []models.PlayedGames
	result := config.DB.Find(&playedGames, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Played games not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, playedGames)
}

// CreateUserPlayedGame godoc
// @Summary Create a played game for a user
// @Description Create a new played game entry for a user
// @Tags users
// @Accept json
// @Produce json
// @Param playedGame body CreateUserPlayedGameInput true "Played Game Input"
// @Success 201 {object} models.PlayedGames
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users/{id}/played-games/{gameid} [post]
func CreateUserPlayedGame(c *gin.Context) {
	var input CreateUserPlayedGameInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	playedGame := models.PlayedGames{
		GameID:    input.GameID,
		PlayerID:  input.UserID,
		BuyIn:     input.BuyIn,
		EndAmount: 0,
	}

	result := config.DB.Create(&playedGame)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, playedGame)
}

// UpdateUserPlayedGame godoc
// @Summary Update a played game by user ID and game ID
// @Description Update a played game's details by user ID and game ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param gameid path string true "Game ID"
// @Param playedGame body UpdateUserPlayedGameInput true "Played Game Update Input"
// @Success 200 {object} models.PlayedGames
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users/{id}/played-games/{gameid} [put]
func UpdateUserPlayedGame(c *gin.Context) {
	userID := c.Param("id")
	gameID := c.Param("gameid")

	var playedGame models.PlayedGames
	result := config.DB.First(&playedGame, "game_id = ? AND player_id = ?", gameID, userID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Payment details not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		}
		return
	}

	var input UpdateUserPlayedGameInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedData := models.PlayedGames{
		BuyIn:     input.BuyIn,
		EndAmount: input.EndAmount,
	}

	result = config.DB.Model(&playedGame).Updates(updatedData)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, playedGame)
}

// GetUserTransactions godoc
// @Summary Get transactions by user ID
// @Description Get all transactions by user ID
// @Tags users
// @Param id path string true "User ID"
// @Success 200 {array} models.Transaction
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users/{id}/transactions [get]
func GetUserTransactions(c *gin.Context) {
	id := c.Param("id")
	var transactions []models.Transaction
	result := config.DB.Find(&transactions, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Played games not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, transactions)
}

// GetUserFriends godoc
// @Summary Get friends by user ID
// @Description Get all friends of a user by user ID
// @Tags users
// @Param id path string true "User ID"
// @Success 200 {array} models.FriendRequest
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users/{id}/friends [get]
func GetUserFriends(c *gin.Context) {
	id := c.Param("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid User ID format"})
		return
	}

	var friends []models.FriendRequest
	result := config.DB.Where("user_id = ? OR friend_id = ?", userID, userID).Find(&friends)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "No friends found"})
		return
	}

	c.JSON(http.StatusOK, friends)
}
