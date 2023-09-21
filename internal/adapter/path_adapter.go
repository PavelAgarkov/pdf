package adapter

import (
	"errors"
	"fmt"
	storage2 "pdf/internal/storage"
	"slices"
	"strings"
)

const (
	frontendDist = "./pdf-frontend/dist"
	PathAlias    = "path"
)

type SplitDir string
type DirPath string
type OutDir string
type InDir string
type ArchiveDir string
type Path string

type PathAdapter struct{}

func NewPathAdapter() *PathAdapter {
	return &PathAdapter{}
}

func (pa *PathAdapter) GetAlias() string {
	return PathAlias
}

// хранить разрезанные файлы в ./files/Hash2lvl/split/ - так же и генерировать урл на скачивание через Hash2lvl

func (pa *PathAdapter) GenerateDirPathToSplitFiles(hash2lvl storage2.Hash2lvl) SplitDir {
	return SplitDir(fmt.Sprintf("./files/%s/split/", string(hash2lvl)))
}

func (pa *PathAdapter) GenerateDirPathToFiles(hash2lvl storage2.Hash2lvl) DirPath {
	return DirPath(fmt.Sprintf("./files/%s/", string(hash2lvl)))
}

func (pa *PathAdapter) GenerateOutDirPath(hash2lvl storage2.Hash2lvl) OutDir {
	return OutDir(fmt.Sprintf("./files/%s/out/", string(hash2lvl)))
}

func (pa *PathAdapter) GenerateArchiveDirPath(hash2lvl storage2.Hash2lvl) ArchiveDir {
	return ArchiveDir(fmt.Sprintf("./files/%s/archive/", string(hash2lvl)))
}

func (pa *PathAdapter) GenerateInDirPath(hash2lvl storage2.Hash2lvl) InDir {
	return InDir(fmt.Sprintf("./files/%s/in/", string(hash2lvl)))
}

func (pa *PathAdapter) StepBack(path Path) (Path, string, error) {
	pathStr := string(path)
	chunk := strings.Split(pathStr, "/")
	if len(chunk) < 1 || (len(chunk) == 1 && (chunk[0] == "." || chunk[0] == "..")) || pathStr == "" {
		return path, "", errors.New("path is root")
	}
	last := chunk[len(chunk)-1]
	chunk = slices.Delete(chunk, len(chunk)-1, len(chunk))
	pathWithoutLastResource := strings.Join(chunk, "/")

	return Path(pathWithoutLastResource), last, nil
}

func (pa *PathAdapter) StepForward(path Path, next string) (Path, Path, error) {
	pathStr := string(path)
	trimStr := strings.TrimRight(pathStr, "/")
	chunk := strings.Split(trimStr, "/")
	last := chunk[len(chunk)-1]

	if strings.Contains(last, ".") {
		return path, Path(""), errors.New("this path closed for grow")
	}
	newPath := trimStr + "/" + next

	return path, Path(newPath), nil
}

func (pa *PathAdapter) GenerateFrontendDist() string {
	return frontendDist
}
