package controller

import "padrecommendations/service"

type Controller struct {
	recommendationService *service.RecommendationService
}

func NewController(service *service.RecommendationService) *Controller {
	return &Controller{
		recommendationService: service,
	}
}

func (c *Controller) Run() {}