package pdf_operation

import (
	"context"
	"pdf/internal/adapter"
	"pdf/internal/entity"
	"strconv"
	"sync"
	"time"
)

const (
	Timer5  = 5 * time.Minute
	Timer10 = 10 * time.Minute
	Timer15 = 15 * time.Minute
)

type Operation interface {
	GetBaseOperation() *BaseOperation
	Execute(ctx context.Context, locator *adapter.Locator, format string) (string, error)
	GetDestination() string
}

// назначение операции - разделение файла, мерж файлов, сжатие и что придумаем еще

type Destination string
type OperationStatus string
type StoppedReason string

// тут записана операция, которую делает пользователь.
//Это нужно если пользователь решил на пол пути делать новую операцию(например хотел соединить, а потом решил разъединить).
//Это нужно для отмены его старых данных и удобства работы с ними

type BaseOperation struct {
	configuration *OperationConfiguration // конфигурации для выполнения операций, например диапазоны разбиения файла
	ud            *entity.UserData
	files         sync.Map
	dirPathFile   adapter.DirPath // путь до директории файла
	inDir         adapter.InDir
	outDir        adapter.OutDir
	archiveDir    adapter.ArchiveDir
	destination   Destination
	status        OperationStatus //статус операции нужен для контоля отмены токена и очистки памяти
	stoppedReason StoppedReason
}

func NewBaseOperation(
	configuration *OperationConfiguration,
	ud *entity.UserData,
	files []string,
	dirPathFile adapter.DirPath,
	idDir adapter.InDir,
	outDIr adapter.OutDir,
	archiveDir adapter.ArchiveDir,
	destination Destination,
) *BaseOperation {
	bo := &BaseOperation{
		configuration: configuration,
		ud:            ud,
		dirPathFile:   dirPathFile,
		inDir:         idDir,
		outDir:        outDIr,
		archiveDir:    archiveDir,
		destination:   destination,
		status:        OperationStatus(StatusStarted),
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

func (bo *BaseOperation) GetStatus() OperationStatus {
	return bo.status
}

func (bo *BaseOperation) GetStoppedReason() StoppedReason {
	return bo.stoppedReason
}

func (bo *BaseOperation) GetInDir() adapter.InDir {
	return bo.inDir
}

func (bo *BaseOperation) GetOutDir() adapter.OutDir {
	return bo.outDir
}

func (bo *BaseOperation) GetArchiveDir() adapter.ArchiveDir {
	return bo.archiveDir
}

func (bo *BaseOperation) SetStatus(status OperationStatus) *BaseOperation {
	bo.status = status
	return bo
}

func (bo *BaseOperation) SetStoppedReason(reason StoppedReason) *BaseOperation {
	bo.stoppedReason = reason
	return bo
}
