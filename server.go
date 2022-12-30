package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Phonlakid/assessment/controller"
	"github.com/Phonlakid/assessment/db"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func main() {
	db.Connect()

	defer db.Conn.Close()

	e := echo.New()
	e.Logger.SetLevel(log.INFO)
	e.POST("/expenses", controller.CreateexpensesHandler)
	e.GET("/expenses/:id", controller.GetUserHandler)
	go func() {
		if err := e.Start(":2565"); err != nil && err != http.ErrServerClosed { // Start server
			e.Logger.Fatal("shutting down the server")
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
