package pdf_operation

import (
	"errors"
	"pdf/internal/adapter"
)

const (
	DestinationDownload = "download"
)

type DownloadOperation struct {
	baseOperation *BaseOperation
}

func NewDownloadOperation(
	bo *BaseOperation,
) *DownloadOperation {
	return &DownloadOperation{baseOperation: bo}
}

func (do *DownloadOperation) GetDestination() string {
	return DestinationDownload
}

func (do *DownloadOperation) GetBaseOperation() *BaseOperation {
	return do.baseOperation
}

func (do *DownloadOperation) Execute(locator *adapter.Locator) error {
	return errors.New("download")
}
