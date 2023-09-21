package storage

import (
	"context"
	"fmt"
	"pdf/internal/adapter"
	"pdf/internal/entity"
	"pdf/internal/hash"
	"pdf/internal/pdf_operation"
	"testing"
	"time"
)

func Test_user_in_memory_storage_test(t *testing.T) {
	adapterLocator := adapter.NewAdapterLocator(
		adapter.NewFileAdapter(),
		adapter.NewPathAdapter(),
		adapter.NewPdfAdapter(),
		adapter.NewRarAdapterAdapter(),
	)
	pathAdapter := adapterLocator.Locate(adapter.PathAlias).(*adapter.PathAdapter)

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	defer cancel()
	uStorage := NewInMemoryOperationStorage()
	uStorage.Run(ctx, pdf_operation.Timer5)

	firstLevelHash := hash.GenerateFirstLevelHash()
	secondLevelHash := hash.GenerateNextLevelHashByPrevious(firstLevelHash, true)

	conf := pdf_operation.NewConfiguration(nil, nil, nil)
	expired := time.Now().Add(pdf_operation.Timer5)

	inDir := pathAdapter.GenerateInDirPath(secondLevelHash)
	dirPath := pathAdapter.GenerateDirPathToFiles(secondLevelHash)
	outDir := pathAdapter.GenerateOutDirPath(secondLevelHash)
	archiveDir := pathAdapter.GenerateArchiveDirPath(secondLevelHash)

	files := []string{""}
	ud := entity.NewUserData(firstLevelHash, secondLevelHash, expired)

	operationFactory := pdf_operation.NewOperationFactory()
	mergePagesOperation := operationFactory.CreateNewOperation(conf, ud, files, dirPath, inDir, outDir, archiveDir, "", pdf_operation.DestinationMerge)

	uStorage.Insert(secondLevelHash, mergePagesOperation)

	op, _ := uStorage.Get(secondLevelHash)
	fmt.Println(op.GetBaseOperation().GetUserData().GetHash2Lvl())

	uStorage.Delete(secondLevelHash)

	op, ok := uStorage.Get(secondLevelHash)
	fmt.Println(ok)
}
