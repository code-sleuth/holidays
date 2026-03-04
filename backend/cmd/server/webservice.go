package main

import (
	"net/http"
	"strings"

	"holidays/internal/api/holiday"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
)

func buildWebService(db *mongo.Client, dbName string) *echo.Echo {
	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	origins := make([]string, 0)
	origins = append(origins, "http://localhost:8088")
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		UnsafeAllowOriginFunc: func(c *echo.Context, origin string) (string, bool, error) {
			if origin == "http://localhost" || strings.HasPrefix(origin, "http://localhost:") {
				return origin, true, nil
			}
			for _, allowed := range origins {
				if origin == allowed {
					return origin, true, nil
				}
			}
			return "", false, nil
		},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{
			echo.HeaderContentType,
			echo.HeaderAccept,
		},
	}))

	logger := log.With().Str("handler", "holidays").Logger()

	// handle manager business logic
	handler := holiday.NewManager(holiday.Handler{
		DB:           db,
		DatabaseName: dbName,
		Logger:       logger,
	})

	api := e.Group("/api/v1")

	api.GET("/holidays", handler.GetHolidays)
	api.POST("/lookups", handler.CreateLookUp)
	api.GET("/lookups", handler.GetLookups)

	return e
}
