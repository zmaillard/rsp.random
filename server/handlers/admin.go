package handlers

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"rsp.random/services"
)

func HandleRefreshCounts(counterService services.CounterService, backgroundChannel chan services.UpdateCounterProcess) echo.HandlerFunc {
	return func(c *echo.Context) error {
		backgroundChannel <- counterService.UpdateData
		return c.NoContent(http.StatusAccepted)
	}
}
