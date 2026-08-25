package delivery

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/irsyadjpp/tracking-diet-backend/internal/auth"
	"github.com/irsyadjpp/tracking-diet-backend/internal/config"
	"github.com/irsyadjpp/tracking-diet-backend/internal/delivery/handlers"
	"github.com/irsyadjpp/tracking-diet-backend/internal/delivery/middleware"
	"github.com/irsyadjpp/tracking-diet-backend/internal/domain"
	"github.com/irsyadjpp/tracking-diet-backend/internal/usecase"
	"github.com/irsyadjpp/tracking-diet-backend/pkg/logger"
)

type HTTPServer struct {
	cfg                 *config.Config
	log                 *logger.Logger
	router              *chi.Mux
	server              *http.Server
	userRepo            domain.UserRepository
	jwtManager          *auth.JWTManager
	userUC              *usecase.UserUsecase
	bodyUC              *usecase.BodyMeasurementUsecase
	nutritionUC         *usecase.NutritionUsecase
	metabolicUC         *usecase.MetabolicUsecase
	fitnessUC           *usecase.FitnessUsecase
	wellbeingUC         *usecase.WellbeingUsecase
	labTestUC           *usecase.LabTestUsecase
	aiRecommendationUC  *usecase.AIRecommendationUsecase
	dashboardUC         *usecase.DashboardUsecase
}

func NewHTTPServer(cfg *config.Config, log *logger.Logger, userRepo domain.UserRepository, jwtManager *auth.JWTManager, userUC *usecase.UserUsecase, bodyUC *usecase.BodyMeasurementUsecase, nutritionUC *usecase.NutritionUsecase, metabolicUC *usecase.MetabolicUsecase, fitnessUC *usecase.FitnessUsecase, wellbeingUC *usecase.WellbeingUsecase, labTestUC *usecase.LabTestUsecase, aiRecommendationUC *usecase.AIRecommendationUsecase, dashboardUC *usecase.DashboardUsecase) *HTTPServer {
	return &HTTPServer{
		cfg:                cfg,
		log:                log,
		userRepo:           userRepo,
		jwtManager:         jwtManager,
		userUC:             userUC,
		bodyUC:             bodyUC,
		nutritionUC:        nutritionUC,
		metabolicUC:        metabolicUC,
		fitnessUC:          fitnessUC,
		wellbeingUC:        wellbeingUC,
		labTestUC:          labTestUC,
		aiRecommendationUC: aiRecommendationUC,
		dashboardUC:        dashboardUC,
	}
}

