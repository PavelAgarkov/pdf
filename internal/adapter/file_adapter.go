package adapter

import (
	"os"
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
// директории GenerateOutDirPath() GenerateDirPathToFiles() GenerateDirPathToFiles() и записать их в операцию,
// после чего положить операцию в хранилище

func (fa *FileAdapter) CreateDir(dirPath string, perm os.FileMode) error {
	err := os.Mkdir(dirPath, perm)
	if err != nil {
		return err
	}
	return nil
}

func (fa *FileAdapter) GetAllEntriesFromDir(path, format string) (map[string]string, error) {
	entries, err := os.ReadDir(path)
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
