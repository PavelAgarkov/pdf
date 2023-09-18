package pdf_operation

import (
	"errors"
	"pdf/internal/adapter"
)

const (
	DestinationCut = "cut"
)

type CutOperation struct {
	baseOperation *BaseOperation
}

func NewCutOperation(
	bo *BaseOperation,
) *CutOperation {
	return &CutOperation{baseOperation: bo}
}

func (so *CutOperation) GetDestination() string {
	return DestinationCut
}

func (so *CutOperation) GetBaseOperation() *BaseOperation {
	return so.baseOperation
}

func (so *CutOperation) Execute(locator *adapter.Locator) error {
	return errors.New("cut_operation")
}
