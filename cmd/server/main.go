package main

import (
	"log"
	"mobile-banking/internal/config"
	"mobile-banking/internal/handler"
	"mobile-banking/internal/repository"
	"mobile-banking/internal/service"
	"mobile-banking/pkg/database"

	"github.com/gin-gonic/gin"
)

func main(){
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	db, err := database.NewConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	defer db.Close()

	userRepo := repository.NewUserRepository(db.DB)

	userService := service.NewUserService(userRepo)

	userHandler := handler.NewUserHandler(userService)

	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		v1.POST("/users/register", userHandler.Register)
		v1.POST("/users/login", userHandler.Login)
		
	} 
	
	if err := router.Run(cfg.ServerAddress); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

}

