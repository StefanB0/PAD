package controller

import (
	"log"
	"padimage/models"
	"padimage/service"

	"github.com/gofiber/fiber/v2"
)

type ImageController struct {
	imageService *service.ImageService
	tokenService *service.TokenService
}

func NewImageController(is *service.ImageService, ts *service.TokenService) *ImageController {
	return &ImageController{
		imageService: is,
		tokenService: ts,
	}
}

func (c *ImageController) Run() {
	app := fiber.New()

	app.Post("/getImage", c.getImage)
	app.Post("/getImageInfo", c.getImageInfo)
	app.Post("/upload", c.upload)
	app.Post("/update", c.update)
	app.Post("/delete", c.delete)

	app.Listen(":8081")
}

func (c *ImageController) getImage(ctx *fiber.Ctx) error {
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

	return ctx.Status(200).Send(image.ImageChunk)
}

func (c *ImageController) getImageInfo(ctx *fiber.Ctx) error {
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
	req := new(uploadRequest)

	form, err := ctx.MultipartForm()
	if err != nil {
		return ctx.Status(400).SendString("Bad request")
	}

	req.Token = form.Value["token"][0]
	req.Author = form.Value["author"][0]
	req.Title = form.Value["title"][0]
	req.Description = form.Value["description"][0]
	req.Tags = form.Value["tags"]
	log.Println(req.Tags[1])

	fileheader, err := ctx.FormFile("image")
	if err != nil {
		return ctx.Status(400).SendString("Bad request")
	}

	file, err := fileheader.Open()
	if err != nil {
		return ctx.Status(400).SendString("Bad request")
	}

	req.ImageBytes = make([]byte, fileheader.Size)
	_, err = file.Read(req.ImageBytes)
	if err != nil {
		return ctx.Status(400).SendString("Bad request")
	}

	

	image := models.Image{
		ImageID:     0,
		Author:      req.Author,
		Title:       req.Title,
		Description: req.Description,
		Tags:        req.Tags,
		ImageChunk:  req.ImageBytes,
	}

	err = c.imageService.CreateImage(image, req.Token)

	ctx.Set("Content-Type", "image/jpeg")
	return ctx.Status(200).Send(image.ImageChunk)
}

func (c *ImageController) update(ctx *fiber.Ctx) error {
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

	err = c.imageService.UpdateImage(image, token)
	if err != nil {
		return ctx.Status(404).SendString(err.Error())
	}

	return ctx.Status(200).SendString("Image updated")
}

func (c *ImageController) delete(ctx *fiber.Ctx) error {
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

	err := c.imageService.DeleteAllImages()
	if err != nil {
		return ctx.Status(404).SendString(err.Error())
	}

	return ctx.Status(200).SendString("All images deleted")
}
