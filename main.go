package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/prabhat-xs/library-management-backend/config"
	"github.com/prabhat-xs/library-management-backend/routes"
	cors "github.com/rs/cors/wrapper/gin"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	config.ConnectDatabase()

	r := gin.Default()
	
	routes.SetupRoutes(r)
	r.Use(cors.Default())
	
	r.Run()
}
