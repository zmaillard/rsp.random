package handlers

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
	"rsp.random/config"
	"rsp.random/services"
)

func HandleRandomSign(search services.SearchService, cfg *config.Config) echo.HandlerFunc {
	return func(c *echo.Context) error {
		res, err := search.RandomSign()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "unable to retrieve result")
		}

		redirectUrl, err := res.GetRedirectUrl(cfg)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "unable to retrieve result")
		}

		return c.Redirect(http.StatusTemporaryRedirect, redirectUrl)
	}

}

func HandleRandomSignByCounty(search services.SearchService, cfg *config.Config) echo.HandlerFunc {
	type HandleRandomSignByCountyDto struct {
		CountySlug string `param:"statesubdivisionslug" validate:"required"`
	}
	return func(c *echo.Context) error {
		var body HandleRandomSignByCountyDto
		if err := c.Bind(&body); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid County")
		}
		if err := c.Validate(body); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		tokens := strings.Split(body.CountySlug, "_")
		if len(tokens) != 2 {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid county")
		}

		res, err := search.RandomSignByCounty(tokens[0], tokens[1])
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "unable to retrieve result")
		}

		redirectUrl, err := res.GetRedirectUrl(cfg)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "unable to retrieve result")
		}

		return c.Redirect(http.StatusTemporaryRedirect, redirectUrl)
	}

}
func HandleRandomSignByPlace(search services.SearchService, cfg *config.Config) echo.HandlerFunc {
	type HandleRandomSignByPlaceDto struct {
		PlaceSlug string `param:"placeslug" validate:"required"`
	}
	return func(c *echo.Context) error {
		var body HandleRandomSignByPlaceDto
		if err := c.Bind(&body); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid Place")
		}
		if err := c.Validate(body); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		tokens := strings.Split(body.PlaceSlug, "_")
		if len(tokens) != 2 {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid place")
		}

		res, err := search.RandomSignByPlace(tokens[0], tokens[1])
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "unable to retrieve result")
		}

		redirectUrl, err := res.GetRedirectUrl(cfg)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "unable to retrieve result")
		}

		return c.Redirect(http.StatusTemporaryRedirect, redirectUrl)
	}

}

func HandleRandomSignByState(search services.SearchService, cfg *config.Config) echo.HandlerFunc {
	type HandleRandomSignByStateDto struct {
		StateSlug string `param:"stateslug" validate:"required"`
	}
	return func(c *echo.Context) error {
		var body HandleRandomSignByStateDto
		if err := c.Bind(&body); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid State")
		}
		if err := c.Validate(body); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		res, err := search.RandomSignByState(body.StateSlug)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "unable to retrieve result")
		}

		redirectUrl, err := res.GetRedirectUrl(cfg)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "unable to retrieve result")
		}

		return c.Redirect(http.StatusTemporaryRedirect, redirectUrl)
	}

}
