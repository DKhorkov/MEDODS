package httpcontroller

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/DKhorkov/hmtm-sso/pkg/logging"
	"github.com/DKhorkov/medods/internal/interfaces"
)

type Controller struct {
	httpServer *http.ServeMux
	host       string
	port       int
	logger     *slog.Logger
}

// Run HTTP server.
func (controller *Controller) Run() {
	controller.logger.Info(
		fmt.Sprintf("Starting HTTP Server at http://%s:%d", controller.host, controller.port),
		"Traceback",
		logging.GetLogTraceback(),
	)

	if err := http.ListenAndServe(
		fmt.Sprintf("%s:%d", controller.host, controller.port),
		controller.httpServer,
	); err != nil {
		controller.logger.Error(
			"Error occurred while listening to HTTP Server",
			"Traceback",
			logging.GetLogTraceback(),
			"Error",
			err,
		)
		panic(err)
	}

	controller.logger.Info("Stopped serving new connections.")
}

// Stop HTTP server gracefully (graceful shutdown).
func (controller *Controller) Stop() {
	controller.logger.Info("Graceful shutdown completed.")
}

// New creates an instance of HTTP Controller.
func New(host string, port int, useCases interfaces.UseCases, logger *slog.Logger) *Controller {
	server := http.NewServeMux()
	server.HandleFunc("/tokens", TokensHandler{UseCases: useCases, Logger: logger}.GetHandleFunc())

	return &Controller{
		httpServer: server,
		port:       port,
		host:       host,
		logger:     logger,
	}
}
