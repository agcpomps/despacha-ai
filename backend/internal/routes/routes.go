package routes

import (
	"github.com/agcpomps/despacha-ai/backend/internal/config"
	"github.com/agcpomps/despacha-ai/backend/internal/handler"
	appmiddleware "github.com/agcpomps/despacha-ai/backend/internal/middleware"

	"github.com/labstack/echo/v5"
)

type RouteHandlers struct {
	AuthHandler     *handler.AuthHandler
	CategoryHandler *handler.CategoryHandler
	ListingHandler  *handler.ListingHandler
}

func RegisterRoutes(e *echo.Echo, cfg *config.Config, h RouteHandlers) {
	api := e.Group("/api/v1")

	registerHealthRoutes(e)
	registerAuthRoutes(api, h.AuthHandler)
	registercategoryRoutes(api, h.CategoryHandler)
	registerListingRountes(api, cfg, h.ListingHandler)
	registerUserRoutes(api, cfg)
	registerAdminRoutes(api, cfg)
	registerModerationRoutes(api, cfg)
}

func registerHealthRoutes(e *echo.Echo) {
	e.GET("/health", func(c *echo.Context) error {
		return c.JSON(200, map[string]string{
			"status": "ok",
			"app":    "despacha-ai",
		})
	})
}

func registerListingRountes(api *echo.Group, cfg *config.Config, listingHandler *handler.ListingHandler) {
	listings := api.Group("/listings")

	// public routes
	listings.GET("", listingHandler.GetListings)
	listings.GET("/:id", listingHandler.GetListingByID)

	// proctedted routes
	listings.POST(
		"",
		listingHandler.CreateListing,
		appmiddleware.AuthMiddleware(cfg.JWTSecret),
	)
}

func registerAuthRoutes(api *echo.Group, authHandler *handler.AuthHandler) {
	auth := api.Group("/auth")

	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
}

func registerUserRoutes(api *echo.Group, cfg *config.Config) {
	private := api.Group("")

	private.Use(appmiddleware.AuthMiddleware(cfg.JWTSecret))

	private.GET("/me", func(c *echo.Context) error {
		userID := c.Get("user_id").(string)
		role := c.Get("user_role").(string)
		phone := c.Get("user_phone").(string)

		return c.JSON(200, map[string]string{
			"user_id": userID,
			"role":    role,
			"phone":   phone,
		})
	})
}

func registerAdminRoutes(api *echo.Group, cfg *config.Config) {
	admin := api.Group("/admin")

	admin.Use(appmiddleware.AuthMiddleware(cfg.JWTSecret))
	admin.Use(appmiddleware.RequireRole("admin"))

	admin.GET("/dashboard", func(c *echo.Context) error {
		return c.JSON(200, map[string]string{
			"message": "admin dashboard",
		})
	})
}

func registerModerationRoutes(api *echo.Group, cfg *config.Config) {
	moderation := api.Group("/moderation")

	moderation.Use(appmiddleware.AuthMiddleware(cfg.JWTSecret))
	moderation.Use(appmiddleware.RequireRole("admin", "moderator"))

	moderation.GET("/listings", func(c *echo.Context) error {
		return c.JSON(200, map[string]string{
			"message": "moderation listings",
		})
	})
}

func registercategoryRoutes(api *echo.Group, categegoryHandler *handler.CategoryHandler) {
	categories := api.Group("/categories")
	categories.GET("", categegoryHandler.GetCategories)
}
