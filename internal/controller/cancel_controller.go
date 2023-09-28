package controller

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"pdf/internal/logger"
)

type CancelController struct {
	bc *BaseController
}

type CancelResponse struct {
	str string
	err error
}

func (r *CancelResponse) GetStr() string {
	return r.str
}

func (r *CancelResponse) GetErr() error {
	return r.err
}

func NewCancelController(bc *BaseController) *CancelController {
	return &CancelController{
		bc: bc,
	}
}

func (f *CancelController) Handle(
	ctx context.Context,
	filesPath string,
	loggerFactory *logger.Factory,
) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return errors.New("")
	}
}
