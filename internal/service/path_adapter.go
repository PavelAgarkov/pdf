package service

import (
	"fmt"
	"pdf/internal/storage"
)

const frontendDist = "./pdf-frontend/dist"

type PathSplitFiles string
type DirPathFile string
type OutDir string

// хранить разрезанные файлы в ./files/Hash2lvl/split/ - так же и генерировать урл на скачивание через Hash2lvl

func GenerateDirPathToSplitFiles(hash2lvl storage.Hash2lvl) PathSplitFiles {
	return PathSplitFiles(fmt.Sprintf("./files/%s/split/", string(hash2lvl)))
}

func GenerateDirPathToFiles(hash2lvl storage.Hash2lvl) DirPathFile {
	return DirPathFile(fmt.Sprintf("./files/%s/", string(hash2lvl)))
}

func GenerateOutDirPath(hash2lvl storage.Hash2lvl) OutDir {
	return OutDir(fmt.Sprintf("./files/%s/out/", string(hash2lvl)))
}

func GenerateFrontendDist() string {
	return frontendDist
}
