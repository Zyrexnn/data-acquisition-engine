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

type LocationRequest struct {
	Query string `json:"query" example:"Paper.id Jakarta"`
}

// Find godoc
// @Summary Find location by query
// @Description Mencari informasi lokasi geografis berdasarkan query teks menggunakan OpenStreetMap Nominatim.
// @Tags Location
// @Accept json
// @Produce json
// @Param request body handler.LocationRequest true "Query pencarian lokasi"
// @Success 200 {object} response.APIResponse{data=service.LocationData}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /extract/location [post]
func (h *LocationHandler) Find(c *fiber.Ctx) error {
	var body LocationRequest

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
