package handlers

import (
	_ "backend/docs"
	"backend/internal/config"
	"backend/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GetGameByID godoc
// @Summary Get a game by ID
// @Description Get a game by ID
// @Tags games
// @Param id path string true "Game ID"
// @Success 200 {object} models.Game
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /games/{id} [get]
func GetGameByID(c *gin.Context) {
	id := c.Param("id")
	var game models.Game
	result := config.DB.First(&game, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, game)
}

// CreateGame godoc
// @Summary Create a new game
// @Description Create a new game
// @Tags games
// @Success 201 {object} models.Game
// @Failure 500 {object} models.ErrorResponse
// @Router /games [post]
func CreateGame(c *gin.Context) {
	game := models.Game{
		ID:   uuid.New(),
		Date: time.Now(),
	}

	result := config.DB.Create(&game)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, game)
}

// DeleteGame godoc
// @Summary Delete a game by ID
// @Description Delete a game by ID and all related played games
// @Tags games
// @Param id path string true "Game ID"
// @Success 204
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /games/{id} [delete]
func DeleteGame(c *gin.Context) {
	id := c.Param("id")

	tx := config.DB.Begin()
	if tx.Error != nil {
		respondWithError(c, http.StatusInternalServerError, tx.Error.Error())
		return
	}

	if !deleteWithTransaction(tx, &models.PlayedGames{}, "game_id = ?", id, c) {
		return
	}

	if !deleteWithTransaction(tx, &models.Game{}, "id = ?", id, c) {
		return
	}

	if err := tx.Commit().Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func deleteWithTransaction(tx *gorm.DB, model interface{}, query string, args interface{}, c *gin.Context) bool {
	result := tx.Delete(model, query, args)
	if result.Error != nil {
		tx.Rollback()
		respondWithError(c, http.StatusInternalServerError, result.Error.Error())
		return false
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		respondWithError(c, http.StatusNotFound, "Resource not found")
		return false
	}

	return true
}

func respondWithError(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"error": message})
}
