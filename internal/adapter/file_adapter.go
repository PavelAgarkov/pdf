package adapter

import (
	"os"
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
