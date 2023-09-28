package controller

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"pdf/internal"
	"pdf/internal/logger"
	"pdf/internal/service"
	"pdf/internal/storage"
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
	operationStorage *storage.OperationStorage,
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

		newHashToBearer := ""
		authToken := service.ParseBearerHeader(c.GetReqHeaders()[internal.AuthenticationHeader])
		operationData, hit := operationStorage.Get(internal.Hash2lvl(authToken))
		if !hit {
			errMsg := fmt.Sprintf("cancel controller: can't find hit %s from storage", authToken)
			loggerFactory.ErrorLog(errMsg, "")
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": errMsg,
			})
		}
		ok, err := service.IsAuthenticated(operationData.GetUserData().GetHash2Lvl(), internal.Hash1lvl(authToken))
		if err != nil {
			errMsg := fmt.Sprintf("cancel controller: can't delete %s from storage", authToken)
			loggerFactory.ErrorLog(fmt.Sprintf(errMsg+" %s", err.Error()), "")
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": errMsg,
			})
		}

		if !ok {
			newHashToBearer = service.GenerateBearerToken()
		}

		return c.JSON(fiber.Map{
			"one":  payload.Key,
			"hash": newHashToBearer,
		})
	}
}
