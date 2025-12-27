package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"motico-api/config"
	authdomain "motico-api/internal/domain/auth"
	categorydomain "motico-api/internal/domain/category"
	productdomain "motico-api/internal/domain/product"
	stockdomain "motico-api/internal/domain/stock"
	storedomain "motico-api/internal/domain/store"
	transferdomain "motico-api/internal/domain/transfer"
	"motico-api/internal/repository"
	"motico-api/internal/rest"
	categoryhandler "motico-api/internal/rest/category"
	producthandler "motico-api/internal/rest/product"
	stockhandler "motico-api/internal/rest/stock"
	storehandler "motico-api/internal/rest/store"
	transferhandler "motico-api/internal/rest/transfer"
	"motico-api/pkg/logger"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	cfg, err := config.Load("config/config.json")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	appLogger, err := logger.New(cfg.Logging.Level, cfg.Logging.Format)
	if err != nil {
		log.Fatalf("Error creating logger: %v", err)
	}
	defer func() {
		if err := appLogger.Sync(); err != nil {
			log.Printf("Error syncing logger: %v", err)
		}
	}()

	ctx := context.Background()

	pool, err := repository.NewConnectionPool(ctx, cfg)
	if err != nil {
		appLogger.Fatal("Error creating database connection pool", logger.Error(err))
	}
	defer pool.Close()

	appLogger.Info("Database connection established")

	authService := authdomain.NewService(cfg)

	storeRepo := repository.NewStoreRepository(pool)
	categoryRepo := repository.NewCategoryRepository(pool)
	productRepo := repository.NewProductRepository(pool)
	stockRepo := repository.NewStockRepository(pool)
	transferRepo := repository.NewTransferRepository(pool)

	storeService := storedomain.NewService(storeRepo, cfg, appLogger)
	categoryService := categorydomain.NewService(categoryRepo, cfg, appLogger)
	productService := productdomain.NewService(productRepo, cfg, appLogger)
	stockService := stockdomain.NewService(stockRepo, cfg, appLogger)
	transferService := transferdomain.NewService(transferRepo, stockService, storeRepo, cfg, appLogger)

	categoryHandler := categoryhandler.NewHandler(categoryService, cfg)
	storeHandler := storehandler.NewHandler(storeService, cfg)
	productHandler := producthandler.NewHandler(productService, stockService, cfg)
	stockHandler := stockhandler.NewHandler(stockService, cfg)
	transferHandler := transferhandler.NewHandler(transferService, cfg)

	router := rest.NewRouter(rest.RouterDependencies{
		AuthService:     authService,
		CategoryHandler: categoryHandler,
		StoreHandler:    storeHandler,
		ProductHandler:  productHandler,
		StockHandler:    stockHandler,
		TransferHandler: transferHandler,
	})

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  parseDuration(cfg.Server.ReadTimeout),
		WriteTimeout: parseDuration(cfg.Server.WriteTimeout),
	}

	go func() {
		appLogger.Info("Server starting", logger.String("address", server.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatal("Server failed to start", logger.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Server shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		appLogger.Error("Error during server shutdown", logger.Error(err))
	}

	appLogger.Info("Server stopped")
}

func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 30 * time.Second
	}
	return d
}
