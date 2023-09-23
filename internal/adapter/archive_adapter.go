package adapter

import (
	"context"
	"fmt"
	"github.com/mholt/archiver/v4"
	"os"
	"pdf/internal/hash"
)

const (
	ArchiveAlias = "archive"

	ZipFormat    = ".zip"
	ZipZstFormat = ".zip.zst"
	TarFormat    = ".tar"
	TarGzFormat  = ".tar.gz"
)

type ArchiveAdapter struct {
	pathAdapter *PathAdapter
}

type Compressor struct {
	format            string
	compressedArchive archiver.CompressedArchive
}

func NewArchiveAdapter(pathAdapter *PathAdapter) *ArchiveAdapter {
	return &ArchiveAdapter{
		pathAdapter: pathAdapter,
	}
}

func (aa *ArchiveAdapter) GetAlias() string {
	return ArchiveAlias
}

func (aa *ArchiveAdapter) CreateCompressor(format string) (*Compressor, error) {
	switch format {
	case TarFormat:
		return &Compressor{
			format: format,
			compressedArchive: archiver.CompressedArchive{
				Compression: nil,
				Archival:    archiver.Tar{},
			},
		}, nil
	case TarGzFormat:
		return &Compressor{
			format: format,
			compressedArchive: archiver.CompressedArchive{
				Compression: archiver.Gz{},
				Archival:    archiver.Tar{},
			},
		}, nil
	case ZipZstFormat:
		return &Compressor{
			format: format,
			compressedArchive: archiver.CompressedArchive{
				Compression: archiver.Zstd{},
				Archival:    archiver.Zip{},
			},
		}, nil
	case ZipFormat:
		return &Compressor{
			format: format,
			compressedArchive: archiver.CompressedArchive{
				Compression: nil,
				Archival:    archiver.Zip{},
			},
		}, nil
	default:
		return &Compressor{}, nil
	}
}

func (aa *ArchiveAdapter) Archive(
	ctx context.Context,
	compressor *Compressor,
	outDirFilesMap map[string]string,
	hash2lvl hash.Hash2lvl,
	archiveDir ArchiveDir,
) (string, error) {

	//files, err := archiver.FilesFromDisk(nil, map[string]string{
	//"/path/on/disk/file1.txt": "file1.txt",
	//"./files/ServiceAgreement_template.pdf":  "ServiceAgreement_template.pdf",
	//"./files/ServiceAgreement_template1.pdf": "ServiceAgreement_template1.pdf",
	//"/path/on/disk/file3.txt": "",              // put in root of archive as file3.txt
	//"/path/on/disk/file4.txt": "subfolder/",    // put in subfolder as file4.txt
	//"/path/on/disk/folder":    "Custom Folder", // contents added recursively
	//})

	files, err := archiver.FilesFromDisk(nil, outDirFilesMap)
	if err != nil {
		return "", fmt.Errorf("can't prepare files for archive: %w", err)
	}

	_, archiveName, err := aa.pathAdapter.StepForward(Path(archiveDir), string(hash2lvl)+compressor.format)
	if err != nil {
		return "", fmt.Errorf("can't step forward with path %s: %w", string(archiveDir), err)
	}

	out, err := os.Create(string(archiveName))
	if err != nil {
		return "", fmt.Errorf("can't create archive file %s: %w", archiveName, err)
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {

		}
	}(out)

	err = compressor.compressedArchive.Archive(ctx, out, files)
	if err != nil {
		return "", fmt.Errorf("can't archivate files %s: %w", archiveName, err)
	}

	return string(archiveName), nil
}
