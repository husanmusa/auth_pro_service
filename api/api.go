package api

import (
	"encoding/json"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/casbin/casbin/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/husanmusa/auth_pro_service/api/controller"
	_ "github.com/husanmusa/auth_pro_service/api/docs"
	"github.com/husanmusa/auth_pro_service/grpc/client"
	"github.com/husanmusa/auth_pro_service/middleware"
)

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// SetUpRouter godoc
// @description This is a api gateway
// @termsOfService https://udevs.io
func SetUpRouter(svc client.ServiceManagerI, enforcer *casbin.Enforcer) *fiber.App {

	router := fiber.New(fiber.Config{JSONEncoder: json.Marshal, BodyLimit: 100 * 1024 * 1024})
	router.Use(logger.New(), cors.New())

	router.Get("/api/swagger/*", swagger.HandlerDefault)
	r := router.Group("/api")

	userController := controller.NewUserController(svc)

	r.Post("/user/register", userController.AddUser(enforcer))
	r.Post("/user/signin", userController.SignInUser)

	r.Use(middleware.AuthorizeJWT())

	// APPOINTMENT
	r.Get("/user", middleware.Authorize("userGet", "read", enforcer), userController.GetAllUser)
	r.Get("/user/:user_id", middleware.Authorize("userGet", "read", enforcer), userController.GetUser)
	r.Put("/user/:user_id", middleware.Authorize("userUpt", "write", enforcer), userController.UpdateUser)
	r.Delete("/user/:user_id", middleware.Authorize("userUpt", "write", enforcer), userController.DeleteUser)
	return router
}
