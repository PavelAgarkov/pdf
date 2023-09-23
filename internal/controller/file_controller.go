package controller

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"io/fs"
	"path/filepath"
	"pdf/internal/logger"
	"time"
)

type FileController struct {
	bc *BaseController
}

type Response struct {
	str string
}

func (r *Response) GetStr() string {
	return r.str
}

func NewFileController(bc *BaseController) *FileController {
	return &FileController{
		bc: bc,
	}
}

func (f *FileController) Handle(
	ctx context.Context,
	filesPath string,
	loggerFactory logger.Logger,
) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		defer RestoreController(loggerFactory, c)
		ctxC, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		cr := make(chan ResponseInterface)
		start := make(chan struct{})

		filename := filesPath + c.Params("filename")
		go realHandler(start, cr, filename)

		res := f.bc.SelectResult(ctxC, cr, start)

		// context cancelled
		if res == nil {
			return c.RedirectToRoute("root", map[string]interface{}{})
		}

		if res.GetStr() == "redirect" {
			return c.RedirectToRoute("root", map[string]interface{}{})
		}

		if res.GetStr() == "download" {
			return c.Download(filename)
		}

		return nil
	}
}

func realHandler(start chan struct{}, ch chan ResponseInterface, filename string) {
	<-start

	err := filepath.WalkDir(
		filename,
		func(s string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return errors.New("file must be f")
			}

			return nil
		},
	)

	if err != nil {
		ch <- &Response{str: "redirect"}
		return
	}
	ch <- &Response{str: "download"}

	return
}
