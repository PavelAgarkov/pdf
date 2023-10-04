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

func Test_remove_pages(t *testing.T) {
	p := adapter.NewPathAdapter()
	adapterLocator := locator.NewAdapterLocator(
		adapter.NewFileAdapter(),
		p,
		adapter.NewPdfAdapter(),
		adapter.NewArchiveAdapter(p),
	)
	pathAdapter := adapterLocator.Locate(adapter.PathAlias).(*adapter.PathAdapter)

	conf := NewConfiguration(nil, []string{"1", "2-11", "13-14", "1-14"})

	firstLevelHash := hash.GenerateFirstLevelHash()
	secondLevelHash := hash.GenerateNextLevelHashByPrevious(firstLevelHash, true)
	expired := time.Now().Add(internal.Timer5 * internal.Minute)

	inDir := pathAdapter.GenerateInDirPath(secondLevelHash)
	rootDir := pathAdapter.GenerateRootDir(secondLevelHash)
	outDir := pathAdapter.GenerateOutDirPath(secondLevelHash)
	archiveDir := pathAdapter.GenerateArchiveDirPath(secondLevelHash)

	ud := entity.NewUserData(firstLevelHash, secondLevelHash, expired)

	filesForReplace := []string{filepath.FromSlash("./files/ServiceAgreement_template.pdf")}
	_, file, _ := pathAdapter.StepBack(internal.Path(filesForReplace[0]))
	f, _ := os.ReadFile(filesForReplace[0])

	files := []string{string(inDir) + file}

	fileAdapter := adapterLocator.Locate(adapter.FileAlias).(*adapter.FileAdapter)
	err := fileAdapter.CreateDir(string(rootDir), 0777)
	err = fileAdapter.CreateDir(string(inDir), 0777)
	err = fileAdapter.CreateDir(string(outDir), 0777)
	err = fileAdapter.CreateDir(string(archiveDir), 0777)
	err = os.WriteFile(string(inDir)+file, f, 0777)

	operationFactory := NewOperationFactory()
	removePagesOperation := operationFactory.CreateNewOperation(conf, ud, files, rootDir, inDir, outDir, archiveDir, "", DestinationRemovePages)

	ctx := context.Background()
	_, err = removePagesOperation.Execute(ctx, adapterLocator, internal.TarFormat)

	if err != nil {
		fmt.Println(err.Error())
	}
}
