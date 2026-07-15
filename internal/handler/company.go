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
