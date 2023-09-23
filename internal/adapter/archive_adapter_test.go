package adapter

import (
	"context"
	"fmt"
	"github.com/mholt/archiver/v4"
	"os"
	"testing"
)

func Test_archive_rar_adapter_test(t *testing.T) {
	files, err := archiver.FilesFromDisk(nil, map[string]string{
		"./files/ServiceAgreement_template.pdf":  "ServiceAgreement_template.pdf",
		"./files/ServiceAgreement_template1.pdf": "ServiceAgreement_template1.pdf",
	})
	if err != nil {
		fmt.Println(err.Error())
	}

	out, err := os.Create("./files/archive.zip.zst")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {

		}
	}(out)

	zipFormat := archiver.CompressedArchive{
		Compression: archiver.Zstd{},
		Archival:    archiver.Zip{},
	}

	err = zipFormat.Archive(context.Background(), out, files)
	if err != nil {
		fmt.Println(err.Error())
	}
}
