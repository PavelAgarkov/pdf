package controller

import (
	"context"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"pdf/internal"
	"pdf/internal/hash"
	"pdf/internal/logger"
	"pdf/internal/pdf_operation"
	"pdf/internal/service"
	"pdf/internal/storage"
)

type BaseController struct{}

type ResponseInterface interface {
	GetStr() string
	GetErr() error
}

func NewBaseController() *BaseController {
	return &BaseController{}
}

func (bc *BaseController) Select(
	ctx context.Context,
	ch chan ResponseInterface,
) ResponseInterface {
	select {
	case <-ctx.Done():
		res, ok := <-ch
		if !ok {
			return nil
		}
		return res
	case res := <-ch:
		return res
	}
}

func RestoreController(loggerFactory *logger.Factory, destination string) {
	if r := recover(); r != nil {
		panicStr := fmt.Sprintf(destination+" : Recovered. Panic: %s\n", r)
		loggerFactory.PanicLog(panicStr, "")
	}
}

func (bc *BaseController) isAuthenticated(
	operationStorage *storage.OperationStorage,
	c *fiber.Ctx,
) (pdf_operation.OperationDataInterface, error) {
	authHeader := c.GetReqHeaders()[internal.AuthenticationHeader][0]
	authToken := service.ParseBearerHeader(authHeader)
	operationData, hit := operationStorage.Get(hash.GenerateNextLevelHashByPrevious(internal.Hash1lvl(authToken), true))
	if !hit {
		errMsg := fmt.Sprintf("can't find hit from storage")
		return nil, errors.New(errMsg)
	}

	ok, err := service.IsAuthenticated(operationData.GetUserData().GetHash2Lvl(), internal.Hash1lvl(authToken))
	if err != nil {
		errMsg := fmt.Sprintf("can't access to storage")
		return nil, errors.New(errMsg)
	}
	if !ok {
		errMsg := fmt.Sprintf("can't access to files by hash")
		return nil, errors.New(errMsg)
	}

	return operationData, nil
}

func (bc *BaseController) isOverAuthenticated(
	operationStorage *storage.OperationStorage,
	c *fiber.Ctx,
) error {
	authHeader := c.GetReqHeaders()[internal.AuthenticationHeader][0]
	authToken := service.ParseBearerHeader(authHeader)
	_, hit := operationStorage.Get(hash.GenerateNextLevelHashByPrevious(internal.Hash1lvl(authToken), true))
	if hit {
		errMsg := fmt.Sprintf("can't process already in storage")
		return errors.New(errMsg)
	}

	return nil
}
