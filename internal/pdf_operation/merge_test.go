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

func Test_cut(t *testing.T) {
	p := adapter.NewPathAdapter()
	adapterLocator := adapter.NewAdapterLocator(
		adapter.NewFileAdapter(),
		p,
		adapter.NewPdfAdapter(),
		adapter.NewArchiveAdapter(p),
	)
	pathAdapter := adapterLocator.Locate(adapter.PathAlias).(*adapter.PathAdapter)

	conf := NewConfiguration(nil, nil, nil)

	firstLevelHash := hash.GenerateFirstLevelHash()
	secondLevelHash := hash.GenerateNextLevelHashByPrevious(firstLevelHash, true)
	expired := time.Now().Add(Timer5)

	inDir := pathAdapter.GenerateInDirPath(secondLevelHash)
	dirPath := pathAdapter.GenerateDirPathToFiles(secondLevelHash)
	outDir := pathAdapter.GenerateOutDirPath(secondLevelHash)
	archiveDir := pathAdapter.GenerateArchiveDirPath(secondLevelHash)

	filesForReplace := []string{"./files/ServiceAgreement_template.pdf", "./files/ServiceAgreement_template.pdf"}

	_, file0, _ := pathAdapter.StepBack(adapter.Path(filesForReplace[0]))
	_, file1, _ := pathAdapter.StepBack(adapter.Path(filesForReplace[1]))
	f0, _ := os.ReadFile(filesForReplace[0])
	f1, _ := os.ReadFile(filesForReplace[1])

	fileAdapter := adapterLocator.Locate(adapter.FileAlias).(*adapter.FileAdapter)
	err := fileAdapter.CreateDir(string(dirPath), 0777)
	err = fileAdapter.CreateDir(string(inDir), 0777)
	err = fileAdapter.CreateDir(string(outDir), 0777)
	err = fileAdapter.CreateDir(string(archiveDir), 0777)
	err = os.WriteFile(string(inDir)+file0, f0, 0777)
	err = os.WriteFile(string(inDir)+file1, f1, 0777)

	files := []string{string(inDir) + file0, string(inDir) + file1}
	ud := entity.NewUserData(firstLevelHash, secondLevelHash, expired)

	operationFactory := NewOperationFactory()
	mergePagesOperation := operationFactory.CreateNewOperation(conf, ud, files, dirPath, inDir, outDir, archiveDir, "", DestinationMerge)

	ctx := context.Background()
	_, err = mergePagesOperation.Execute(ctx, adapterLocator, adapter.ZipFormat)

	if err != nil {
		fmt.Println(err.Error())
	}
}
