package controller

import (
	"padimage/models"
	"padimage/service"
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

func (c *ImageController) GetImage(imageID int64) (*models.Image, error) {
	return c.imageService.GetImage(imageID)
}

func (c *ImageController) CreateImage(image models.Image, token string) error {
	return c.imageService.CreateImage(image, token)
}

func (c *ImageController) UpdateImage(image models.Image, token string) error {
	return c.imageService.UpdateImage(image, token)
}

func (c *ImageController) DeleteImage(imageID int64, token string) error {
	return c.imageService.DeleteImage(imageID, token)
}

func (c *ImageController) GetImagesByAuthor(author string) ([]int64, error) {
	return c.imageService.GetImagesByAuthor(author)
}

func (c *ImageController) GetImagesByTag(tag string) ([]int64, error) {
	return c.imageService.GetImagesByTag(tag)
}

func (c *ImageController) DeleteAllImages() error {
	return c.imageService.DeleteAllImages()
}
