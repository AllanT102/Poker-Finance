package main

import (
	"backend/internal/api"
	"backend/internal/config"
	"backend/internal/services/email"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
)

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

func perClientRateLimiter() gin.HandlerFunc {
	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)
	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		mu.Lock()
		if _, found := clients[ip]; !found {
			clients[ip] = &client{limiter: rate.NewLimiter(2, 4)}
		}
		clients[ip].lastSeen = time.Now()
		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			c.JSON(http.StatusTooManyRequests, gin.H{
				"status": "Request Failed",
				"body":   "The API is at capacity, try again later.",
			})
			c.Abort()
			return
		}
		mu.Unlock()
		c.Next()
	}
}

func validFloat(fl validator.FieldLevel) bool {
	float := fl.Field().Float()
	fmt.Println(float >= 0.0)
	return true
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	config.InitDatabase()
	email.CreateEmailChannel()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("valid_float", validFloat)
	}

	router := gin.Default()
	router.Use(perClientRateLimiter())
	api.SetupRoutes(router)

	log.Fatal(router.Run(":8080"))
}
