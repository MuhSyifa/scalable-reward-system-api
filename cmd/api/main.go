package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"scalable-reward-system/internal/cache"
	"scalable-reward-system/internal/config"
	"scalable-reward-system/internal/database"
	"scalable-reward-system/internal/delivery/http/handler"
	"scalable-reward-system/internal/delivery/http/middleware"
	"scalable-reward-system/internal/logger"
	"scalable-reward-system/internal/repository"
	"scalable-reward-system/internal/service"
	"scalable-reward-system/internal/worker"
)

func main() {
	// Load config
	cfg := config.LoadConfig()

	// Init logger
	logger.InitLogger(cfg.AppEnv)

	// Init dependencies
	db := database.NewPostgresDB(cfg)
	defer db.Close()

	redisClient := cache.NewRedisClient(cfg)
	defer redisClient.Close()

	// Start Background Workers
	workerCtx, workerCancel := context.WithCancel(context.Background())
	defer workerCancel()

	pointWorker := worker.NewPointExpirationWorker(db)
	pointWorker.Start(workerCtx, 24*time.Hour) // Run daily

	// Init Repositories
	userRepo := repository.NewUserRepository(db)
	campaignRepo := repository.NewCampaignRepository(db)
	voucherRepo := repository.NewVoucherRepository(db)
	txRepo := repository.NewTransactionRepository(db)

	// Init Services
	userService := service.NewUserService(userRepo, cfg)
	redemptionService := service.NewRedemptionService(db, redisClient, userRepo, voucherRepo, txRepo)

	// Init Handlers
	authHandler := handler.NewAuthHandler(userService)
	campaignHandler := handler.NewCampaignHandler(campaignRepo, voucherRepo)
	redemptionHandler := handler.NewRedemptionHandler(redemptionService)

	// Setup Gin router
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	// Apply Rate Limiter globally for public routes (e.g. 100 req per minute)
	router.Use(middleware.RateLimiter(redisClient, 100, time.Minute))

	// Routes
	v1 := router.Group("/api/v1")
	{
		// Public routes
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "up"})
		})
		v1.POST("/register", authHandler.Register)
		v1.POST("/login", authHandler.Login)
		v1.GET("/campaigns", campaignHandler.ListActiveCampaigns)
		v1.GET("/campaigns/:id/vouchers", campaignHandler.ListVouchers)

		// Protected routes
		protected := v1.Group("/")
		protected.Use(middleware.RequireAuth(cfg.JWTSecret))
		{
			// User endpoints
			protected.POST("/redeem", middleware.Idempotency(redisClient), redemptionHandler.Redeem)
			
			// Admin endpoints
			admin := protected.Group("/admin")
			admin.Use(middleware.RequireRole("admin"))
			{
				// admin routes would go here
			}
		}
	}

	// Server config
	srv := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		log.Info().Msgf("Starting server on port %s", cfg.AppPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("Failed to listen and serve")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exiting")
}
