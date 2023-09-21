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
	dirPathFile adapter.DirPath,
	inDir adapter.InDir,
	outDit adapter.OutDir,
	archiveDir adapter.ArchiveDir,
	splitPath adapter.SplitDir,
	destination Destination,
) Operation {
	bo := NewBaseOperation(configuration, ud, files, dirPathFile, inDir, outDit, archiveDir, destination)

	switch destination {
	case DestinationMerge:
		return NewMergeOperation(bo)
	case DestinationSplit:
		return NewSplitOperation(bo, splitPath)
	case DestinationRemovePages:
		return NewRemovePagesOperation(bo)
	default:
		return nil
	}
}
