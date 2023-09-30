package adapter

import (
	"context"
	"fmt"
	"github.com/mholt/archiver/v4"
	"os"
	"pdf/internal"
)

const (
	ArchiveAlias = "archive"
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
	case internal.TarFormat:
		return &Compressor{
			format: format,
			compressedArchive: archiver.CompressedArchive{
				Compression: nil,
				Archival:    archiver.Tar{},
			},
		}, nil
	case internal.TarGzFormat:
		return &Compressor{
			format: format,
			compressedArchive: archiver.CompressedArchive{
				Compression: archiver.Gz{},
				Archival:    archiver.Tar{},
			},
		}, nil
	case internal.ZipZstFormat:
		return &Compressor{
			format: format,
			compressedArchive: archiver.CompressedArchive{
				Compression: archiver.Zstd{},
				Archival:    archiver.Zip{},
			},
		}, nil
	case internal.ZipFormat:
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
	hash2lvl internal.Hash2lvl,
	archiveDir internal.ArchiveDir,
) (string, error) {
	files, err := archiver.FilesFromDisk(nil, outDirFilesMap)
	if err != nil {
		return "", fmt.Errorf("can't prepare files for archive: %w", err)
	}

	_, archiveName, err := aa.pathAdapter.StepForward(internal.Path(archiveDir), string(hash2lvl)+compressor.format)
	if err != nil {
		return "", fmt.Errorf("can't step forward with path : %w", err)
	}

	out, err := os.Create(string(archiveName))
	if err != nil {
		return "", fmt.Errorf("can't create archive file : %w", err)
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {

		}
	}(out)

	err = compressor.compressedArchive.Archive(ctx, out, files)
	if err != nil {
		return "", fmt.Errorf("can't archivate files : %w", err)
	}

	return string(archiveName), nil
}
