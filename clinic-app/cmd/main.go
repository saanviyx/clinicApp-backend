package main

import (
	"clinic-app/cmd/rest"
	"clinic-app/internal/config"
	"clinic-app/internal/constants"
	"clinic-app/pkg/adapters"
	"clinic-app/pkg/infra"
	adminRepo "clinic-app/pkg/repository/admin"
	appointmentsRepo "clinic-app/pkg/repository/appointments"
	authenticationRepo "clinic-app/pkg/repository/authentication"
	doctorRepo "clinic-app/pkg/repository/doctor"
	"clinic-app/pkg/services"
	adminUsecase "clinic-app/pkg/usecase/admin"
	appointmentsUsecase "clinic-app/pkg/usecase/appointments"
	authenticationUsecase "clinic-app/pkg/usecase/authentication"
	doctorUsecase "clinic-app/pkg/usecase/doctor"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	// ========= Load Environment Variables =========
	cfg := config.LoadConfig()

	// ========= Setup Adapters =========
	adpt, err := adapters.SetupAdapters(&adapters.Options{
		DBConnStr:      cfg.DBConnStr,
		DBMaxIdleConns: cfg.DBMaxIdleConns,
		DBMaxOpenConns: cfg.DBMaxOpenConns,
	})
	if err != nil {
		log.Fatal("Error setting up adapters", zap.Error(err))
	}

	// ========= Setup Infrastructure =========
	infrastructure, err := infra.NewInfrastructure(&infra.Options{
		DB:            adpt.DB,
		Logger:        adpt.Logger,
		MigrationPath: constants.MigrationPath,
	})
	if err != nil {
		log.Fatal("Error initializing infrastructure", zap.Error(err))
	}

	// ========= Setup Repositories =========
	adminRepo := adminRepo.New()
	authRepo := authenticationRepo.New()
	aptmtRepo := appointmentsRepo.New()
	doctorRepo := doctorRepo.New()

	// ========= Setup Services =========
	err = services.SetupService(&services.Options{
		DB:     infrastructure.DB,
		Logger: infrastructure.Logger,
	})
	if err != nil {
		log.Fatal("Error setting up services", zap.Error(err))
	}

	// ========= Setup Usecases =========
	adminUsecase := adminUsecase.New(
		adminRepo,
	)
	authUsecase := authenticationUsecase.New(
		authRepo,
	)
	aptmtsUsecase := appointmentsUsecase.New(
		aptmtRepo,
	)
	doctorUsecase := doctorUsecase.New(
		doctorRepo,
	)

	// ========= Setup Handler =========
	restHandler := rest.NewRestHandler(
		authUsecase, aptmtsUsecase, doctorUsecase, adminUsecase)

	// ========= Setup Router =========
	r := restHandler.SetupRouter(infrastructure.Logger)

	// ========= Start Server =========
	go func() {
		infrastructure.Logger.Info("Starting server on port 8080")
		if err := r.Run(":8080"); err != nil {
			log.Fatal("Error starting server", zap.Error(err))
		}
	}()

	// ========= Handle Shutdown =========
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	infrastructure.Logger.Info("Shutting down server")

}
