package routes

import (
	"context"
	"time"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"

	"mondash-backend/api"
	"mondash-backend/logger"
	"mondash-backend/repository"
	"mondash-backend/repository/inmemory"
	mongorepo "mondash-backend/repository/mongo"
	"mondash-backend/routes/middlewares"
	"mondash-backend/services"
)

// NewRouter sets up the application routes and returns a chi router.
func NewRouter(db *mongo.Database) *chi.Mux {
	if logger.Log == nil {
		_ = logger.Init()
	}
	router := chi.NewRouter()
	router.Use(middlewares.CORSMiddleware)
	router.Use(middlewares.LoggingMiddleware)

	var (
		nodeRepo   repository.NodeRepository
		appRepo    repository.AppRepository
		alertRepo  repository.AlertRepository
		mapRepo    repository.MapRepository
		deviceRepo repository.DeviceRepository
		authRepo   repository.AuthRepository
		userRepo   repository.UserRepository
	)

	if db == nil {
		logger.Log.Info("Using in-memory repositories")
		nodeRepo = inmemory.NewNodeRepo()
		appRepo = inmemory.NewAppRepo()
		alertRepo = inmemory.NewAlertRepo()
		mapRepo = inmemory.NewMapRepo()
		deviceRepo = inmemory.NewDeviceRepo(nodeRepo.(*inmemory.NodeRepo))
		authRepo = inmemory.NewAuthRepo()
		userRepo = inmemory.NewUserRepo(authRepo.(*inmemory.AuthRepo))
	} else {
		logger.Log.Info("Using MongoDB repositories")
		nodeRepo = mongorepo.NewNodeRepo(db)
		appRepo = mongorepo.NewAppRepo(db)
		alertRepo = mongorepo.NewAlertRepo(db)
		mapRepo = mongorepo.NewMapRepo(db)
		deviceRepo = mongorepo.NewDeviceRepo(db)
		authRepo = mongorepo.NewAuthRepo(db)
		userRepo = mongorepo.NewUserRepo(db)
	}

	nodeService := &services.NodeService{Repo: nodeRepo, DeviceRepo: deviceRepo}
	appService := &services.AppService{Repo: appRepo}
	alertService := &services.AlertService{
		Repo:       alertRepo,
		DeviceRepo: deviceRepo,
	}
	alertService.InitFromEnv()
	mapService := &services.MapService{Repo: mapRepo}
	deviceService := &services.DeviceService{Repo: deviceRepo}
	userService := &services.UserService{Repo: userRepo}
	authService := &services.AuthService{Repo: authRepo}

	_ = alertService.Load()
	alertService.StartMonitoring(context.Background(), time.Second*5)

	router.Get("/healthcheck", api.HealthcheckHandler)

	// API routes used by the frontend
	router.Route("/api", func(r chi.Router) {
		r.Post("/login", api.LoginHandler(authService))
		r.Post("/register", api.RegisterHandler(authService))

		r.Group(func(pr chi.Router) {
			pr.Use(middlewares.CookieAuthMiddleware)
			pr.Get("/apps", api.AppsHandler(appService))
			pr.Get("/apps-timeline", api.AppsTimelineHandler(appService))
			pr.Get("/alerts", api.AlertsHandler(deviceService, alertService))
			pr.Post("/alert", api.RegisterAlertHandler(alertService))
			pr.Get("/active-alerts", api.ActiveAlertsHandler(alertService))
			pr.Get("/nodes", api.NodesHandler(nodeService))
			pr.Get("/map", api.MapHandler(mapService))
			pr.Get("/devices", api.DevicesHandler(deviceService))
			pr.Get("/users", api.UsersHandler(userService))
		})
	})

	router.Group(func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware)
		r.Post("/update-node", api.UpdateNodeHandler(nodeService))
		r.Post("/update-app", api.UpdateAppHandler(appService))
	})

	return router
}
