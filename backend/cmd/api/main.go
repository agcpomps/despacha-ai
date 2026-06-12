package main

import (
	"context"
	"log"
	"net/http"

	"github.com/agcpomps/despacha-ai/backend/internal/config"
	"github.com/agcpomps/despacha-ai/backend/internal/database"
	"github.com/agcpomps/despacha-ai/backend/internal/handler"

	"github.com/agcpomps/despacha-ai/backend/internal/repository"
	"github.com/agcpomps/despacha-ai/backend/internal/routes"
	"github.com/agcpomps/despacha-ai/backend/internal/service"
	"github.com/labstack/echo/v5"
	echomiddleware "github.com/labstack/echo/v5/middleware"
)

func main() {
	cfg := config.Load()

	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}

	defer db.Close()

	if err := database.RunMigrations(context.Background(), db, cfg.MigrationsPath); err != nil {
		log.Fatal("failed to run migrations: ", err)
	}

	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	listingRepo := repository.NewListingRepository(db)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	categoryService := service.NewCategoryService(categoryRepo)
	listingService := service.NewListingService(listingRepo, categoryRepo, userRepo)
	userService := service.NewUserService(userRepo)
	authHandler := handler.NewAuthHandler(authService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	listingHandler := handler.NewListingHandler(listingService)
	userHandler := handler.NewUserHandler(userService)

	e := echo.New()

	e.Use(echomiddleware.RequestLogger())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins:     cfg.CORSAllowedOrigins,
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))
	e.Static("/uploads", "uploads")

	// routes
	routes.RegisterRoutes(e, cfg, routes.RouteHandlers{
		AuthHandler:     authHandler,
		CategoryHandler: categoryHandler,
		ListingHandler:  listingHandler,
		UserHandler:     userHandler,
	})

	log.Println("Despacha Aí Api running on port", cfg.Port)

	if err := e.Start(":" + cfg.Port); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
