package handler

import (
	"data-acquisition-engine/internal/response"
	"data-acquisition-engine/internal/service"

	"github.com/gofiber/fiber/v2"
)

type DomainHandler struct {
	svc *service.DomainService
}

func NewDomainHandler(svc *service.DomainService) *DomainHandler {
	return &DomainHandler{svc: svc}
}

func (h *DomainHandler) Extract(c *fiber.Ctx) error {
	var body struct {
		Domain string `json:"domain"`
	}

	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	if body.Domain == "" {
		return response.BadRequest(c, "domain is required")
	}

	data, err := h.svc.Extract(body.Domain)
	if err != nil {
		return response.InternalError(c, "failed to extract domain intelligence")
	}

	return response.Success(c, data)
}
