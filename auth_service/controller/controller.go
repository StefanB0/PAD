package controller

import (
	"padauth/service"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	authService *service.AuthService
}

func NewUserController(authService *service.AuthService) *UserController {
	return &UserController{
		authService: authService,
	}
}

func (c *UserController) Run() {
	app := fiber.New()

	app.Post("/login", c.login)
	app.Post("/register", c.register)
	app.Post("/delete", c.delete)

	app.Listen(":8080")
}

func (c *UserController) login(ctx *fiber.Ctx) error {
	req := new(loginRequest)
	err := ctx.BodyParser(req)
	if err != nil {
		return ctx.Status(400).SendString("Bad request")
	}

	accessToken, refreshToken, err := c.authService.Login(req.Username, req.Password)
	if err != nil {
		return ctx.Status(404).SendString(err.Error())
	}

	return ctx.Status(200).JSON(fiber.Map{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func (c *UserController) register(ctx *fiber.Ctx) error {
	req := new(registerRequest)
	err := ctx.BodyParser(req)
	if err != nil {
		return ctx.Status(400).SendString("Bad request")
	}

	err = c.authService.Register(req.Username, req.Password)
	if err != nil {
		return ctx.Status(404).SendString(err.Error())
	}

	return ctx.Status(200).SendString("Success")
}

func (c *UserController) delete(ctx *fiber.Ctx) error {
	req := new(deleteRequest)
	err := ctx.BodyParser(req)
	if err != nil {
		return ctx.Status(400).SendString("Bad request")
	}

	err = c.authService.Delete(req.RefreshToken)
	if err != nil {
		return ctx.Status(404).SendString(err.Error())
	}

	return ctx.Status(200).SendString("Success")
}
