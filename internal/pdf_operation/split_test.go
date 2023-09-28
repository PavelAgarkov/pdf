package pdf_operation

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"pdf/internal"
	"pdf/internal/adapter"
	"pdf/internal/entity"
	"pdf/internal/hash"
	"pdf/internal/locator"
	"testing"
	"time"
)

func Test_split(t *testing.T) {
	p := adapter.NewPathAdapter()
	adapterLocator := locator.NewAdapterLocator(
		adapter.NewFileAdapter(),
		p,
		adapter.NewPdfAdapter(),
		adapter.NewArchiveAdapter(p),
	)
	pathAdapter := adapterLocator.Locate(adapter.PathAlias).(*adapter.PathAdapter)

	conf := NewConfiguration([]string{"1-3", "1-2"}, nil, nil)
	fileAdapter := adapterLocator.Locate(adapter.FileAlias).(*adapter.FileAdapter)

	firstLevelHash := hash.GenerateFirstLevelHash()
	secondLevelHash := hash.GenerateNextLevelHashByPrevious(firstLevelHash, true)
	expired := time.Now().Add(internal.Timer5)

	inDir := pathAdapter.GenerateInDirPath(secondLevelHash)
	rootDir := pathAdapter.GenerateRootDir(secondLevelHash)
	outDir := pathAdapter.GenerateOutDirPath(secondLevelHash)
	splitDir := pathAdapter.GenerateDirPathToSplitFiles(secondLevelHash)
	archiveDir := pathAdapter.GenerateArchiveDirPath(secondLevelHash)

	ud := entity.NewUserData(firstLevelHash, secondLevelHash, expired)

	filesForReplace := []string{filepath.FromSlash("./files/ServiceAgreement_template.pdf")}
	_, file0, _ := pathAdapter.StepBack(internal.Path(filesForReplace[0]))
	f, _ := os.ReadFile(filesForReplace[0])

	files := []string{string(inDir) + file0}

	err := fileAdapter.CreateDir(string(rootDir), 0777)
	err = fileAdapter.CreateDir(string(inDir), 0777)
	err = fileAdapter.CreateDir(string(outDir), 0777)
	err = fileAdapter.CreateDir(string(splitDir), 0777)
	err = fileAdapter.CreateDir(string(archiveDir), 0777)
	err = os.WriteFile(string(inDir)+file0, f, 0777)

	operationFactory := NewOperationFactory()
	mergePagesOperation := operationFactory.CreateNewOperation(conf, ud, files, rootDir, inDir, outDir, archiveDir, splitDir, DestinationSplit)

	ctx := context.Background()
	_, err = mergePagesOperation.Execute(ctx, adapterLocator, internal.ZipZstFormat)

	if err != nil {
		fmt.Println(err.Error())
	}
}
