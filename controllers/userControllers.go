package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cglavin50/go-jwt/initializers"
	"github.com/cglavin50/go-jwt/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// var key *ecdsa.PrivateKey

// func FetchKey() {
// 	var err error
// 	ecdsa := os.Getenv("ECDSA_PRV")
// 	fmt.Println("ecdsa:", ecdsa)
// 	key, err = x509.ParseECPrivateKey([]byte(ecdsa))
// 	if err != nil {
// 		log.Fatal("Failed to parse ECDSA key")
// 	}
// }

func SignUp(c *fiber.Ctx) error {
	user := models.User{}
	extractUser(c, &user)
	fmt.Println("Creating user: ", user.Email)
	result := initializers.DB.Create(&user) // create record (post) user following user model
	if result.Error != nil {
		c.Status(400).JSON(&fiber.Map{
			"error": "Failed to create user, email already in use",
		}) // sends HTTP status 400, bad request
		return nil
	}

	// if no fails, set status ok and return
	c.Status(200).JSON(&fiber.Map{
		"Status": "OK",
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
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil { // check match of hashes, relying on strong collision resistance
		c.Status(400).JSON(&fiber.Map{
			"error": "Incorrect credentials provided",
		})
		return nil
	} // wrong-password case

	// generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "go-jwt-backend",
		"exp": time.Now().Add(time.Hour).Unix(),
		"sub": user.ID, // is this safe to publish?
	})

	// update this to use an actual asymmetric key system
	tokenString, err := token.SignedString([]byte("12i3bkajsckl23ekljncoa9sid"))
	if err != nil {
		c.Status(400).JSON(&fiber.Map{
			"error": "Failed to create token",
		}) // sends HTTP status 400, bad request
		return nil
	}
	// send JWT as cookie
	c.Cookie(&fiber.Cookie{
		Name:    "Authorization",
		Value:   tokenString,
		Expires: time.Now().Add(time.Minute * 15),
	})
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

func Validate(c *fiber.Ctx) error {
	fmt.Println("Printing in validate")
	for key, value := range c.GetRespHeaders() {
		fmt.Println(key, ":", value)
	}
	err := c.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "logged in",
		"user":    c.GetRespHeader("Id"),
	})
	return err
}
