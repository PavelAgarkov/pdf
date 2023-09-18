package pdf_operation

import (
	"pdf/internal"
	"pdf/internal/adapter"
)

type OperationsFactory struct{}

func NewOperationFactory() *OperationsFactory {
	return &OperationsFactory{}
}

func (*OperationsFactory) CreateNewOperation(
	configuration *OperationConfiguration,
	ud *internal.UserData,
	files []string,
	dirPathFile adapter.DirPathFile,
	outDit adapter.OutDir,
	destination string,
) Operation {
	bo := NewBaseOperation(configuration, ud, files, dirPathFile, outDit, Destination(destination))

	switch destination {
	case DestinationMerge:
		return NewMergeOperation(bo)
	case DestinationSplit:
		return NewSplitOperation(bo)
	case DestinationCut:
		return NewCutOperation(bo)
	case DestinationRemovePages:
		return NewRemovePagesOperation(bo)
	case DestinationDownload:
		return NewDownloadOperation(bo)
	default:
		return nil
	}
}
