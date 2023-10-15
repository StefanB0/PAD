package controller

import "padauth/service"

type UserController struct {
	authService *service.AuthService
}

func NewUserController(authService *service.AuthService) *UserController {
	return &UserController{
		authService: authService,
	}
}

func (c *UserController) Run() {

}
