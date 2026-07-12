package handlers

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v5"
	"rsp.random/services"
)

func HandleRefreshCounts(counterService services.CounterService, backgroundChannel chan func(context.Context) error) echo.HandlerFunc {
	return func(c *echo.Context) error {
		backgroundChannel <- counterService.UpdateCounts
		return c.NoContent(http.StatusAccepted)
	}

}
