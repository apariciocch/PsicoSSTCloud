package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/apariciocch/psicosstcloud/internal/config"
	"github.com/apariciocch/psicosstcloud/internal/domain/entities"
	"github.com/apariciocch/psicosstcloud/internal/handler"
	"github.com/apariciocch/psicosstcloud/internal/middleware"
	"github.com/apariciocch/psicosstcloud/internal/pkg/logger"
	"github.com/apariciocch/psicosstcloud/internal/repository"
	"github.com/apariciocch/psicosstcloud/internal/router"
	authService "github.com/apariciocch/psicosstcloud/internal/service/auth"
	"github.com/apariciocch/psicosstcloud/internal/service/jwt"
	"github.com/apariciocch/psicosstcloud/internal/service/password"
	"github.com/go-chi/chi/v5"
)

func main() {
	// Cargar configuración
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	// Inicializar logger
	if err := logger.Init(cfg.Logs.Level, cfg.Logs.Format); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing logger: %v\n", err)
		os.Exit(1)
	}

	logger.Info("Starting PsicoSST Cloud Backend",
		logger.WithRequest("", "", 0, 0)...,
	)

	// Conectar a PocketBase
	pbClient, err := repository.NewPocketBaseClient(
		cfg.PocketBase.URL,
		cfg.PocketBase.AdminEmail,
		cfg.PocketBase.AdminPassword,
	)
	if err != nil {
		logger.Fatal("Error connecting to PocketBase", logger.WithError(err))
	}
	defer pbClient.Close()

	logger.Info("Connected to PocketBase successfully")

	// Inicializar servicios
	jwtManager := jwt.NewManager(
		cfg.JWT.Secret,
		cfg.JWT.Expiration,
		cfg.JWT.RefreshExpiration,
	)

	passwordManager := password.NewManager()

	// Inicializar repositorios en memoria (TODO: conectar a PocketBase cuando esté listo)
	userRepo := repository.NewInMemoryUserRepository()
	roleRepo := repository.NewInMemoryRoleRepository()

	// Crear un audit service simple
	auditService := &SimpleAuditService{}

	// Crear servicio de autenticación
	authServiceImpl := authService.NewAuthService(
		userRepo,
		roleRepo,
		jwtManager,
		passwordManager,
		auditService,
	)

	// Crear handlers
	authHandler := handler.NewAuthHandler(authServiceImpl)

	// Crear router
	r := chi.NewRouter()

	// Rate limiter
	rateLimiter := middleware.NewRateLimiter(
		float64(cfg.Security.RateLimitRequests)/float64(cfg.Security.RateLimitWindow),
		1000, // max visitors
	)

	// Configurar rutas
	router.SetupRoutes(r, authHandler, jwtManager, rateLimiter)

	// Crear servidor HTTP
	srv := &http.Server{
		Addr:           fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:        r,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	// Iniciar servidor en goroutine
	go func() {
		logger.Info(fmt.Sprintf("Server starting on %s", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server error", logger.WithError(err))
		}
	}()

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	logger.Info("Shutdown signal received, gracefully shutting down...")

	// Context con timeout para shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown error", logger.WithError(err))
	}

	logger.Info("Server stopped")
}

// SimpleAuditService implementación simple de auditoría
type SimpleAuditService struct{}

func (s *SimpleAuditService) Log(ctx context.Context, audit *entities.AuditLog) error {
	// Por ahora solo log, implementación real cuando se integre con base de datos
	return nil
}
