package handler

import (
	"data-acquisition-engine/internal/response"
	"data-acquisition-engine/internal/service"

	"github.com/gofiber/fiber/v2"
)

type WebsiteHandler struct {
	svc *service.WebsiteService
}

func NewWebsiteHandler(svc *service.WebsiteService) *WebsiteHandler {
	return &WebsiteHandler{svc: svc}
}

func (h *WebsiteHandler) Extract(c *fiber.Ctx) error {
	var body struct {
		URL string `json:"url"`
	}

	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	if body.URL == "" {
		return response.BadRequest(c, "url is required")
	}

	data, err := h.svc.Extract(body.URL)
	if err != nil {
		return response.InternalError(c, "failed to extract website metadata")
	}

	return response.Success(c, data)
}
