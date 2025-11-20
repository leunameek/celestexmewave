package router

import (
	"github.com/gin-gonic/gin"
	"github.com/leunameek/celestexmewave/handlers"
	"github.com/leunameek/celestexmewave/internal/middleware"
)

// SetupRouter sets up all routes
func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Apply global middleware
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.ErrorHandlingMiddleware())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API routes
	api := router.Group("/api")
	{
		// Auth routes (no auth required)
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
			auth.POST("/refresh-token", handlers.RefreshToken)
			auth.POST("/logout", handlers.Logout)
			auth.POST("/request-password-reset", handlers.RequestPasswordReset)
			auth.POST("/verify-reset-code", handlers.VerifyResetCode)
		}

		// User routes (auth required)
		users := api.Group("/users")
		users.Use(middleware.AuthMiddleware())
		{
			users.GET("/profile", handlers.GetProfile)
			users.PUT("/profile", handlers.UpdateProfile)
			users.PUT("/change-password", handlers.ChangePassword)
			users.DELETE("/profile", handlers.DeleteProfile)
		}

		// Product routes (no auth required)
		products := api.Group("/products")
		{
			products.GET("", handlers.GetAllProducts)
			products.GET("/:id", handlers.GetProductByID)
			products.GET("/store/:store_id", handlers.GetProductsByStore)
			products.GET("/category/:category", handlers.GetProductsByCategory)
			products.GET("/images/*filename", handlers.ServeImage)
		}

		// Cart routes (optional auth)
		cart := api.Group("/cart")
		cart.Use(middleware.OptionalAuthMiddleware())
		{
			cart.GET("", handlers.GetCart)
			cart.POST("/items", handlers.AddItem)
			cart.PUT("/items/:item_id", handlers.UpdateItem)
			cart.DELETE("/items/:item_id", handlers.RemoveItem)
			cart.DELETE("", handlers.ClearCart)
		}

		// Order routes (optional auth)
		orders := api.Group("/orders")
		orders.Use(middleware.OptionalAuthMiddleware())
		{
			orders.POST("", handlers.CreateOrder)
			orders.GET("/:id", handlers.GetOrder)
			orders.GET("", handlers.GetOrders)
			orders.POST("/:id/payment", handlers.ProcessPayment)
			orders.GET("/:id/confirmation", handlers.GetConfirmation)
		}
	}

	return router
}
