package handlers

import (
	_ "backend/docs"
	"backend/internal/config"
	"backend/internal/models"

	"backend/internal/constants"
	"errors"
	"fmt"
	"math"
	"net/http"
	"sort"
	"time"

	"backend/internal/services/email"

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

// EndGame ends a game by setting its status to 'ended'.
// @Summary End a game
// @Description Ends the game with the specified ID by updating its status.
// @Tags game
// @Accept json
// @Produce json
// @Param id path string true "Game ID"
// @Success 200 {object} models.Game "Game ended successfully"
// @Failure 400 {object} models.ErrorResponse "Bad Request if the request body is incorrect"
// @Failure 404 {object} gin.H "Not Found if the game does not exist"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error for database issues"
// @Router /game/{id}/end [post]
func EndGame(c *gin.Context) {
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

	// if cant calculate, then someone hasn't updated their shit properly
	err := payoutPlayers(uuid.MustParse(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
	}

	result = config.DB.Model(&game).Updates(models.Game{
		Status: "ended",
	})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, nil)
}

func payoutPlayers(gameID uuid.UUID) error {
	var playedGames []models.PlayedGames

	result := config.DB.Find(&playedGames, "game_id = ?", gameID)
	if result.Error != nil {
		return result.Error
	}

	idToPay := make(map[uuid.UUID]float64)
	totalLosses := 0.0
	totalWins := 0.0
	for _, pg := range playedGames {
		if pg.EndAmount < pg.BuyIn {
			totalLosses += pg.EndAmount - pg.BuyIn
		} else {
			totalWins += pg.EndAmount - pg.BuyIn
		}
		idToPay[pg.PlayerID] = pg.EndAmount - pg.BuyIn
	}
	if math.Abs(totalLosses) != totalWins {
		fmt.Printf("totalLoss %.2f, totalWins: %.2f", totalLosses, totalWins)
		return errors.New("total winnings should equal total losses - something doesn't add up")
	}

	type kv struct {
		Key   uuid.UUID
		Value float64
	}

	var ss []kv
	for k, v := range idToPay {
		ss = append(ss, kv{k, v})
	}
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value < ss[j].Value
	})

	left, right := 0, len(ss)-1
	for {
		if left >= right {
			break
		}

		amtToPay := math.Min(ss[right].Value, math.Abs(ss[left].Value))

		details := *models.NewPaymentDetails(
			gameID, ss[left].Key, ss[right].Key, amtToPay, "they lost lol", constants.PAYMENT_DETAILS_PENDING,
		)

		if ss[right].Value < math.Abs(ss[left].Value) {
			ss[left].Value += ss[right].Value
			right -= 1
		} else {
			ss[right].Value += ss[left].Value
			left += 1
		}
		if amtToPay == 0 {
			continue
		}

		savePaymentDetails(&details)
		email.QueueEmail(details)
	}

	return nil
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
		ID:     uuid.New(),
		Date:   time.Now(),
		Status: "started",
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

func savePaymentDetails(details *models.PaymentDetails) error {
	result := config.DB.Create(&details)
	return result.Error
}
