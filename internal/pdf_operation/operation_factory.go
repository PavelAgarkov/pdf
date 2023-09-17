package pdf_operation

import (
	"pdf/internal/service"
	"pdf/internal/storage"
)

type OperationsFactory struct{}

func NewOperationFactory() *OperationsFactory {
	return &OperationsFactory{}
}

func (*OperationsFactory) GetNewOperation(
	ud *storage.UserData,
	files []string,
	dirPathFile service.DirPathFile,
	outDit service.OutDir,
	destination string,
) Operation {
	bo := NewBaseOperation(ud, files, dirPathFile, outDit, Destination(destination))

	switch destination {
	case DestinationMerge:
		return NewMergeOperation(bo)
	case DestinationSplit:
		return NewSplitOperation(bo)
	case DestinationCut:
		return NewCutOperation(bo)
	case DestinationRemovePages:
		return NewRemovePagesOperation(bo)
	default:
		return nil
	}
}