func (s *HTTPServer) SetupRoutes() {
	s.router = chi.NewRouter()

	// Middleware
	authMiddleware := middleware.NewAuthMiddleware(s.jwtManager)
	loggingMiddleware := middleware.NewLoggingMiddleware(s.log)
	errorHandler := middleware.NewErrorHandler()
	corsConfig := middleware.DefaultCORSConfig()

	s.router.Use(middleware.CORS(corsConfig))
	s.router.Use(loggingMiddleware.RequestLogger)
	s.router.Use(errorHandler.Recover)

	// Health check
	s.router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// API routes
	s.router.Route("/api/v1", func(r chi.Router) {
		// Auth routes (public)
		authHandler := handlers.NewAuthHandler(s.userRepo, s.jwtManager, s.log)
		r.Post("/auth/register", authHandler.Register)
		r.Post("/auth/login", authHandler.Login)
		r.Post("/auth/refresh", authHandler.Refresh)
		
		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.Authenticate)
			r.Get("/auth/me", authHandler.Me)
			
			// User routes
			userHandler := handlers.NewUserHandler(s.userUC, s.log)
			r.Get("/users/me", userHandler.GetProfile)
			r.Put("/users/me", userHandler.UpdateProfile)
			
			// Body measurement routes
			bodyHandler := handlers.NewBodyMeasurementHandler(s.bodyUC, s.log)
			r.Post("/measurements/body", bodyHandler.Create)
			r.Get("/measurements/body", bodyHandler.List)
			r.Get("/measurements/body/{id}", bodyHandler.GetByID)
			r.Put("/measurements/body/{id}", bodyHandler.Update)
			r.Delete("/measurements/body/{id}", bodyHandler.Delete)
			
			// Nutrition routes
			nutritionHandler := handlers.NewNutritionHandler(s.nutritionUC, s.log)
			r.Post("/measurements/nutrition", nutritionHandler.Create)
			r.Get("/measurements/nutrition", nutritionHandler.List)
			r.Get("/measurements/nutrition/daily", nutritionHandler.GetDailySummary)
			r.Get("/measurements/nutrition/{id}", nutritionHandler.GetByID)
			r.Put("/measurements/nutrition/{id}", nutritionHandler.Update)
			r.Delete("/measurements/nutrition/{id}", nutritionHandler.Delete)
			
			// Metabolic measurement routes
			metabolicHandler := handlers.NewMetabolicHandler(s.metabolicUC, s.log)
			r.Post("/measurements/metabolic", metabolicHandler.Create)
			r.Get("/measurements/metabolic", metabolicHandler.List)
			r.Get("/measurements/metabolic/{id}", metabolicHandler.GetByID)
			r.Put("/measurements/metabolic/{id}", metabolicHandler.Update)
			r.Delete("/measurements/metabolic/{id}", metabolicHandler.Delete)
			
			// Fitness measurement routes
			fitnessHandler := handlers.NewFitnessHandler(s.fitnessUC, s.log)
			r.Post("/measurements/fitness", fitnessHandler.Create)
			r.Get("/measurements/fitness", fitnessHandler.List)
			r.Get("/measurements/fitness/{id}", fitnessHandler.GetByID)
			r.Put("/measurements/fitness/{id}", fitnessHandler.Update)
			r.Delete("/measurements/fitness/{id}", fitnessHandler.Delete)
			
			// Wellbeing measurement routes
			wellbeingHandler := handlers.NewWellbeingHandler(s.wellbeingUC, s.log)
			r.Post("/measurements/wellbeing", wellbeingHandler.Create)
			r.Get("/measurements/wellbeing", wellbeingHandler.List)
			r.Get("/measurements/wellbeing/{id}", wellbeingHandler.GetByID)
			r.Put("/measurements/wellbeing/{id}", wellbeingHandler.Update)
			r.Delete("/measurements/wellbeing/{id}", wellbeingHandler.Delete)
			
			// Lab test routes
			labTestHandler := handlers.NewLabTestHandler(s.labTestUC, s.log)
			r.Post("/lab-tests", labTestHandler.Create)
			r.Get("/lab-tests", labTestHandler.List)
			r.Get("/lab-tests/{id}", labTestHandler.GetByID)
			r.Put("/lab-tests/{id}", labTestHandler.Update)
			r.Delete("/lab-tests/{id}", labTestHandler.Delete)
			
			// AI recommendation routes
			aiHandler := handlers.NewAIRecommendationHandler(s.aiRecommendationUC, s.log)
			r.Post("/recommendations/generate", aiHandler.Generate)
			r.Get("/recommendations", aiHandler.List)
			r.Get("/recommendations/{id}", aiHandler.GetByID)
			r.Delete("/recommendations/{id}", aiHandler.Delete)
			
			// Dashboard routes
			dashboardHandler := handlers.NewDashboardHandler(s.dashboardUC, s.log)
			r.Get("/dashboard/stats", dashboardHandler.GetStats)
			r.Get("/dashboard/trend", dashboardHandler.GetProgressTrend)
			r.Get("/dashboard/weekly", dashboardHandler.GetWeeklySummary)
			r.Get("/dashboard/goals", dashboardHandler.GetGoalProgress)
			r.Get("/dashboard/comparison", dashboardHandler.GetComparisonReport)
		})
	})
}

func (s *HTTPServer) Run() error {
	s.SetupRoutes()

	s.server = &http.Server{
		Addr:         ":" + s.cfg.Port,
		Handler:      s.router,
		ReadTimeout:  s.cfg.ReadTimeout,
		WriteTimeout: s.cfg.WriteTimeout,
		IdleTimeout:  s.cfg.IdleTimeout,
	}

	s.log.Info("Starting HTTP server", map[string]interface{}{
		"port": s.cfg.Port,
		"env":  s.cfg.Env,
	})

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	s.log.Info("Shutting down HTTP server", nil)
	return s.server.Shutdown(ctx)
}