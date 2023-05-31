package main

import (
	"log"
	"os"

	"github.com/cglavin50/go-jwt/controllers"
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
	app.Post("/signup", controllers.SignUp) // validate syntax here (vs passing parameter to this func and handling error elsewhere)
	// ^ handles post requests to create a new user, request: header: POST body: email, password (handles any type of encoding)
	app.Post("/login", controllers.Login)

	log.Fatal(app.Listen(":3000")) //listening on localhost:3000
}
