package handler

import (
	"data-acquisition-engine/internal/response"
	"data-acquisition-engine/internal/service"

	"github.com/gofiber/fiber/v2"
)

type CompanyHandler struct {
	svc *service.CompanyService
}

func NewCompanyHandler(svc *service.CompanyService) *CompanyHandler {
	return &CompanyHandler{svc: svc}
}

// GetInfo godoc
// @Summary Get unified company information
// @Description Menggabungkan data website, domain, dan lokasi secara paralel berdasarkan domain perusahaan. Mendukung fallback lokasi jika title website tidak menghasilkan hasil.
// @Tags Company
// @Accept json
// @Produce json
// @Param domain query string true "Domain perusahaan (contoh: paper.id)"
// @Success 200 {object} response.APIResponse{data=service.CompanyData}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /company-information [get]
func (h *CompanyHandler) GetInfo(c *fiber.Ctx) error {
	domain := service.CleanDomain(c.Query("domain"))

	if domain == "" {
		return response.BadRequest(c, "domain query param is required")
	}

	data, err := h.svc.GetInfo(domain)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.Success(c, data)
}
