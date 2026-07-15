package main

import (
	"log"

	_ "data-acquisition-engine/docs"
	"data-acquisition-engine/internal/handler"
	"data-acquisition-engine/internal/response"
	"data-acquisition-engine/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	fiberswagger "github.com/swaggo/fiber-swagger"
)

// @title Data Acquisition Engine API
// @version 1.0
// @description API untuk ekstraksi metadata website, domain intelligence, dan lokasi geografis.
// @host localhost:8080
// @BasePath /
func main() {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return response.InternalError(c, err.Error())
		},
	})

	app.Use(logger.New())
	app.Use(recover.New())

	websiteSvc := service.NewWebsiteService()
	domainSvc := service.NewDomainService()
	locationSvc := service.NewLocationService()
	companySvc := service.NewCompanyService()

	websiteH := handler.NewWebsiteHandler(websiteSvc)
	domainH := handler.NewDomainHandler(domainSvc)
	locationH := handler.NewLocationHandler(locationSvc)
	companyH := handler.NewCompanyHandler(companySvc)

	app.Get("/swagger/*", fiberswagger.WrapHandler)

	api := app.Group("/extract")
	api.Post("/website", websiteH.Extract)
	api.Post("/domain", domainH.Extract)
	api.Post("/location", locationH.Find)

	app.Get("/company-information", companyH.GetInfo)

	log.Fatal(app.Listen(":8080"))
}
