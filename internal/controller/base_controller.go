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

func (bc *BaseController) SelectResponse(
	ctx context.Context,
	ch chan ResponseInterface,
	start chan struct{},
) ResponseInterface {
	start <- struct{}{}
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
	loggerFactory *logger.Factory,
) (pdf_operation.OperationDataInterface, error) {
	authToken := service.ParseBearerHeader(c.GetReqHeaders()[internal.AuthenticationHeader])
	operationData, hit := operationStorage.Get(hash.GenerateNextLevelHashByPrevious(internal.Hash1lvl(authToken), true))
	if !hit {
		errMsg := fmt.Sprintf("can't find hit %s from storage", authToken)
		loggerFactory.ErrorLog(errMsg, "")
		return nil, errors.New(errMsg)
	}

	ok, err := service.IsAuthenticated(operationData.GetUserData().GetHash2Lvl(), internal.Hash1lvl(authToken))
	if err != nil {
		errMsg := fmt.Sprintf("can't acces %s to storage", authToken)
		loggerFactory.ErrorLog(fmt.Sprintf(errMsg+" %s", err.Error()), "")
		return nil, errors.New(errMsg)
	}
	if !ok {
		errMsg := fmt.Sprintf("can't acces to %s files by hash", authToken)
		loggerFactory.ErrorLog(errMsg, "")
		return nil, errors.New(errMsg)
	}

	return operationData, nil
}

func (bc *BaseController) isOverAuthenticated(
	operationStorage *storage.OperationStorage,
	c *fiber.Ctx,
	loggerFactory *logger.Factory,
) error {
	authToken := service.ParseBearerHeader(c.GetReqHeaders()[internal.AuthenticationHeader])
	_, hit := operationStorage.Get(hash.GenerateNextLevelHashByPrevious(internal.Hash1lvl(authToken), true))
	if hit {
		errMsg := fmt.Sprintf("can't process %s already in storage", authToken)
		loggerFactory.ErrorLog(errMsg, "")
		return errors.New(errMsg)
	}

	return nil
}
