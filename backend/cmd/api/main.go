package main

import (
	"log"

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

	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	listingRepo := repository.NewListingRepository(db)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	categoryService := service.NewCategoryService(categoryRepo)
	listingService := service.NewListingService(listingRepo)
	authHandler := handler.NewAuthHandler(authService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	listingHandler := handler.NewListingHandler(listingService)

	e := echo.New()

	e.Use(echomiddleware.RequestLogger())
	e.Use(echomiddleware.Recover())
	//e.Use(middleware.CORS())

	// routes
	routes.RegisterRoutes(e, cfg, routes.RouteHandlers{
		AuthHandler:     authHandler,
		CategoryHandler: categoryHandler,
		ListingHandler:  listingHandler,
	})

	log.Println("Despacha Aí Api running on port", cfg.Port)

	if err := e.Start(":" + cfg.Port); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
