package pdf_operation

import (
	"errors"
	"pdf/internal/service"
)

type CutOperation struct {
	baseOperation *BaseOperation
}

func NewCutOperation(
	bo *BaseOperation,
) *CutOperation {
	return &CutOperation{baseOperation: bo}
}

func (so *CutOperation) GetBaseOperation() *BaseOperation {
	return so.baseOperation
}

func (so *CutOperation) Execute(pdfAdapter *service.PdfAdapter) error {
	return errors.New("cut_operation")
}
