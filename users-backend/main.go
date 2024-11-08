package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"
	"users-backend/controller"
	"users-backend/handler"
	"users-backend/repo/postgres"

	_ "users-backend/docs"

	"github.com/labstack/echo/v4"
)

func main() {
	repo, cleanup := postgres.NewPostgresRepo()
	defer cleanup()

	c := controller.NewUserController(repo)

	e := echo.New()
	handler.InitRouter(e, c)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
