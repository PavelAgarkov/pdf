package controller

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"io/fs"
	"path/filepath"
)

type FileController struct{}

func GetFC() *FileController {
	return &FileController{}
}

func (f *FileController) GetCallback(filesPath string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Route()
		filename := filesPath + c.Params("filename")
		err := filepath.WalkDir(filename, walk)
		if err != nil {
			return c.RedirectToRoute("root", map[string]interface{}{})
		}
		return c.Download(filename)
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
