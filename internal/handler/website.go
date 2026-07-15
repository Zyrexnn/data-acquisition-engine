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

type WebsiteRequest struct {
	URL string `json:"url" example:"https://paper.id"`
}

// Extract godoc
// @Summary Extract website metadata
// @Description Mengekstrak metadata dari URL website, termasuk title, deskripsi, Open Graph, email, telepon, dan media sosial.
// @Tags Website
// @Accept json
// @Produce json
// @Param request body handler.WebsiteRequest true "URL website yang akan diekstrak"
// @Success 200 {object} response.APIResponse{data=service.WebsiteData}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /extract/website [post]
func (h *WebsiteHandler) Extract(c *fiber.Ctx) error {
	var body WebsiteRequest

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
	