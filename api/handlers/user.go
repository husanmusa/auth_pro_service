package handlers

import (
	"github.com/casbin/casbin/v2"
	"github.com/gofiber/fiber/v2"
	pb "github.com/husanmusa/auth_pro_service/genproto/auth_service"
	"github.com/husanmusa/auth_pro_service/pkg/utils"
	"net/http"
	"strconv"
)

// GetAllUser godoc
// @Security ApiKeyAuth
// @Summary Get users
// @Description This API for getting users
// ID get_all_user
// @Tags User
// @Accept json
// @Produce json
// @Param offset query integer false "offset"
// @Param limit query integer false "limit"
// Success 201 {object} http.Response{data=string} "User data"
// @Failure 400 {object} http.Response{data=string} "Bad request"
// @Failure 500 {object} http.Response{data=string} "Internal server error"
// @Router /api/user/ [GET]
func (h Handler) GetAllUser(ctx *fiber.Ctx) error {
	limit, err := getLimitParam(ctx)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	offset, err := getOffsetParam(ctx)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	users, err := h.services.UserService().GetUserList(ctx.Context(), &pb.GetUserListRequest{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})

	}
	return ctx.Status(http.StatusOK).JSON(users)
}

// GetUser godoc
// @Security ApiKeyAuth
// @Summary Get user by user_id
// ID get_user
// @Tags User
// @Accept json
// @Produce json
// @Param user_id path string true "user_id"
// @Success 200 {object} http.Response{data=auth_service.User} "GetUser ResponseBody"
// @Failure 400 {object} http.Response{data=string} "Bad request"
// @Failure 500 {object} http.Response{data=string} "Internal server error"
// @Router /api/user/{user_id} [GET]
func (h Handler) GetUser(ctx *fiber.Ctx) error {
	id := ctx.Params("user")

	user, err := h.services.UserService().GetUser(ctx.Context(), &pb.ById{Id: id})
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(http.StatusOK).JSON(user)
}

// SignInUser godoc
// @Summary sign in a user
// @Description This API for sign in a user
// ID create_user
// @Tags User
// @Accept json
// @Produce json
// @Param user body auth_service.SignInReq true "SignInReq"
// Success 201 {object} http.Response{data=string} "User data"
// @Failure 400 {object} http.Response{data=string} "Bad request"
// @Failure 500 {object} http.Response{data=string} "Internal server error"
// @Router /api/user/signin [POST]
func (h Handler) SignInUser(ctx *fiber.Ctx) error {
	var user pb.SignInReq
	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	token, err := h.services.UserService().SignInUser(ctx.Context(), &user)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"msg": "No Such User Found"})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{"token": token})
}

// AddUser godoc
// @Summary Create a new user
// @Description This API for creating a new user
// ID create_user
// @Tags User
// @Accept json
// @Produce json
// @Param user body auth_service.User true "UserCreateRequest"
// Success 201 {object} http.Response{data=string} "User data"
// @Failure 400 {object} http.Response{data=string} "Bad request"
// @Failure 500 {object} http.Response{data=string} "Internal server error"
// @Router /api/user/register [POST]
func (h Handler) AddUser(enforcer *casbin.Enforcer) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var user pb.User
		if err := ctx.BodyParser(&user); err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		utils.HashPassword(&user.Password)
		_, err := h.services.UserService().CreateUser(ctx.Context(), &user)
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		//_, err = enforcer.AddGroupingPolicy(fmt.Sprint(user.Id), user.Role)
		//if err != nil {
		//	return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		//}

		return ctx.Status(http.StatusOK).JSON("SUCCESS")
	}
}

// UpdateUser godoc
// @Security ApiKeyAuth
// @Summary Update user by_id
// @ID update_user
// @Tags User
// @Accept json
// @Produce json
// @Param user body auth_service.User true "UserUpdateRequest"
// @Success 200 {object} http.Response{data=string} "Success Update"
// @Failure 400 {object} http.Response{data=string} "Bad request"
// @Failure 500 {object} http.Response{data=string} "Internal server error"
// @Router /api/user/{user_id} [PUT]
func (h Handler) UpdateUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	_, err := h.services.UserService().UpdateUser(ctx.Context(), &pb.User{Id: id})
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(http.StatusOK).JSON("SUCCESS")
}

// DeleteUser godoc
// @Security ApiKeyAuth
// @Summary Delete user by_id
// @ID delete_user
// @Tags User
// @Accept json
// @Produce json
// @Param user_id path string false "user_id"
// @Success 200 {object} http.Response{data=string} "Success DeleteUser"
// @Failure 400 {object} http.Response{data=string} "Bad request"
// @Failure 500 {object} http.Response{data=string} "Internal server error"
// @Router /api/user/{user_id} [DELETE]
func (h Handler) DeleteUser(ctx *fiber.Ctx) error {

	id := ctx.Params("user_id")

	_, err := h.services.UserService().DeleteUser(ctx.Context(), &pb.ById{Id: id})
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(http.StatusOK).JSON("SUCCESS")
}

func getOffsetParam(c *fiber.Ctx) (offset int, err error) {
	offsetStr := c.Query("offset", "0")
	return strconv.Atoi(offsetStr)
}

func getLimitParam(c *fiber.Ctx) (offset int, err error) {
	offsetStr := c.Query("limit", "10")
	return strconv.Atoi(offsetStr)
}
