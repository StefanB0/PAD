package controller

import (
	"fmt"
	"padrecommendations/models"
	"padrecommendations/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type Controller struct {
	rs        *service.RecommendationService
	semaphore *service.Semaphore
}

func NewController(rs *service.RecommendationService) *Controller {
	return &Controller{
		rs:        rs,
		semaphore: service.NewSemaphore(20),
	}
}

func (c *Controller) Run(port string) {
	app := fiber.New()

	app.Get("/status", c.status)

	app.Post("/getTags", c.getTags)
	app.Post("/getRecommendations", c.getRecommendations)
	app.Post("/addImage", c.addImage)
	app.Post("/updateImage", c.updateImage)
	app.Post("/deleteALL", c.deleteAll)

	app.Delete("/transaction/:id", c.revertTransaction)

	app.Listen(":8083")
}

func (c *Controller) status(ctx *fiber.Ctx) error {
	return ctx.Status(200).SendString("OK")
}

func (c *Controller) getTags(ctx *fiber.Ctx) error {
	c.semaphore.Acquire()
	defer c.semaphore.Release()

	tags := c.rs.GetTags()

	return ctx.Status(200).JSON(fiber.Map{
		"tags": tags,
	})
}

func (c *Controller) getRecommendations(ctx *fiber.Ctx) error {
	c.semaphore.Acquire()
	defer c.semaphore.Release()

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
	c.semaphore.Acquire()
	defer c.semaphore.Release()

	req := new(addImageRequest)
	err := ctx.BodyParser(req)
	if err != nil {
		return ctx.Status(400).SendString("Bad request")
	}

	image := models.Image{
		ImageID: req.ID,
		Tags:    req.Tags,
	}

	err = c.rs.AddImage(image)
	if err != nil {
		c.rs.CancelSagaTransaction(req.SagaID)
		return ctx.Status(500).SendString("Server Error")
	}

	log.Info().Msg("Added image " + strconv.Itoa(req.ID) + " with tags " + fmt.Sprintf("%v", req.Tags))
	c.rs.ConfirmSagaTransaction(req.SagaID)

	return ctx.Status(201).SendString("Image added")
}

func (c *Controller) updateImage(ctx *fiber.Ctx) error {
	c.semaphore.Acquire()
	defer c.semaphore.Release()

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
	c.semaphore.Acquire()
	defer c.semaphore.Release()

	c.rs.DeleteAll()

	return ctx.Status(200).SendString("All data deleted")
}

func (c *Controller) revertTransaction(ctx *fiber.Ctx) error {
	c.semaphore.Acquire()
	defer c.semaphore.Release()

	id := ctx.Params("id")

	err := c.rs.RevertSagaTransaction(id)
	if err != nil {
		return ctx.Status(404).SendString(err.Error())
	}

	log.Info().Msg(fmt.Sprintf("Transaction %s reverted", id))

	return ctx.Status(200).SendString("Transaction reverted")
}
