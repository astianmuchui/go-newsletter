package main

import (
	"github.com/astianmuchui/go-newsletter/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/django/v3"
)

func main () {

	engine := django.New("./views",".django")
	app := fiber.New(fiber.Config{
		Prefork: true,
		Views: engine,
	})

	/* Use logger and recover middleware */
	app.Use(recover.New())
	app.Use(logger.New())


	app.Get("/", handlers.HomeHandler)
	app.Post("/subscribe", handlers.SubscribeHandler)

	app.Get("/send", handlers.SendHandler)
	app.Post("/send-emails", handlers.SendEmailHandler)
	app.Listen(":8080")
}