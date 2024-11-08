package handler

import (
	"fmt"
	"net/http"
	"users-backend/controller"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// Swagger spec:
//
//	@title			User Service
//	@version		1.0
//	@description	User Management Service
//	@host			localhost:8080
//	@BasePath		/api/v1
func InitRouter(e *echo.Echo, userController *controller.UserControllerImpl) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:4200"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	echo.NotFoundHandler = func(c echo.Context) error {
		return respError(c, http.StatusNotFound, "Invalid endpoint", fmt.Sprintf("Endpoint %s does not exist", c.Request().URL.Path))
	}

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	api := e.Group("/api/v1")
	user := api.Group("/users")

	userHttpHandler := NewUserHttpHandler(user, userController)
	userHttpHandler.RegisterRoutes()
}
