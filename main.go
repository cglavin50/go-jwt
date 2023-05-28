package main

import (
	"os"

	"github.com/cglavin50/go-jwt/initializers"
	"github.com/gofiber/fiber/v2"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.InitDB(os.Getenv("DSN"))
} // called upon main instantiation

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong!")
	})

	app.Listen(":3000")
}
