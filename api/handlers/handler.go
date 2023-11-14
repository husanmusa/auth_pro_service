package handlers

import (
	"github.com/husanmusa/auth_pro_service/api/http"
	"github.com/husanmusa/auth_pro_service/config"
	"github.com/husanmusa/auth_pro_service/grpc/client"
	"strconv"

	"github.com/saidamir98/udevs_pkg/logger"

	fiber "github.com/gofiber/fiber/v2"
)

type Handler struct {
	cfg      config.Config
	log      logger.LoggerI
	services client.ServiceManagerI
}

func NewHandler(cfg config.Config, log logger.LoggerI, svcs client.ServiceManagerI) Handler {
	return Handler{
		cfg:      cfg,
		log:      log,
		services: svcs,
	}
}

func (h *Handler) handleResponse(c *fiber.Ctx, status http.Status, data interface{}) error {
	switch code := status.Code; {
	case code < 300:
		h.log.Info(
			"response",
			logger.Int("code", status.Code),
			logger.String("status", status.Status),
			logger.Any("description", status.Description),
			logger.Any("data", data),
		)
	case code < 400:
		h.log.Warn(
			"response",
			logger.Int("code", status.Code),
			logger.String("status", status.Status),
			logger.Any("description", status.Description),
			logger.Any("data", data),
		)
	default:
		h.log.Error(
			"response",
			logger.Int("code", status.Code),
			logger.String("status", status.Status),
			logger.Any("description", status.Description),
			logger.Any("data", data),
		)
	}

	return c.Status(status.Code).JSON(data)
}

func (h *Handler) getOffsetParam(c *fiber.Ctx) (offset int, err error) {
	offsetStr := c.Query("offset", h.cfg.DefaultOffset)
	return strconv.Atoi(offsetStr)
}

func (h *Handler) getLimitParam(c *fiber.Ctx) (offset int, err error) {
	offsetStr := c.Query("limit", h.cfg.DefaultLimit)
	return strconv.Atoi(offsetStr)
}
