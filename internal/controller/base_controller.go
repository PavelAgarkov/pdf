package controller

import (
	"context"
)

type BaseController struct{}

type ResponseInterface interface {
	GetStr() string
}

func getBaseController() *BaseController {
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
