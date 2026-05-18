package router

import (
	"net/http"

	"github.com/apariciocch/psicosstcloud/internal/config"
	"github.com/apariciocch/psicosstcloud/internal/handler"
	"github.com/apariciocch/psicosstcloud/internal/middleware"
	"github.com/apariciocch/psicosstcloud/internal/service/jwt"
	"github.com/go-chi/chi/v5"
)

// SetupRoutes configura todas las rutas de la API
func SetupRoutes(
	router *chi.Mux,
	authHandler *handler.AuthHandler,
	jwtMgr *jwt.Manager,
	rateLimiter *middleware.RateLimiter,
) {
	// Middlewares globales
	router.Use(middleware.SecurityHeaders)
	router.Use(middleware.CORSMiddleware([]string{"*"})) // Configurar origins en config
	router.Use(middleware.RequestLogger)
	router.Use(middleware.RateLimit(rateLimiter))

	// ==================== AUTH ====================
	router.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/login", authHandler.Login)
		r.Post("/refresh", authHandler.RefreshToken)

		// Rutas protegidas
		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTAuth(jwtMgr))
			r.Post("/logout", authHandler.Logout)
			r.Post("/change-password", authHandler.ChangePassword)
		})
	})

	// ==================== HEALTH ====================
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// ==================== API V1 ====================
	router.Route("/api/v1", func(r chi.Router) {
		// Middleware de autenticación en todas las rutas v1
		r.Use(middleware.JWTAuth(jwtMgr))

		// ==================== USERS ====================
		r.Route("/users", func(r chi.Router) {
			r.Get("/", middleware.CheckPermission(config.PermUsersView)(
				http.HandlerFunc(nil), // Será reemplazado por handler real
			))
			r.Post("/", middleware.CheckPermission(config.PermUsersCreate)(
				http.HandlerFunc(nil),
			))
			r.Get("/{id}", middleware.CheckPermission(config.PermUsersView)(
				http.HandlerFunc(nil),
			))
			r.Put("/{id}", middleware.CheckPermission(config.PermUsersEdit)(
				http.HandlerFunc(nil),
			))
			r.Delete("/{id}", middleware.CheckPermission(config.PermUsersDelete)(
				http.HandlerFunc(nil),
			))
		})

		// ==================== EMPLOYEES ====================
		r.Route("/employees", func(r chi.Router) {
			r.Get("/", middleware.CheckPermission(config.PermEmployeesView)(
				http.HandlerFunc(nil),
			))
			r.Post("/", middleware.CheckPermission(config.PermEmployeesCreate)(
				http.HandlerFunc(nil),
			))
			r.Get("/{id}", middleware.CheckPermission(config.PermEmployeesView)(
				http.HandlerFunc(nil),
			))
			r.Put("/{id}", middleware.CheckPermission(config.PermEmployeesEdit)(
				http.HandlerFunc(nil),
			))
		})

		// ==================== OBSERVATIONS ====================
		r.Route("/observations", func(r chi.Router) {
			r.Get("/", middleware.CheckPermission(config.PermObservationsView)(
				http.HandlerFunc(nil),
			))
			r.Post("/", middleware.CheckPermission(config.PermObservationsCreate)(
				http.HandlerFunc(nil),
			))
			r.Get("/{id}", middleware.CheckPermission(config.PermObservationsView)(
				http.HandlerFunc(nil),
			))
			r.Put("/{id}", middleware.CheckPermission(config.PermObservationsEdit)(
				http.HandlerFunc(nil),
			))
			r.Delete("/{id}", middleware.CheckPermission(config.PermObservationsDelete)(
				http.HandlerFunc(nil),
			))
		})

		// ==================== INCIDENTS ====================
		r.Route("/incidents", func(r chi.Router) {
			r.Get("/", middleware.CheckPermission(config.PermIncidentsView)(
				http.HandlerFunc(nil),
			))
			r.Post("/", middleware.CheckPermission(config.PermIncidentsCreate)(
				http.HandlerFunc(nil),
			))
			r.Get("/{id}", middleware.CheckPermission(config.PermIncidentsView)(
				http.HandlerFunc(nil),
			))
			r.Put("/{id}", middleware.CheckPermission(config.PermIncidentsEdit)(
				http.HandlerFunc(nil),
			))
		})

		// ==================== CORRECTIVE ACTIONS ====================
		r.Route("/corrective-actions", func(r chi.Router) {
			r.Get("/", middleware.CheckPermission(config.PermCorrectiveActionsView)(
				http.HandlerFunc(nil),
			))
			r.Post("/", middleware.CheckPermission(config.PermCorrectiveActionsCreate)(
				http.HandlerFunc(nil),
			))
			r.Get("/{id}", middleware.CheckPermission(config.PermCorrectiveActionsView)(
				http.HandlerFunc(nil),
			))
			r.Put("/{id}", middleware.CheckPermission(config.PermCorrectiveActionsEdit)(
				http.HandlerFunc(nil),
			))
		})

		// ==================== ANALYTICS ====================
		r.Route("/analytics", func(r chi.Router) {
			r.Get("/dashboard", middleware.CheckPermission(config.PermAnalyticsView)(
				http.HandlerFunc(nil),
			))
			r.Get("/kpis", middleware.CheckPermission(config.PermAnalyticsView)(
				http.HandlerFunc(nil),
			))
			r.Get("/trends", middleware.CheckPermission(config.PermAnalyticsView)(
				http.HandlerFunc(nil),
			))
		})

		// ==================== REPORTS ====================
		r.Route("/reports", func(r chi.Router) {
			r.Get("/", middleware.CheckPermission(config.PermReportsGenerate)(
				http.HandlerFunc(nil),
			))
			r.Post("/", middleware.CheckPermission(config.PermReportsGenerate)(
				http.HandlerFunc(nil),
			))
			r.Get("/{id}/download", middleware.CheckPermission(config.PermReportsExport)(
				http.HandlerFunc(nil),
			))
		})

		// ==================== AUDIT LOGS ====================
		r.Route("/audit-logs", func(r chi.Router) {
			r.Get("/", middleware.CheckPermission(config.PermAuditLogsView)(
				http.HandlerFunc(nil),
			))
			r.Get("/{id}", middleware.CheckPermission(config.PermAuditLogsView)(
				http.HandlerFunc(nil),
			))
		})

		// ==================== NOTIFICATIONS ====================
		r.Route("/notifications", func(r chi.Router) {
			r.Get("/", http.HandlerFunc(nil))
			r.Get("/unread", http.HandlerFunc(nil))
			r.Put("/{id}/read", http.HandlerFunc(nil))
			r.Put("/mark-all-as-read", http.HandlerFunc(nil))
		})

		// ==================== WORK AREAS ====================
		r.Route("/work-areas", func(r chi.Router) {
			r.Get("/", http.HandlerFunc(nil))
			r.Get("/{id}", http.HandlerFunc(nil))
		})
	})
}
