package controller

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"pdf/internal/hash"
	"pdf/internal/logger"
	"time"
)

type EmptyController struct {
	bc *BaseController
}

type EmptyResponse struct {
	str string
	err error
}

func NewEmptyController(bc *BaseController) *EmptyController {
	return &EmptyController{
		bc: bc,
	}
}

func (r *EmptyResponse) GetStr() string {
	return r.str
}

func (r *EmptyResponse) GetErr() error {
	return r.err
}

func (f *EmptyController) Handle(
	ctx context.Context,
	loggerFactory logger.Logger,
) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
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

		fv := form.Value["key"]
		_ = c.BodyParser(&payload)
		fmt.Println(payload, fv)
		defer RestoreController(loggerFactory, c)
		_, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		return c.JSON(fiber.Map{
			"one":  "goooo",
			"hash": hash.GenerateFirstLevelHash(),
		})
	}
}
