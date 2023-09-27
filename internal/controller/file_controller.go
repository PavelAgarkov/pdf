package controller

import (
	"context"
	"errors"
	"fmt"
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
	err error
}

func (r *Response) GetStr() string {
	return r.str
}

func (r *Response) GetErr() error {
	return r.err
}

func NewFileController(bc *BaseController) *FileController {
	return &FileController{
		bc: bc,
	}
}

func (f *FileController) Handle(
	ctx context.Context,
	filesPath string,
	loggerFactory *logger.Factory,
) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		defer RestoreController(loggerFactory, c)
		ctxC, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		cr := make(chan ResponseInterface)
		start := make(chan struct{})

		//v := c.Cookies("X-HASH")
		//fmt.Println(v)
		filename := filesPath + c.Params("filename")
		go realHandler(start, cr, filename)

		res := f.bc.SelectResult(ctxC, cr, start)

		// context cancelled
		if res == nil {
			return c.RedirectToRoute("root", map[string]interface{}{})
		}

		if res.GetStr() == "redirect" {
			fmt.Println(res.GetErr().Error())
			name := filesPath + "ServiceAgreement_template.zip"
			//return c.RedirectToRoute("root", map[string]interface{}{})
			c.Accepts("application/pdf")
			c.Accepts("application/zip")
			c.Accepts("application/x-bzip")
			c.Accepts("application/x-tar")
			//c.Accepts("application/x-7z-compressed")
			return c.Download(name, "ServiceAgreement_template.zip")
		}
		c.Request()

		//if res.GetStr() == "download" {
		//	return c.Download(filesPath)
		//}

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
		ch <- &Response{str: "redirect", err: err}
		return
	}
	ch <- &Response{str: "download"}

	return
}
