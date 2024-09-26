package main

import (
	"github.com/DKhorkov/hmtm-sso/pkg/logging"
	"github.com/DKhorkov/medods/internal/app"
	"github.com/DKhorkov/medods/internal/config"
	httpcontroller "github.com/DKhorkov/medods/internal/controllers/http"
	"github.com/DKhorkov/medods/internal/database"
	"github.com/DKhorkov/medods/internal/repositories"
	"github.com/DKhorkov/medods/internal/services"
	"github.com/DKhorkov/medods/internal/usecases"
)

func main() {
	settings := config.New()

	logger := logging.GetInstance(
		settings.Logging.Level,
		settings.Logging.LogFilePath,
	)

	dbConnector, err := database.New(
		settings.Databases.PostgreSQL,
		logger,
	)

	if err != nil {
		panic(err)
	}

	defer dbConnector.CloseConnection()

	authRepository := &repositories.CommonAuthRepository{DBConnector: dbConnector}
	authService := &services.CommonAuthService{AuthRepository: authRepository}
	useCases := &usecases.CommonUseCases{
		AuthService: authService,
		HashCost:    settings.Security.HashCost,
		JWTConfig:   settings.Security.JWT,
	}

	controller := httpcontroller.New(
		settings.HTTP.Host,
		settings.HTTP.Port,
		useCases,
		logger,
	)

	application := app.New(controller)
	application.Run()
}
