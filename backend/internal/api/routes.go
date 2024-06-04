package api

import (
	"backend/internal/api/handlers"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		// swagger
		v1.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		// User routes
		v1.GET("/users/:id", handlers.GetUserByID)
		v1.GET("/users/email/:email", handlers.GetUserByEmail)
		v1.POST("/users", handlers.CreateUser)
		v1.PUT("/users/:id", handlers.UpdateUser)

		// Played games routes
		v1.GET("/users/:id/played-games", handlers.GetUserPlayedGames)
		v1.POST("/users/:id/played-games/:gameid", handlers.CreateUserPlayedGame)
		v1.PUT("/users/:id/played-games/:gameid", handlers.UpdateUserPlayedGame)

		// Game routes
		v1.GET("/games/:id", handlers.GetGameByID)
		v1.POST("/games", handlers.CreateGame)
		v1.DELETE("/games/:id", handlers.DeleteGame)
		v1.POST("/games/:id/end", handlers.EndGame)

		// // Transaction routes
		// v1.GET("/users/:id/transactions", handlers.GetUserTransactions)
		// v1.POST("/transactions", handlers.CreateTransaction)
		// v1.PUT("/transactions/:id", handlers.UpdateTransactionStatus)

		// Payment Details routes
		v1.GET("/payment-details/:id", handlers.GetPaymentDetailsByID)
		v1.POST("/payment-details", handlers.CreatePaymentDetails)
		v1.PUT("/payment-details/:id", handlers.UpdatePaymentDetails)

		// friends
		v1.POST("/friend-request", handlers.CreateFriendRequest)
		v1.PUT("/friend-request/:id", handlers.UpdateFriendRequest)
		v1.GET("/users/:id/friends", handlers.GetUserFriends)
	}
}

// 4.2 seconds blocking api calls
// goroutine to 3.6