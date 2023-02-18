package main

import (
	fiber "github.com/gofiber/fiber/v2"
	"jwtauth/service"
	"net/http"
)

func main() {

	app := fiber.New()

	app.Use("/register", func(ctx *fiber.Ctx) error {
		token := service.GenerateToken(ctx)

		if token == nil {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		return ctx.JSON(token)

	})

	app.Use(func(ctx *fiber.Ctx) error {
		isOk := service.CheckToken(ctx)
		if isOk {
			return ctx.Next()
		}
		return ctx.SendStatus(401)
	})

	app.Get("/home", func(ctx *fiber.Ctx) error {
		return ctx.JSON("hello world")
	})

	app.Listen(":8080")
}
