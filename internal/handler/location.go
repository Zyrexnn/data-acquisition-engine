package handler

import (
	"data-acquisition-engine/internal/response"
	"data-acquisition-engine/internal/service"

	"github.com/gofiber/fiber/v2"
)

type LocationHandler struct {
	svc *service.LocationService
}

func NewLocationHandler(svc *service.LocationService) *LocationHandler {
	return &LocationHandler{svc: svc}
}

func (h *LocationHandler) Find(c *fiber.Ctx) error {
	var body struct {
		Query string `json:"query"`
	}

	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	if body.Query == "" {
		return response.BadRequest(c, "query is required")
	}

	data, err := h.svc.Find(body.Query)
	if err != nil {
		return response.InternalError(c, "failed to find location")
	}

	return response.Success(c, data)
}
