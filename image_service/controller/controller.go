package controller

import (
	// "log"
	"fmt"
	"padimage/models"
	"padimage/service"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type ImageController struct {
	imageService *service.ImageService
	tokenService *service.TokenService
	semaphore    *service.Semaphore
}

func NewImageController(is *service.ImageService, ts *service.TokenService) *ImageController {
	return &ImageController{
		imageService: is,
		tokenService: ts,
		semaphore:    service.NewSemaphore(20),
	}
}

func (c *ImageController) Run(port string) {
	app := fiber.New()

	app.Get("/status", c.status)

	app.Post("/getImage", c.getImage)
	app.Post("/getImageInfo", c.getImageInfo)
	app.Post("/uploadImage", c.upload)
	app.Post("/likeImage", c.likeImage)
	app.Post("/updateImage", c.update)
	app.Post("/deleteImage", c.delete)

	app.Delete("/transaction/:id", c.revertTransaction)

	app.Listen(port)
}

func (c *ImageController) status(ctx *fiber.Ctx) error {
	return ctx.Status(200).SendString("OK")
}

func (c *ImageController) getImage(ctx *fiber.Ctx) error {
	c.semaphore.Acquire()
	defer c.semaphore.Release()

	req := new(getImageRequest)
	err := ctx.BodyParser(req)
	if err != nil {
		return ctx.Status(400).SendString("Bad request")
	}

	image, err := c.imageService.GetImage(req.ImageID)
	if err != nil {
		return ctx.Status(404).SendString(err.Error())
	}

	ctx.Set("Content-Type", "image/jpeg")

	c.imageService.AddViews(req.ImageID, 1)

	return ctx.Status(200).Send(image.ImageChunk)
}

func (c *ImageController) getImageInfo(ctx *fiber.Ctx) error {
	c.semaphore.Acquire()
	defer c.semaphore.Release()

	req := new(getImageRequest)
	err := ctx.BodyParser(req)
	if err != nil {
		return ctx.Status(400).SendString("Bad request")
	}

	image, err := c.imageService.GetImage(req.ImageID)
	if err != nil {
		return ctx.Status(404).SendString(err.Error())
	}

	res := getImageResponse{
		ImageID:     image.ImageID,
		Author:      image.Author,
		Title:       image.Title,
		Description: image.Description,
		Tags:        image.Tags,
	}

	return ctx.Status(200).JSON(res)
}

func (c *ImageController) upload(ctx *fiber.Ctx) error {
	c.semaphore.Acquire()
	defer c.semaphore.Release()

	req := new(uploadRequest)

	form, err := ctx.MultipartForm()
	if err != nil {
		return ctx.Status(400).SendString("Bad request. Missing form data." + err.Error())
	}
	if temp := form.Value["author"]; len(temp) == 0 {
		return ctx.Status(400).SendString("Bad request. Missing author field")
	}
	if temp := form.Value["title"]; len(temp) == 0 {
		return ctx.Status(400).SendString("Bad request. Missing title field")
	}
	if temp := form.Value["description"]; len(temp) == 0 {
		return ctx.Status(400).SendString("Bad request. Missing description field")
	}
	if temp := form.Value["tags"]; len(temp) == 0 {
		return ctx.Status(400).SendString("Bad request. Missing tags field")
	}
	if temp := form.Value["sagaid"]; len(temp) == 0 {
		return ctx.Status(400).SendString("Bad request. Missing Saga ID")
	}

	req.Author = form.Value["author"][0]
	req.Title = form.Value["title"][0]
	req.Description = form.Value["description"][0]
	req.Tags = form.Value["tags"]
	req.SagaID = form.Value["sagaid"][0]

	fileheader, err := ctx.FormFile("image")
	if err != nil {
		return ctx.Status(400).SendString("Bad request. Missing image file")
	}

	file, err := fileheader.Open()
	if err != nil {
		return ctx.Status(400).SendString("Bad request. Could not open image file")
	}

	req.ImageBytes = make([]byte, fileheader.Size)
	_, err = file.Read(req.ImageBytes)
	if err != nil {
		return ctx.Status(400).SendString("Bad request. Could not read image file")
	}

	image := models.Image{
		Author:      req.Author,
		Title:       req.Title,
		Description: req.Description,
		Tags:        req.Tags,
		ImageChunk:  req.ImageBytes,
	}

	id, err := c.imageService.CreateImage(image, req.SagaID, req.Token)

	c.imageService.AddSagaTransaction(req.SagaID, id)

	if err != nil {
		c.imageService.CancelSagaTransaction(req.SagaID)
		return ctx.Status(503).SendString("Service unavailable")
	}

	ctx.Set("Content-Type", "application/json")

	return ctx.Status(201).JSON(fiber.Map{
		"imageID": id,
	})
}

func (c *ImageController) likeImage(ctx *fiber.Ctx) error {
	c.semaphore.Acquire()
	defer c.semaphore.Release()

	req := new(likeRequest)
	err := ctx.BodyParser(req)
	if err != nil {
		return ctx.Status(400).SendString("Bad request")
	}

	err = c.imageService.AddLikes(req.ImageID, 1)
	if err != nil {
		return ctx.Status(404).SendString(err.Error())
	}

	return ctx.Status(200).SendString("Image liked")
}

func (c *ImageController) update(ctx *fiber.Ctx) error {
	c.semaphore.Acquire()
	defer c.semaphore.Release()

	req := new(updateRequest)
	err := ctx.BodyParser(req)
	if err != nil {
		return ctx.Status(400).SendString("Bad request")
	}

	token := req.Token

	image := models.Image{
		ImageID:     req.ImageID,
		Author:      req.Author,
		Title:       req.Title,
		Description: req.Description,
	}

	log.Printf("Image: %v", image)

	err = c.imageService.UpdateImage(image, token)
	if err != nil {
		return ctx.Status(404).SendString(err.Error())
	}

	return ctx.Status(200).SendString("Image updated")
}

func (c *ImageController) delete(ctx *fiber.Ctx) error {
	c.semaphore.Acquire()
	defer c.semaphore.Release()

	req := new(deleteRequest)
	err := ctx.BodyParser(req)
	if err != nil {
		return ctx.Status(400).SendString("Bad request")
	}

	token := req.Token

	err = c.imageService.DeleteImage(req.ImageID, token)
	if err != nil {
		return ctx.Status(404).SendString(err.Error())
	}

	return ctx.Status(200).SendString("Image deleted")
}

func (c *ImageController) deleteAll(ctx *fiber.Ctx) error {
	c.semaphore.Acquire()
	defer c.semaphore.Release()

	err := c.imageService.DeleteAllImages()
	if err != nil {
		return ctx.Status(404).SendString(err.Error())
	}

	return ctx.Status(200).SendString("All images deleted")
}

func (c *ImageController) revertTransaction(ctx *fiber.Ctx) error {
	c.semaphore.Acquire()
	defer c.semaphore.Release()

	id := ctx.Params("id")

	err := c.imageService.RevertSagaTransaction(id)
	if err != nil {
		return ctx.Status(404).SendString(err.Error())
	}

	log.Info().Msg(fmt.Sprintf("Transaction %s reverted", id))

	return ctx.Status(200).SendString("Transaction reverted")
}
