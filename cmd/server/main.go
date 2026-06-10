package main

import (
	"log"

	"task-manager-api/internal/auth"
	"task-manager-api/internal/config"
	"task-manager-api/internal/handlers"
	"task-manager-api/internal/models"
	"task-manager-api/internal/repository"
	"task-manager-api/internal/routes"
	"task-manager-api/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	if err := config.ConnectDB(); err != nil {
		log.Fatal(err)
	}

	if err := config.DB.AutoMigrate(&models.User{}, &models.Task{}); err != nil {
		log.Fatal(err)
	}

	jwtManager := auth.NewJWTManager(cfg.JWTSecret, cfg.TokenTTL)

	userRepo := repository.NewUserRepository(config.DB)
	taskRepo := repository.NewTaskRepository(config.DB)

	authService := services.NewAuthService(userRepo, jwtManager)
	taskService := services.NewTaskService(taskRepo)

	authHandler := handlers.NewAuthHandler(authService, jwtManager, cfg.CookieSecure)
	taskHandler := handlers.NewTaskHandler(taskService)

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.FrontendOrigin},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	routes.RegisterRoutes(r, authHandler, taskHandler, jwtManager)

	log.Printf("Server listening on :%s", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
