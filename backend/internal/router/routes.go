package router

import (
	"net/http"

	"github.com/apariciocch/psicosstcloud/internal/config"
	"github.com/apariciocch/psicosstcloud/internal/handler"
	"github.com/apariciocch/psicosstcloud/internal/middleware"
	"github.com/apariciocch/psicosstcloud/internal/service/jwt"
	"github.com/go-chi/chi/v5"
)

// handlerWrapper convierte un http.Handler a http.HandlerFunc
func handlerWrapper(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	}
}

// notImplementedHandler handler que retorna 501 Not Implemented
func notImplementedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte(`{"error":"Not implemented"}`))
}

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

	// ==================== DASHBOARD ====================
	router.Get("/", handler.DashboardHandler)
	router.Get("/dashboard", handler.DashboardHandler)

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
			r.Get("/", handlerWrapper(middleware.CheckPermission(config.PermUsersView)(
				http.HandlerFunc(notImplementedHandler),
			)))
			r.Post("/", handlerWrapper(middleware.CheckPermission(config.PermUsersCreate)(
				http.HandlerFunc(notImplementedHandler),
			)))
			r.Get("/{id}", handlerWrapper(middleware.CheckPermission(config.PermUsersView)(
				http.HandlerFunc(notImplementedHandler),
			)))
			r.Put("/{id}", handlerWrapper(middleware.CheckPermission(config.PermUsersEdit)(
				http.HandlerFunc(notImplementedHandler),
			)))
			r.Delete("/{id}", handlerWrapper(middleware.CheckPermission(config.PermUsersDelete)(
				http.HandlerFunc(notImplementedHandler),
			)))
		})

		// ==================== EMPLOYEES ====================
		r.Route("/employees", func(r chi.Router) {
			r.Get("/", handlerWrapper(middleware.CheckPermission(config.PermEmployeesView)(
				http.HandlerFunc(notImplementedHandler),
			)))
			r.Post("/", handlerWrapper(middleware.CheckPermission(config.PermEmployeesCreate)(
				http.HandlerFunc(notImplementedHandler),
			)))
			r.Get("/{id}", handlerWrapper(middleware.CheckPermission(config.PermEmployeesView)(
				http.HandlerFunc(notImplementedHandler),
			)))
			r.Put("/{id}", handlerWrapper(middleware.CheckPermission(config.PermEmployeesEdit)(
				http.HandlerFunc(notImplementedHandler),
			)))
		})

		// ==================== OBSERVATIONS ====================
		r.Route("/observations", func(r chi.Router) {
			r.Get("/", handlerWrapper(middleware.CheckPermission(config.PermObservationsView)(
				http.HandlerFunc(notImplementedHandler),
			)))
			r.Post("/", handlerWrapper(middleware.CheckPermission(config.PermObservationsCreate)(
				http.HandlerFunc(notImplementedHandler),
			)))
			r.Get("/{id}", handlerWrapper(middleware.CheckPermission(config.PermObservationsView)(
				http.HandlerFunc(notImplementedHandler),
			)))
			r.Put("/{id}", handlerWrapper(middleware.CheckPermission(config.PermObservationsEdit)(
				http.HandlerFunc(notImplementedHandler),
			)))
			r.Delete("/{id}", handlerWrapper(middleware.CheckPermission(config.PermObservationsDelete)(
				http.HandlerFunc(notImplementedHandler),
			)))
		})

		// ==================== INCIDENTS ====================
		r.Route("/incidents", func(r chi.Router) {
			r.Get("/", handlerWrapper(middleware.CheckPermission(config.PermIncidentsView)(
				http.HandlerFunc(notImplementedHandler),
			)))
			r.Post("/", handlerWrapper(middleware.CheckPermission(config.PermIncidentsCreate)(
				http.HandlerFunc(notImplementedHandler),
			)))
			r.Get("/{id}", handlerWrapper(middleware.CheckPermission(config.PermIncidentsView)(
				http.HandlerFunc(notImplementedHandler),
			)))
			r.Put("/{id}", handlerWrapper(middleware.CheckPermission(config.PermIncidentsEdit)(
				http.HandlerFunc(notImplementedHandler),
			)))
		})

		// ==================== CORRECTIVE ACTIONS ====================
		r.Route("/corrective-actions", func(r chi.Router) {
			r.Get("/", handlerWrapper(middleware.CheckPermission(config.PermCorrectiveActionsView)(
				http.HandlerFunc(notImplementedHandler),
			)))
			r.Post("/", handlerWrapper(middleware.CheckPermission(config.PermCorrectiveActionsCreate)(
				http.HandlerFunc(notImplementedHandler),
			)))
			r.Get("/{id}", handlerWrapper(middleware.CheckPermission(config.PermCorrectiveActionsView)(
				http.HandlerFunc(notImplementedHandler),
			)))
			r.Put("/{id}", handlerWrapper(middleware.CheckPermission(config.PermCorrectiveActionsEdit)(
				http.HandlerFunc(notImplementedHandler),
			)))
		})

		// ==================== ANALYTICS ====================
		r.Route("/analytics", func(r chi.Router) {
			r.Get("/dashboard", handlerWrapper(middleware.CheckPermission(config.PermAnalyticsView)(
				http.HandlerFunc(notImplementedHandler),
			)))
			r.Get("/kpis", handlerWrapper(middleware.CheckPermission(config.PermAnalyticsView)(
				http.HandlerFunc(notImplementedHandler),
			)))
			r.Get("/trends", handlerWrapper(middleware.CheckPermission(config.PermAnalyticsView)(
				http.HandlerFunc(notImplementedHandler),
			)))
		})

		// ==================== REPORTS ====================
		r.Route("/reports", func(r chi.Router) {
			r.Get("/", handlerWrapper(middleware.CheckPermission(config.PermReportsGenerate)(
				http.HandlerFunc(notImplementedHandler),
			)))
			r.Post("/", handlerWrapper(middleware.CheckPermission(config.PermReportsGenerate)(
				http.HandlerFunc(notImplementedHandler),
			)))
			r.Get("/{id}/download", handlerWrapper(middleware.CheckPermission(config.PermReportsExport)(
				http.HandlerFunc(notImplementedHandler),
			)))
		})

		// ==================== AUDIT LOGS ====================
		r.Route("/audit-logs", func(r chi.Router) {
			r.Get("/", handlerWrapper(middleware.CheckPermission(config.PermAuditLogsView)(
				http.HandlerFunc(notImplementedHandler),
			)))
			r.Get("/{id}", handlerWrapper(middleware.CheckPermission(config.PermAuditLogsView)(
				http.HandlerFunc(notImplementedHandler),
			)))
		})

		// ==================== NOTIFICATIONS ====================
		r.Route("/notifications", func(r chi.Router) {
			r.Get("/", http.HandlerFunc(notImplementedHandler))
			r.Get("/unread", http.HandlerFunc(notImplementedHandler))
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
