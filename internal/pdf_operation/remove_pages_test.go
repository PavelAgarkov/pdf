package pdf_operation

import (
	"fmt"
	"os"
	"pdf/internal"
	"pdf/internal/adapter"
	"pdf/internal/storage"
	"testing"
	"time"
)

func Test_remove_pages(t *testing.T) {
	adapterLocator := adapter.NewAdapterLocator(
		adapter.NewFileAdapter(),
		adapter.NewPathAdapter(),
		adapter.NewPdfAdapter(),
		adapter.NewRarAdapterAdapter(),
	)
	pathAdapter := adapterLocator.Locate(adapter.PathAlias).(*adapter.PathAdapter)

	conf := NewConfiguration(nil, []string{"2-11", "13-14"}, nil)

	firstLevelHash := storage.GenerateFirstLevelHash()
	secondLevelHash := storage.GenerateNextLevelHashByPrevious(firstLevelHash, true)
	expired := time.Now().Add(Timer)

	inDir := pathAdapter.GenerateInDirPath(secondLevelHash)
	dirPath := pathAdapter.GenerateDirPathToFiles(secondLevelHash)
	outDir := pathAdapter.GenerateOutDirPath(secondLevelHash)

	ud := internal.NewUserData(firstLevelHash, secondLevelHash, expired)

	files := []string{"./files/ServiceAgreement_template.pdf"}

	_, file, _ := pathAdapter.StepBack(adapter.Path(files[0]))

	operationFactory := NewOperationFactory()
	removePagesOperation := operationFactory.CreateNewOperation(conf, ud, files, dirPath, inDir, outDir, "", DestinationRemovePages)

	f, _ := os.ReadFile(files[0])

	fileAdapter := adapterLocator.Locate(adapter.FileAlias).(*adapter.FileAdapter)
	err := fileAdapter.CreateDir(string(dirPath), 0777)
	err = fileAdapter.CreateDir(string(inDir), 0777)
	err = fileAdapter.CreateDir(string(outDir), 0777)

	err = os.WriteFile(string(inDir)+file, f, 0777)

	err = removePagesOperation.Execute(adapterLocator)

	if err != nil {
		fmt.Println(err.Error())
	}
}
