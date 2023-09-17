package pdf_operation

import (
	"errors"
	"pdf/internal/service"
)

type MergeOperation struct {
	baseOperation *BaseOperation
}

func NewMergeOperation(
	bo *BaseOperation,
) *MergeOperation {
	return &MergeOperation{baseOperation: bo}
}

func (mo *MergeOperation) GetBaseOperation() *BaseOperation {
	return mo.baseOperation
}

func (mo *MergeOperation) Execute(pdfAdapter *service.PdfAdapter) error {
	return errors.New("merge_operation")
}
