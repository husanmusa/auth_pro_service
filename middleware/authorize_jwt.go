package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/husanmusa/auth_pro_service/utils"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

// AuthorizeJWT -> to authorize JWT Token
func AuthorizeJWT() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		const BearerSchema string = "Bearer "
		authHeader := ctx.Get("Authorization")
		if authHeader == "" {
			return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "No Authorization header found",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, BearerSchema)

		token, err := utils.ValidateToken(tokenString)
		if err != nil {
			fmt.Println("token", tokenString, err.Error())
			return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "Not Valid Token",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return ctx.Status(http.StatusUnauthorized).SendString("Invalid token claims")
		}

		if token.Valid {
			ctx.Locals("userID", claims["userID"])
		} else {
			return ctx.Status(http.StatusUnauthorized).SendString("Token is not valid")
		}

		return ctx.Next()
	}
}
