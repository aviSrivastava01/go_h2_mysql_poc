package main

import (
	"fmt"
	"log"
	"os"

	"go-redis-mysql-poc/database"
	"go-redis-mysql-poc/handlers"
	"go-redis-mysql-poc/redis"

	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	database.ConnectDB()
	redis.ConnectRedis()

	//router := mux.NewRouter()
	router := gin.Default()

	// Product routes
	router.POST("/products", handlers.CreateProduct)
	router.GET("/products", handlers.ListProducts)
	router.GET("/products/:id", handlers.GetProduct)
	router.PUT("/products/:id", handlers.UpdateProduct)
	router.DELETE("/products/:id", handlers.DeleteProduct)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" // Default port if not specified
	}

	fmt.Println("Server listening on port " + port)
	log.Fatal(router.Run(":" + port)) // Start the server

	//fmt.Println("Server listening on port 8000")
	//log.Fatal(http.ListenAndServe(":8000", router))
}
