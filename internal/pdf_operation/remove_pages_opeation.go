package pdf_operation

import (
	"errors"
	"pdf/internal/adapter"
)

const (
	DestinationRemovePages = "remove_pages"
)

type RemovePagesOperation struct {
	baseOperation *BaseOperation
}

func NewRemovePagesOperation(
	bo *BaseOperation,
) *RemovePagesOperation {
	return &RemovePagesOperation{baseOperation: bo}
}

func (rpo *RemovePagesOperation) GetDestination() string {
	return DestinationRemovePages
}

func (rpo *RemovePagesOperation) GetBaseOperation() *BaseOperation {
	return rpo.baseOperation
}

func (rpo *RemovePagesOperation) Execute(locator *adapter.Locator) error {
	return errors.New("remove_pages")
}
