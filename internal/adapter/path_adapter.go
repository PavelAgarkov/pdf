package adapter

import (
	"fmt"
	storage2 "pdf/internal/storage"
)

const (
	frontendDist = "./pdf-frontend/dist"
	PathAlias    = "path"
)

type PathSplitFiles string
type DirPathFile string
type OutDir string

type PathAdapter struct{}

func NewPathAdapter() *PathAdapter {
	return &PathAdapter{}
}

func (pa *PathAdapter) GetAlias() string {
	return PathAlias
}

// хранить разрезанные файлы в ./files/Hash2lvl/split/ - так же и генерировать урл на скачивание через Hash2lvl

func (pa *PathAdapter) GenerateDirPathToSplitFiles(hash2lvl storage2.Hash2lvl) PathSplitFiles {
	return PathSplitFiles(fmt.Sprintf("./files/%s/split/", string(hash2lvl)))
}

func (pa *PathAdapter) GenerateDirPathToFiles(hash2lvl storage2.Hash2lvl) DirPathFile {
	return DirPathFile(fmt.Sprintf("./files/%s/", string(hash2lvl)))
}

func (pa *PathAdapter) GenerateOutDirPath(hash2lvl storage2.Hash2lvl) OutDir {
	return OutDir(fmt.Sprintf("./files/%s/out/", string(hash2lvl)))
}

func (pa *PathAdapter) GenerateFrontendDist() string {
	return frontendDist
}
