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
	p := adapter.NewPathAdapter()
	adapterLocator := adapter.NewAdapterLocator(
		adapter.NewFileAdapter(),
		adapter.NewPathAdapter(),
		adapter.NewPdfAdapter(),
		adapter.NewArchiveAdapter(p),
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

	operationData := pdf_operation.NewOperationData(
		ud,
		archiveDir,
		mergePagesOperation.GetBaseOperation().GetStatus(),
		mergePagesOperation.GetBaseOperation().GetStoppedReason(),
	)

	uStorage.Insert(secondLevelHash, operationData)

	op, _ := uStorage.Get(secondLevelHash)
	fmt.Println(op.GetUserData().GetHash2Lvl())

	//uStorage.Delete(secondLevelHash)

	//op, ok := uStorage.Get(secondLevelHash)
	//fmt.Println(ok)

	op, ok := uStorage.Put(secondLevelHash, operationData)
	fmt.Println(op, ok)

	uStorage.Delete(secondLevelHash)
	op, ok = uStorage.Get(secondLevelHash)
	fmt.Println(op, ok)
}
