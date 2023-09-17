package pdf_operation

import (
	"pdf/internal/service"
	"pdf/internal/storage"
	"slices"
	"strconv"
	"sync"
)

const (
	DestinationMerge       = "merge"
	DestinationSplit       = "split"
	DestinationCut         = "cut"
	DestinationRemovePages = "remove_pages"
)

type Operation interface {
	GetBaseOperation() *BaseOperation
	Execute(pdfAdapter *service.PdfAdapter) error
}

// назначение операции - разделение файла, мерж файлов, сжатие и что придумаем еще

type Destination string

// тут записана операция, которую делает пользователь.
//Это нужно если пользователь решил на пол пути делать новую операцию(например хотел соединить, а потом решил разъединить).
//Это нужно для отмены его старых данных и удобства работы с ними

type BaseOperation struct {
	ud          *storage.UserData
	files       sync.Map
	dirPathFile service.DirPathFile // путь до директории файла
	outDir      service.OutDir
	destination Destination
}

func NewBaseOperation(
	ud *storage.UserData,
	files []string,
	dirPathFile service.DirPathFile,
	outDIr service.OutDir,
	destination Destination,
) *BaseOperation {
	bo := &BaseOperation{
		ud:          ud,
		dirPathFile: dirPathFile,
		outDir:      outDIr,
		destination: destination,
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
