package pdf_operation

import (
	"pdf/internal"
	"pdf/internal/entity"
)

type OperationsFactory struct{}

func NewOperationFactory() *OperationsFactory {
	return &OperationsFactory{}
}

func (*OperationsFactory) CreateNewOperation(
	configuration *OperationConfiguration,
	ud *entity.UserData,
	files []string,
	rootDir internal.RootDir,
	inDir internal.InDir,
	outDit internal.OutDir,
	archiveDir internal.ArchiveDir,
	splitPath internal.SplitDir,
	destination internal.Destination,
) Operation {
	bo := NewBaseOperation(configuration, ud, files, rootDir, inDir, outDit, archiveDir, destination)

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
