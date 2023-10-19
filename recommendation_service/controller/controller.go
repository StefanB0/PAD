package controller

import (
	"padrecommendations/models"
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

	app.Get("/status", c.status)

	app.Post("/getTags", c.getTags)
	app.Post("/getRecommendations", c.getRecommendations)
	app.Post("/addImage", c.addImage)
	app.Post("/updateImage", c.updateImage)
	app.Post("/deleteALL", c.deleteAll)

	app.Listen(":8083")
}

func (c *Controller) status(ctx *fiber.Ctx) error {
	return ctx.Status(200).SendString("OK")
}

func (c *Controller) getTags(ctx *fiber.Ctx) error {
	tags := c.rs.GetTags()

	return ctx.Status(200).JSON(fiber.Map{
		"tags": tags,
	})
}

func (c *Controller) getRecommendations(ctx *fiber.Ctx) error {
	req := new(recommendRequest)
	err := ctx.BodyParser(req)
	if err != nil {
		return ctx.Status(400).SendString("Bad request")
	}

	imageID, err := c.rs.GetRecommendations(req.Tag)
	if err != nil {
		return ctx.Status(404).SendString(err.Error())
	}

	return ctx.Status(200).JSON(fiber.Map{
		"imageID": imageID,
	})
}

func (c *Controller) addImage(ctx *fiber.Ctx) error {
	req := new(addImageRequest)
	err := ctx.BodyParser(req)
	if err != nil {
		return ctx.Status(400).SendString("Bad request")
	}

	image := models.Image{
		ImageID: req.ID,
		Tags:    req.Tags,
	}
	c.rs.AddImage(image)

	return ctx.Status(201).SendString("Image added")
}

func (c *Controller) updateImage(ctx *fiber.Ctx) error {
	req := new(updateImageRequest)
	err := ctx.BodyParser(req)
	if err != nil {
		return ctx.Status(400).SendString("Bad request")
	}

	c.rs.AddView(req.ID, req.Views)
	c.rs.AddLike(req.ID, req.Likes)

	return ctx.Status(200).SendString("Image updated")
}

func (c *Controller) deleteAll(ctx *fiber.Ctx) error {
	c.rs.DeleteAll()

	return ctx.Status(200).SendString("All data deleted")
}
