package controllers

import (
	"github.com/cglavin50/go-jwt/initializers"
	"github.com/cglavin50/go-jwt/models"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *fiber.Ctx) error {
	// req.body: email/pass pair
	var body struct {
		Email    string
		Password string
	}
	err := c.BodyParser(&body)
	if err != nil {
		c.Status(400).JSON(&fiber.Map{
			"error": "Failed to read body", // => HTTP.Status 400 = bad requests
		}) // sends HTTP status 400, bad request
		return nil
	}

	// hash password
	hpw, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.Status(400).JSON(&fiber.Map{
			"error": "Failed to hash password",
		}) // sends HTTP status 400, bad request
		return nil
	}

	// create User
	user := models.User{
		Email:    body.Email,
		Password: string(hpw),
	}
	result := initializers.DB.Create(&user) // create record (post) user following user model
	if result.Error != nil {
		c.Status(400).JSON(&fiber.Map{
			"error": "Failed to create user",
		}) // sends HTTP status 400, bad request
		return nil
	}

	// if no fails, set status ok and return
	c.Status(200).JSON(&fiber.Map{
		"Status": "OK",
	})
	return nil
}
