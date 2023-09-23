package pdf_operation

import (
	"context"
	"fmt"
	"os"
	"pdf/internal/adapter"
	"pdf/internal/entity"
	"pdf/internal/hash"
	"testing"
	"time"
)

func Test_remove_pages(t *testing.T) {
	p := adapter.NewPathAdapter()
	adapterLocator := adapter.NewAdapterLocator(
		adapter.NewFileAdapter(),
		p,
		adapter.NewPdfAdapter(),
		adapter.NewArchiveAdapter(p),
	)
	pathAdapter := adapterLocator.Locate(adapter.PathAlias).(*adapter.PathAdapter)

	conf := NewConfiguration(nil, []string{"2-11", "13-14"}, nil)

	firstLevelHash := hash.GenerateFirstLevelHash()
	secondLevelHash := hash.GenerateNextLevelHashByPrevious(firstLevelHash, true)
	expired := time.Now().Add(Timer5)

	inDir := pathAdapter.GenerateInDirPath(secondLevelHash)
	dirPath := pathAdapter.GenerateDirPathToFiles(secondLevelHash)
	outDir := pathAdapter.GenerateOutDirPath(secondLevelHash)
	archiveDir := pathAdapter.GenerateArchiveDirPath(secondLevelHash)

	ud := entity.NewUserData(firstLevelHash, secondLevelHash, expired)

	filesForReplace := []string{"./files/ServiceAgreement_template.pdf"}
	_, file, _ := pathAdapter.StepBack(adapter.Path(filesForReplace[0]))
	f, _ := os.ReadFile(filesForReplace[0])

	files := []string{string(inDir) + file}

	fileAdapter := adapterLocator.Locate(adapter.FileAlias).(*adapter.FileAdapter)
	err := fileAdapter.CreateDir(string(dirPath), 0777)
	err = fileAdapter.CreateDir(string(inDir), 0777)
	err = fileAdapter.CreateDir(string(outDir), 0777)
	err = fileAdapter.CreateDir(string(archiveDir), 0777)
	err = os.WriteFile(string(inDir)+file, f, 0777)

	operationFactory := NewOperationFactory()
	removePagesOperation := operationFactory.CreateNewOperation(conf, ud, files, dirPath, inDir, outDir, archiveDir, "", DestinationRemovePages)

	ctx := context.Background()
	_, err = removePagesOperation.Execute(ctx, adapterLocator, adapter.TarFormat)

	if err != nil {
		fmt.Println(err.Error())
	}
}
