package pdf_operation

import (
	"pdf/internal/adapter"
	"pdf/internal/entity"
	"slices"
	"strconv"
	"sync"
	"time"
)

//		 		expired
//		   		/\
//	       		|
//
// started->processed->awaiting_download->completed
//
//	   			|
//	  			\/
//			canceled

const (
	StatusStarted          = "started"
	StatusProcessed        = "processed"
	StatusCompleted        = "completed"
	StatusExpired          = "expired"
	StatusCanceled         = "canceled"
	StatusAwaitingDownload = "awaiting_download"
)

const (
	Timer5  = 5 * time.Minute
	Timer10 = 10 * time.Minute
	Timer15 = 15 * time.Minute
)

type Operation interface {
	GetBaseOperation() *BaseOperation
	Execute(locator *adapter.Locator) error
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
		convertKeyInt, err := strconv.Atoi(key.(string))
		if err != nil {
			return false
		}
		paths = slices.Insert(paths, convertKeyInt, value.(string))
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

func (bo *BaseOperation) CanDeleted() bool {
	return bo.GetStatus() == StatusExpired ||
		bo.GetStatus() == StatusCanceled ||
		bo.GetStatus() == StatusCompleted
}

func (bo *BaseOperation) SetStatus(status OperationStatus) *BaseOperation {
	bo.status = status
	return bo
}

func (bo *BaseOperation) SetStoppedReason(reason StoppedReason) *BaseOperation {
	bo.stoppedReason = reason
	return bo
}
