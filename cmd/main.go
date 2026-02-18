package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Makefolder/cynero/internal/common"
	"github.com/Makefolder/cynero/internal/handler"
	"github.com/Makefolder/cynero/internal/log"
	"github.com/Makefolder/cynero/internal/routes"
	"github.com/Makefolder/cynero/internal/service"
	"github.com/gofiber/fiber/v3"
)

const (
	env  common.Environment = common.EnvironmentDevelopment
	port uint16             = 7654
)

func main() {
	log, err := log.New(env)
	if err != nil {
		panic(err)
	}

	app := fiber.New()
	api := app.Group("/api")

	s := service.New()
	h := handler.New(s)
	routes.Setup(api, h)

	if err := app.Listen(fmt.Sprintf("0.0.0.0:%d", port)); err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	log.Infoln("Gracefully shutting down...")
	if err := app.Shutdown(); err != nil {
		log.Fatal(err)
	}
}
