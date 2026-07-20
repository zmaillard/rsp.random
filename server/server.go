package server

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/dgraph-io/badger/v4"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"rsp.random/config"
	"rsp.random/metrics"
	"rsp.random/services"

	echoprometheus "github.com/labstack/echo-prometheus"
	"github.com/labstack/echo/v5/middleware"
)

func NewEchoServer(config *config.Config, httpClient *http.Client, badgerDb *badger.DB, backgroundChan chan services.UpdateCounterProcess) *echo.Echo {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	httpServer := echo.New()

	// Initialize custom Prometheus metrics
	metrics.Init()

	httpServer.Use(echoprometheus.NewMiddleware("rsprandom"))
	go func() {
		metricEcho := echo.New()                                // this Echo will run on separate port 1334
		metricEcho.GET("/metrics", echoprometheus.NewHandler()) // adds route to serve gathered metrics
		if err := metricEcho.Start(":1334"); err != nil {
			slog.Error("failed to start metrics server", "error", err)
		}
	}()

	httpServer.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		HandleError: true, // forwards the error to the global error handler so it can pick the status code
		LogValuesFunc: func(c *echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))
	httpServer.Validator = &CustomValidator{Validator: validator.New()}
	httpServer.Use(middleware.Recover())

	searchService := services.NewSearchService(httpClient, badgerDb, config)
	updateStoreService := services.NewCounterService(badgerDb, config)

	if config.LoadDataAtStartup {
		backgroundChan <- updateStoreService.UpdateData
	}
	addRoutes(httpServer, config, updateStoreService, searchService, backgroundChan)

	return httpServer
}
