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

func Test_split(t *testing.T) {
	adapterLocator := adapter.NewAdapterLocator(
		adapter.NewFileAdapter(),
		adapter.NewPathAdapter(),
		adapter.NewPdfAdapter(),
		adapter.NewRarAdapterAdapter(),
	)
	pathAdapter := adapterLocator.Locate(adapter.PathAlias).(*adapter.PathAdapter)

	conf := NewConfiguration([]string{"1-3", "6-12"}, nil, nil)

	firstLevelHash := storage.GenerateFirstLevelHash()
	secondLevelHash := storage.GenerateNextLevelHashByPrevious(firstLevelHash, true)
	expired := time.Now().Add(Timer)

	inDir := pathAdapter.GenerateInDirPath(secondLevelHash)
	dirPath := pathAdapter.GenerateDirPathToFiles(secondLevelHash)
	outDir := pathAdapter.GenerateOutDirPath(secondLevelHash)
	splitDit := pathAdapter.GenerateDirPathToSplitFiles(secondLevelHash)

	ud := internal.NewUserData(firstLevelHash, secondLevelHash, expired)

	files := []string{"./files/ServiceAgreement_template.pdf"}

	_, file0, _ := pathAdapter.StepBack(adapter.Path(files[0]))
	//_, file1, _ := pathAdapter.StepBack(adapter.Path(files[1]))

	operationFactory := NewOperationFactory()
	mergePagesOperation := operationFactory.CreateNewOperation(conf, ud, files, dirPath, inDir, outDir, splitDit, DestinationSplit)

	f, _ := os.ReadFile(files[0])

	fileAdapter := adapterLocator.Locate(adapter.FileAlias).(*adapter.FileAdapter)
	err := fileAdapter.CreateDir(string(dirPath), 0777)
	err = fileAdapter.CreateDir(string(inDir), 0777)
	err = fileAdapter.CreateDir(string(outDir), 0777)
	err = fileAdapter.CreateDir(string(splitDit), 0777)

	err = os.WriteFile(string(inDir)+file0, f, 0777)
	//err = os.WriteFile(string(inDir)+file1, f, 0777)

	err = mergePagesOperation.Execute(adapterLocator)

	if err != nil {
		fmt.Println(err.Error())
	}
}
