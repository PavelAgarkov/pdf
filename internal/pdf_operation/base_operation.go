package pdf_operation

import (
	"context"
	"pdf/internal"
	"pdf/internal/entity"
	"pdf/internal/locator"
	"strconv"
	"sync"
)

type Operation interface {
	GetBaseOperation() *BaseOperation
	Execute(ctx context.Context, locator *locator.Locator, format string) (string, error)
	GetDestination() string
}

type BaseOperation struct {
	configuration *OperationConfiguration // конфигурации для выполнения операций, например диапазоны разбиения файла
	ud            *entity.UserData
	files         sync.Map
	rootDir       internal.RootDir
	inDir         internal.InDir
	outDir        internal.OutDir
	archiveDir    internal.ArchiveDir
	destination   internal.Destination
	status        internal.OperationStatus
	stoppedReason internal.StoppedReason
}

func NewBaseOperation(
	configuration *OperationConfiguration,
	ud *entity.UserData,
	files []string,
	rootDir internal.RootDir,
	idDir internal.InDir,
	outDIr internal.OutDir,
	archiveDir internal.ArchiveDir,
	destination internal.Destination,
) *BaseOperation {
	bo := &BaseOperation{
		configuration: configuration,
		ud:            ud,
		rootDir:       rootDir,
		inDir:         idDir,
		outDir:        outDIr,
		archiveDir:    archiveDir,
		destination:   destination,
		status:        internal.OperationStatus(internal.StatusStarted),
	}

	for k, filename := range files {
		key := strconv.Itoa(k)
		bo.files.Store(key, filename)
	}

	return bo
}

func (bo *BaseOperation) GetAllPaths() []string {
	paths := make([]string, 0)

	bo.files.Range(func(key, value any) bool {
		paths = append(paths, value.(string))
		return true
	})

	return paths
}

func (bo *BaseOperation) GetConfiguration() *OperationConfiguration {
	return bo.configuration
}

func (bo *BaseOperation) GetUserData() *entity.UserData {
	return bo.ud
}

func (bo *BaseOperation) GetStatus() internal.OperationStatus {
	return bo.status
}

func (bo *BaseOperation) GetStoppedReason() internal.StoppedReason {
	return bo.stoppedReason
}

func (bo *BaseOperation) GetInDir() internal.InDir {
	return bo.inDir
}

func (bo *BaseOperation) GetOutDir() internal.OutDir {
	return bo.outDir
}

func (bo *BaseOperation) GetArchiveDir() internal.ArchiveDir {
	return bo.archiveDir
}

func (bo *BaseOperation) GetRootDir() internal.RootDir {
	return bo.rootDir
}

func (bo *BaseOperation) SetStatus(status internal.OperationStatus) *BaseOperation {
	bo.status = status
	return bo
}

func (bo *BaseOperation) SetStoppedReason(reason internal.StoppedReason) *BaseOperation {
	bo.stoppedReason = reason
	return bo
}
