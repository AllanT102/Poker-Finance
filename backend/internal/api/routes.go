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
		
		// Transaction routes
		v1.GET("/users/:id/transactions", handlers.GetUserTransactions)
		v1.POST("/transactions", handlers.CreateTransaction)
		v1.PUT("/transactions/:id", handlers.UpdateTransactionStatus)
	
		// Payment Details routes
		v1.GET("/payment-details/:id", handlers.GetPaymentDetailsByID)
		v1.POST("/payment-details", handlers.CreatePaymentDetails)
		v1.PUT("/payment-details/:id", handlers.UpdatePaymentDetails)
	}
}
