package pdf_operation

import (
	"errors"
	"pdf/internal/service"
)

type RemovePagesOperation struct {
	baseOperation *BaseOperation
}

func NewRemovePagesOperation(
	bo *BaseOperation,
) *RemovePagesOperation {
	return &RemovePagesOperation{baseOperation: bo}
}

func (rpo *RemovePagesOperation) GetBaseOperation() *BaseOperation {
	return rpo.baseOperation
}

func (rpo *RemovePagesOperation) Execute(pdfAdapter *service.PdfAdapter) error {
	return errors.New("remove_pages")
}
