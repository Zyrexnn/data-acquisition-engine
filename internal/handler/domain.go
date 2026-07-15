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

type DomainRequest struct {
	Domain string `json:"domain" example:"paper.id"`
}

// Extract godoc
// @Summary Extract domain intelligence
// @Description Mengambil informasi domain melalui protokol RDAP, termasuk registrar, tanggal registrasi, kedaluwarsa, status, dan nameserver.
// @Tags Domain
// @Accept json
// @Produce json
// @Param request body handler.DomainRequest true "Domain yang akan diekstrak"
// @Success 200 {object} response.APIResponse{data=service.DomainData}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /extract/domain [post]
func (h *DomainHandler) Extract(c *fiber.Ctx) error {
	var body DomainRequest

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
