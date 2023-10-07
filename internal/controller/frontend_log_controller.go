package controller

import (
	"github.com/gofiber/fiber/v2"
	"pdf/internal/logger"
)

type FrontendLogController struct {
	bc *BaseController
}

type FrontendLogResponse struct {
	str string
	err error
}

func (r *FrontendLogResponse) GetStr() string {
	return r.str
}

func (r *FrontendLogResponse) GetErr() error {
	return r.err
}

func NewFrontendLogController(bc *BaseController) *FrontendLogController {
	return &FrontendLogController{
		bc: bc,
	}
}

type frontendLog struct {
	ErrorMsg  string `json:"errorMsg"`
	NotifyMsg string `json:"notifyMsg"`
}

func (cc *FrontendLogController) Handle(
	loggerFactory *logger.Factory,
) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		defer RestoreController(loggerFactory, "frontend log controller")

		fLJson := new(frontendLog)
		if err := c.BodyParser(fLJson); err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "yes",
			})
		}

		if fLJson.ErrorMsg != "" {
			loggerFactory.FrontendLog("error: " + fLJson.ErrorMsg)
		}

		if fLJson.NotifyMsg != "" {
			loggerFactory.FrontendLog("notify: " + fLJson.NotifyMsg)
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"error": "no",
		})
	}
}
