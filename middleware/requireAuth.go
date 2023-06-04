package middleware

import (
	"fmt"
	"strconv"
	"time"

	"github.com/cglavin50/go-jwt/controllers"
	"github.com/cglavin50/go-jwt/initializers"
	"github.com/cglavin50/go-jwt/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(c *fiber.Ctx) error {
	// get cookie from request
	tokenString := c.Cookies("Authorization")
	if tokenString == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Error, incorrect cookies provided")
	}
	// decode/validate
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return controllers.Pub_key, nil
	})
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Error, cookies could not be parsed")
	}
	// check expiration
	exp := (int64)(claims["exp"].(float64))
	if exp < time.Now().Unix() {
		return fiber.NewError(fiber.StatusUnauthorized, "Error, expired cookie provided")
	}
	// find user with token sub field
	var user models.User
	result := initializers.DB.First(&user, claims["sub"])
	if result.Error != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Error, no user found")
	}
	// attach to request
	id := int(user.ID) // must be cleaner way to convert uint to int without ParseUInt
	fmt.Println("user ID in requireAuth:", id)
	c.Set("id", strconv.Itoa(id)) // avoid uint clash with uint64 conv
	// continue
	c.Next()
	return nil
}
