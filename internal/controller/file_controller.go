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
}
type Response struct {
	str string
}

func (r *Response) GetStr() string {
	return r.str
}

func GetFC() *FileController {
	return &FileController{}
}

func (f *FileController) FileController(filesPath string, loggerFactory *logger.LoggerFactory) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.Context(), 3*time.Second)
		defer cancel()

		ch := make(chan ResponseInterface)
		start := make(chan struct{})
		frontendLogger := loggerFactory.GetFrontendLogger()
		frontendLogger.Error("frontend")
		loggerFactory.GetErrorLogger().Error("errror")
		loggerFactory.GetInfoLogger().Info("Info")
		loggerFactory.GetWarningLogger().Warn("warning")

		filename := filesPath + c.Params("filename")
		go func() {
			<-start
			err := filepath.WalkDir(filename, walk)
			if err != nil {
				ch <- &Response{str: "redirect"}
				return
			}
			ch <- &Response{str: "download"}
			return
		}()

		res := getBaseController().SelectResult(ctx, ch, start)

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

func walk(s string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	if d.IsDir() {
		return errors.New("file must be f")
	}

	return nil
}
