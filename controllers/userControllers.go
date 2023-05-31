package controllers

import (
	"fmt"
	"os"
	"time"

	"github.com/cglavin50/go-jwt/initializers"
	"github.com/cglavin50/go-jwt/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *fiber.Ctx) error {
	user := models.User{}
	extractUser(c, &user)
	result := initializers.DB.Create(&user) // create record (post) user following user model
	if result.Error != nil {
		c.Status(400).JSON(&fiber.Map{
			"error": "Failed to create user, email already in use",
		}) // sends HTTP status 400, bad request
		return nil
	}

	// if no fails, set status ok and return
	c.Status(200).JSON(&fiber.Map{
		"Status":   "OK",
		"Password": user.Password,
	})
	return nil
}

func Login(c *fiber.Ctx) error {
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
	var user models.User
	result := initializers.DB.First(&user, "email = ?", body.Email) // order by ID, find match on email
	if result.Error != nil {
		c.Status(400).JSON(&fiber.Map{
			"error": "Bad request",
		}) // sends HTTP status 400, bad request
		return nil
	}
	// check reqUser vs User
	// $2a$10$J/FMfxr22jKobnC4wkNtluGggdR1iwzn2QtS7.eM.qxtZKTwL0jUa"
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil { // check match of hashes, relying on strong collision resistance
		bp, _ := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
		c.Status(400).JSON(&fiber.Map{
			"error":           "Incorrect credentials provided",
			"hashed password": bp,
			"user email":      user.Email,
			"user password":   user.Password,
		})
		return nil
	} // wrong-password case

	// generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"iss": "go-jwt-backend",
		"exp": time.Now().Add(time.Hour).Unix(),
		"sub": user.ID, // is this safe to publish?
	})
	fmt.Println(os.Getenv("ECDSA_PRV"))
	tokenString, err := token.SignedString([]byte(os.Getenv("ECDSA_PRV")))
	if err != nil {
		c.Status(400).JSON(&fiber.Map{
			"error": "Failed to create token",
		}) // sends HTTP status 400, bad request
		return nil
	}
	// send JWT
	c.Status(200).JSON(&fiber.Map{
		"token": tokenString,
	})
	return nil
}

func extractUser(c *fiber.Ctx, user *models.User) {
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
		return
	}

	// hash password
	hpw, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.Status(400).JSON(&fiber.Map{
			"error": "Failed to hash password",
		}) // sends HTTP status 400, bad request
		return
	}

	// Populate User field
	user.Email = body.Email
	user.Password = string(hpw)
}
