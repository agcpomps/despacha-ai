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
	UserHandler     *handler.UserHandler
}

func RegisterRoutes(e *echo.Echo, cfg *config.Config, h RouteHandlers) {
	api := e.Group("/api/v1")

	registerHealthRoutes(e)
	registerAuthRoutes(api, h.AuthHandler)
	registercategoryRoutes(api, h.CategoryHandler)
	registerListingRountes(api, cfg, h.ListingHandler)
	registerUserRoutes(api, cfg, h.UserHandler, h.ListingHandler)
	registerAdminRoutes(api, cfg, h.ListingHandler, h.UserHandler)
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
	listings.PUT("/:id", listingHandler.UpdateListing, appmiddleware.AuthMiddleware(cfg.JWTSecret))
	listings.DELETE("/:id", listingHandler.DeleteListing, appmiddleware.AuthMiddleware(cfg.JWTSecret))
}

func registerAuthRoutes(api *echo.Group, authHandler *handler.AuthHandler) {
	auth := api.Group("/auth")

	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
}

func registerUserRoutes(api *echo.Group, cfg *config.Config, userHandler *handler.UserHandler, listingHandler *handler.ListingHandler) {
	private := api.Group("")

	private.Use(appmiddleware.AuthMiddleware(cfg.JWTSecret))

	private.GET("/me", userHandler.GetMe)
	private.GET("/me/listings", listingHandler.GetMyListings)
	private.PUT("/me/password", userHandler.ChangePassword)
}

func registerAdminRoutes(api *echo.Group, cfg *config.Config, listingHandler *handler.ListingHandler, userHandler *handler.UserHandler) {
	admin := api.Group("/admin")

	admin.Use(appmiddleware.AuthMiddleware(cfg.JWTSecret))
	admin.Use(appmiddleware.RequireRole("admin"))

	// user management
	admin.GET("/users", userHandler.GetUsers)
	admin.PUT("/users/:id/role", userHandler.UpdateUserRole)
	admin.PUT("/users/:id/password", userHandler.ResetUserPassword)

	// listing promotion (manual monetization: featured + bump)
	admin.PUT("/listings/:id/feature", listingHandler.FeatureListing)
	admin.DELETE("/listings/:id/feature", listingHandler.UnfeatureListing)
	admin.PUT("/listings/:id/bump", listingHandler.BumpListing)
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
