package main

import (
	"log"

	"github.com/agcpomps/despacha-ai/backend/internal/config"
	"github.com/agcpomps/despacha-ai/backend/internal/database"
	"github.com/agcpomps/despacha-ai/backend/internal/handler"
	"github.com/agcpomps/despacha-ai/backend/internal/repository"
	"github.com/agcpomps/despacha-ai/backend/internal/service"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	cfg := config.Load()

	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}

	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	authHandler := handler.NewAuthHandler(authService)

	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	//e.Use(middleware.CORS())

	// routes
	api := e.Group("/api/v1")

	auth := api.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)

	e.GET("/health", func(c *echo.Context) error {
		return c.JSON(200, map[string]string{
			"status": "ok",
			"app":    "despacha-ai",
		})
	})

	log.Println("Despacha Aí Api running on port", cfg.Port)

	if err := e.Start(":" + cfg.Port); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
