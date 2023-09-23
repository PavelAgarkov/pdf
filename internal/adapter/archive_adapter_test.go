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
		//"/path/on/disk/file1.txt": "file1.txt",
		"./files/ServiceAgreement_template.pdf":  "ServiceAgreement_template.pdf",
		"./files/ServiceAgreement_template1.pdf": "ServiceAgreement_template1.pdf",
		//"/path/on/disk/file3.txt": "",              // put in root of archive as file3.txt
		//"/path/on/disk/file4.txt": "subfolder/",    // put in subfolder as file4.txt
		//"/path/on/disk/folder":    "Custom Folder", // contents added recursively
	})
	if err != nil {
		fmt.Println(err.Error())
	}

	// create the output file we'll write to
	//out, err := os.Create("./files/archive.tar.gz")

	// формат .zip.zst сжимает архив .zip более чем в 2 раза! но нужно 2 раза разархивировать
	out, err := os.Create("./files/archive.zip.zst")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {

		}
	}(out)

	// we can use the CompressedArchive type to gzip a tarball
	// (compression is not required; you could use Tar directly)

	zipFormat := archiver.CompressedArchive{
		//nil,
		//nil,
		Compression: archiver.Zstd{},
		Archival:    archiver.Zip{},
	}

	//tarGzFormat := archiver.CompressedArchive{
	//	//nil,
	//	//nil,
	//	Compression: archiver.Gz{},
	//	Archival:    archiver.Tar{},
	//}

	// create the archive
	err = zipFormat.Archive(context.Background(), out, files)
	if err != nil {
		fmt.Println(err.Error())
	}
}
