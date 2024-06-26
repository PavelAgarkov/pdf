package adapter

import (
	"context"
	"os"
	"path/filepath"
	"strings"
)

type FileAdapter struct{}

const (
	FileAlias = "file"
)

func NewFileAdapter() *FileAdapter {
	return &FileAdapter{}
}

func (fa *FileAdapter) GetAlias() string {
	return FileAlias
}

// перед операцией необходимо создать все необходимые
// директории GenerateOutDirPath() GenerateRootDir() и записать их в операцию,
// после чего положить операцию в хранилище

func (fa *FileAdapter) CreateDir(dirPath string, perm os.FileMode) error {
	err := os.Mkdir(filepath.FromSlash(dirPath), perm)
	if err != nil {
		return err
	}
	return nil
}

func (fa *FileAdapter) GetAllEntriesFromDir(ctx context.Context, path, format string) (map[string]string, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(filepath.FromSlash(path))
	mapFiles := make(map[string]string)

	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() || entry.Name() == "." || entry.Name() == ".." {
			continue
		}

		if strings.Contains(entry.Name(), format) {
			chunks := strings.Split(entry.Name(), "_")
			suffix := chunks[len(chunks)-1]
			suffixChunk := strings.Split(suffix, ".")
			fileNumber := suffixChunk[len(suffixChunk)-2]
			mapFiles[fileNumber] = entry.Name()
		}
	}

	return mapFiles, nil
}
