package api

import (
	"encoding/json"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/casbin/casbin/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/husanmusa/auth_pro_service/api/docs"
	"github.com/husanmusa/auth_pro_service/api/handlers"
	middleware2 "github.com/husanmusa/auth_pro_service/api/middleware"
)

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// SetUpRouter godoc
// @description This is a api gateway
// @termsOfService https://udevs.io
func SetUpRouter(h handlers.Handler, enforcer *casbin.Enforcer) *fiber.App {

	router := fiber.New(fiber.Config{JSONEncoder: json.Marshal, BodyLimit: 100 * 1024 * 1024})
	router.Use(logger.New(), cors.New())

	router.Get("/api/swagger/*", swagger.HandlerDefault)
	r := router.Group("/api")

	r.Post("/user/register", h.AddUser(enforcer))
	r.Post("/user/signin", h.SignInUser)

	r.Use(middleware2.AuthorizeJWT())

	// APPOINTMENT
	r.Get("/user", middleware2.Authorize("userGet", "read", enforcer), h.GetAllUser)
	r.Get("/user/:user_id", middleware2.Authorize("userGet", "read", enforcer), h.GetUser)
	r.Put("/user/:user_id", middleware2.Authorize("userUpt", "write", enforcer), h.UpdateUser)
	r.Delete("/user/:user_id", middleware2.Authorize("userUpt", "write", enforcer), h.DeleteUser)
	return router
}
