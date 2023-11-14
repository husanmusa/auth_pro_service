package middleware

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/husanmusa/auth_pro_service/utils"
	"net/http"
	"strings"

	"github.com/casbin/casbin/v2"
)

// Authorize determines if current user has been authorized to take an action on an object.
func Authorize(obj string, act string, enforcer *casbin.Enforcer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get current user/subject
		const BearerSchema string = "Bearer "
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "No Authorization header found",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, BearerSchema)

		token, err := utils.ValidateToken(tokenString)
		if err != nil {
			fmt.Println("token", tokenString, err.Error())
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "Not Valid Token",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(http.StatusUnauthorized).SendString("Invalid token claims")
		}
		role := claims["role"]
		// Load policy from Database
		err = enforcer.LoadPolicy()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"msg": "Failed to load policy from DB"})
		}

		// Casbin enforces policy
		ok, err = enforcer.Enforce(role, obj, act)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"msg": "Error occurred when authorizing user"})
		}

		if !ok {
			return c.Status(403).JSON(fiber.Map{"msg": "You are not authorized"})
		}

		return c.Next()
	}
}
