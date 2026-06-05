// Package main is the entry point for the backend server.
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	// Internal packets import
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/auth"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/config"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/database"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/models"
	notificationHandlers "github.com/isw2-unileon/proyect-scaffolding/backend/internal/notifications"
	handlers "github.com/isw2-unileon/proyect-scaffolding/backend/internal/product"
	producthandlers "github.com/isw2-unileon/proyect-scaffolding/backend/internal/product"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func main() {
	ctx := context.Background()

	// 1. Load configuration
	cfg := config.Load()

	// 2. PostgreSQL connection and table migration
	database.Connect(cfg.DatabaseURL)

	err := database.DB.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Tracking{},
		&models.PriceHistory{},
		&models.Notification{},
	)
	if err != nil {
		logger.Error("Error al migrar la base de datos", "error", err)
		os.Exit(1)
	}
	slog.Info("Tablas sincronizadas correctamente en PostgreSQL")

	// 3. Web service configuration
	gin.SetMode(cfg.GinMode)

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"}, // Añade aquí el puerto exacto de tu React (ej: Vite usa 5173)
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// --- ENRUTAMIENTO GENERAL DE LA API ---
	api := r.Group("/api")
	{
		api.GET("/hello", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Hello from the API"})
		})

		// Rutas Públicas (No requieren login)
		api.POST("/login", auth.LoginHandler)
		api.POST("/logout", auth.LogoutHandler)
		api.POST("/register", auth.RegisterHandler)

		// Rutas Protegidas por JWT
		protected := api.Group("/")
		protected.Use(auth.AuthMiddleware())
		{
			// Endpoints de configuración del perfil de usuario
			protected.GET("/user/profile", auth.GetProfileHandler)
			protected.PUT("/user/profile", auth.UpdateProfileHandler)
			protected.PUT("/user/profile/password", auth.UpdatePasswordHandler)

			// Rutas del tracker
			protected.POST("/track", producthandlers.AddProduct)
			protected.GET("/products/search", producthandlers.SearchProducts)
			protected.POST("/tracking", producthandlers.UpdateTracking)

			// Centro de Notificaciones
			protected.GET("/user/notifications", notificationHandlers.GetUserNotifications)
			protected.PATCH("/user/notifications/:id", notificationHandlers.MarkNotificationAsRead)
			protected.DELETE("/user/notifications/:id", notificationHandlers.DeleteNotification)

			// protected.POST("/track", handlers.AddProduct)
			// protected.GET("/products/search", handlers.SearchProducts)
			protected.GET("/products/:id", handlers.GetProductByID)
			protected.DELETE("/tracking/:id", handlers.UnfollowProduct)
			// protected.POST("/tracking", handlers.UpdateTracking)
			protected.GET("/check-tracking/:id", handlers.CheckTracking)
			protected.GET("/user/saved-products", handlers.GetSavedProducts)
			protected.GET("/products/:id/history", handlers.GetPriceHistory)

		}
	}

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		slog.Info("server listening", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	slog.Info("shutting down server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("shutdown error", "error", err)
	}

	logger.Info("server stopped")
}
