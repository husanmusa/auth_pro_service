package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"

	"github.com/casbin/casbin/v2"
)

// Authorize determines if current user has been authorized to take an action on an object.
func Authorize(obj string, act string, enforcer *casbin.Enforcer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get current user/subject
		sub := c.Get("user_id")

		// Load policy from Database
		err := enforcer.LoadPolicy()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"msg": "Failed to load policy from DB"})
		}

		// Casbin enforces policy
		ok, err := enforcer.Enforce(fmt.Sprint(sub), obj, act)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"msg": "Error occurred when authorizing user"})
		}

		if !ok {
			return c.Status(403).JSON(fiber.Map{"msg": "You are not authorized"})
		}

		return c.Next()
	}
}
