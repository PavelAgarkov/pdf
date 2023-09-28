package controller

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"pdf/internal/hash"
	"pdf/internal/logger"
	"time"
)

type MergeController struct {
	bc *BaseController
}

type MergeResponse struct {
	str string
	err error
}

func NewMergeController(bc *BaseController) *MergeController {
	return &MergeController{
		bc: bc,
	}
}

func (r *MergeResponse) GetStr() string {
	return r.str
}

func (r *MergeResponse) GetErr() error {
	return r.err
}

func (f *MergeController) Handle(
	ctx context.Context,
	loggerFactory *logger.Factory,
) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		defer RestoreController(loggerFactory, c, "empty controller")

		payload := struct {
			Key string `json:"key"`
		}{}

		files := make([]string, 0)
		form, err := c.MultipartForm()
		if err != nil {
			/* handle error */
		}
		for _, fileHeaders := range form.File {
			for _, fileHeader := range fileHeaders {
				files = append(files, fileHeader.Filename)
			}
		}
		loggerFactory.ErrorLog("errrprrrrrrr", zap.Stack("").String)
		loggerFactory.WarningLog("errrprrrrrrr")

		fv := form.Value["key"]
		_ = c.BodyParser(&payload)
		fmt.Println(payload, fv)
		_, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		hashCookie := ""
		if v := c.Cookies("X-HASH"); v == "" {
			hashCookie = string(hash.GenerateFirstLevelHash())
		}

		return c.JSON(fiber.Map{
			"one":  payload.Key,
			"hash": hashCookie,
		})
	}
}
