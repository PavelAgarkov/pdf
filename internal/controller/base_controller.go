package controller

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"pdf/internal/logger"
)

type BaseController struct{}

type ResponseInterface interface {
	GetStr() string
}

func NewBaseController() *BaseController {
	return &BaseController{}
}

func (bc *BaseController) SelectResult(
	ctx context.Context,
	ch chan ResponseInterface,
	start chan struct{},
) ResponseInterface {
	start <- struct{}{}
	select {
	case <-ctx.Done():
		return nil
	case res := <-ch:
		return res
	}
}

func RestoreController(loggerFactory logger.Logger, c *fiber.Ctx) {
	if r := recover(); r != nil {
		errStr := fmt.Sprintf("Recovered. Panic: %s\n", r)
		loggerFactory.GetLogger(logger.ErrorName).Error(errStr)
		err := c.RedirectToRoute("root", map[string]interface{}{})
		if err != nil {
			return
		}
	}
}
