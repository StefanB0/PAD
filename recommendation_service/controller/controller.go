package controller

import (
	"padrecommendations/service"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	rs *service.RecommendationService
}

func NewController(service *service.RecommendationService) *Controller {
	return &Controller{
		rs: service,
	}
}

func (c *Controller) Run() {
	app := fiber.New()

	app.Get("/getTags", c.getTags)

	app.Get("/recommendations/:tagname", func(ctx *fiber.Ctx) error {
		tagname := ctx.Params("tagname")
		imageID, err := c.rs.GetRecommendations(tagname)
		if err != nil {
			return ctx.Status(404).SendString(err.Error())
		}

		return ctx.Status(200).JSON(fiber.Map{
			"imageID": imageID,
		})

	})

	app.Listen(":8083")
}

func (c *Controller) getTags(ctx *fiber.Ctx) error {
	tags := c.rs.GetTags()

	return ctx.Status(200).JSON(fiber.Map{
		"tags": tags,
	})
}


