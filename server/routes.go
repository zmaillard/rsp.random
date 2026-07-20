package server

import (
	"encoding/gob"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"rsp.random/config"
	"rsp.random/server/handlers"
	"rsp.random/services"
)

func addRoutes(server *echo.Echo, c *config.Config, counterService services.CounterService, searchService services.SearchService, backgroundChan chan services.UpdateCounterProcess) {

	gob.Register(map[string]interface{}{})

	// Health Check
	server.GET("/healthz", func(c *echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	server.GET("/health", func(ctx *echo.Context) error {
		return ctx.JSON(http.StatusOK, map[string]interface{}{
			"version": c.VersionNumber,
		})
	})

	server.GET("/", handlers.HandleRandomSign(searchService, c))
	server.GET("/statesubdivision/:statesubdivisionslug", handlers.HandleRandomSignByCounty(searchService, c))
	server.GET("/place/:placeslug", handlers.HandleRandomSignByPlace(searchService, c))
	server.GET("/state/:stateslug", handlers.HandleRandomSignByState(searchService, c))

	adminRoutes := server.Group("admin")
	adminRoutes.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "header:x-api-key",
		Validator: func(ctx *echo.Context, key string, source middleware.ExtractorSource) (bool, error) {
			return key == c.AdminApiKey, nil
		},
	}))
	adminRoutes.POST("/refresh", handlers.HandleRefreshCounts(counterService, backgroundChan))

}
