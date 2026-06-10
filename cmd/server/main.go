package main

import (
	"log"

	"task-manager-api/internal/auth"
	"task-manager-api/internal/config"
	"task-manager-api/internal/handlers"
	"task-manager-api/internal/models"
	"task-manager-api/internal/realtime"
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

	if err := config.DB.AutoMigrate(
		&models.User{},
		&models.Task{},
		&models.TaskActivity{},
		&models.Attachment{},
	); err != nil {
		log.Fatal(err)
	}

	jwtManager := auth.NewJWTManager(cfg.JWTSecret, cfg.TokenTTL)
	hub := realtime.NewHub()

	userRepo := repository.NewUserRepository(config.DB)
	taskRepo := repository.NewTaskRepository(config.DB)
	activityRepo := repository.NewActivityRepository(config.DB)
	attachmentRepo := repository.NewAttachmentRepository(config.DB)

	authService := services.NewAuthService(userRepo, jwtManager, cfg.AdminEmails)
	taskService := services.NewTaskService(taskRepo, activityRepo, hub)
	attachmentService := services.NewAttachmentService(
		attachmentRepo, taskRepo, taskService, cfg.UploadDir, cfg.MaxUploadMB<<20,
	)

	r := gin.Default()
	r.MaxMultipartMemory = 8 << 20
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.FrontendOrigin},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	routes.RegisterRoutes(r, routes.Handlers{
		Auth:        handlers.NewAuthHandler(authService, jwtManager, cfg.CookieSecure),
		Tasks:       handlers.NewTaskHandler(taskService),
		Attachments: handlers.NewAttachmentHandler(attachmentService),
		Events:      handlers.NewEventsHandler(hub),
	}, jwtManager)

	log.Printf("Server listening on :%s", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
