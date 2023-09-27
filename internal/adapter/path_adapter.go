package adapter

import (
	"errors"
	"fmt"
	"path/filepath"
	"pdf/internal/hash"
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

func GenerateFrontendDist() string {
	return filepath.FromSlash(frontendDist)
}

func NewPathAdapter() *PathAdapter {
	return &PathAdapter{}
}

func (pa *PathAdapter) GetAlias() string {
	return PathAlias
}

// хранить разрезанные файлы в ./files/Hash2lvl/split/ - так же и генерировать урл на скачивание через Hash2lvl

func (pa *PathAdapter) GenerateDirPathToSplitFiles(hash2lvl hash.Hash2lvl) SplitDir {
	return SplitDir(filepath.FromSlash(fmt.Sprintf("./files/%s/split/", string(hash2lvl))))
}

func (pa *PathAdapter) GenerateDirPathToFiles(hash2lvl hash.Hash2lvl) DirPath {
	return DirPath(filepath.FromSlash(fmt.Sprintf("./files/%s/", string(hash2lvl))))
}

func (pa *PathAdapter) GenerateOutDirPath(hash2lvl hash.Hash2lvl) OutDir {
	return OutDir(filepath.FromSlash(fmt.Sprintf("./files/%s/out/", string(hash2lvl))))
}

func (pa *PathAdapter) GenerateOutDirFile(hash2lvl hash.Hash2lvl, file string) string {
	return filepath.FromSlash(fmt.Sprintf("./files/%s/out/%s", string(hash2lvl), file))
}

func (pa *PathAdapter) GenerateArchiveDirPath(hash2lvl hash.Hash2lvl) ArchiveDir {
	return ArchiveDir(filepath.FromSlash(fmt.Sprintf("./files/%s/archive/", string(hash2lvl))))
}

func (pa *PathAdapter) GenerateInDirPath(hash2lvl hash.Hash2lvl) InDir {
	return InDir(filepath.FromSlash(fmt.Sprintf("./files/%s/in/", string(hash2lvl))))
}

func (pa *PathAdapter) StepBack(path Path) (Path, string, error) {
	pathStr := string(path)
	chunk := strings.Split(pathStr, filepath.FromSlash("/"))
	if len(chunk) < 1 || (len(chunk) == 1 && (chunk[0] == "." || chunk[0] == "..")) || pathStr == "" {
		return path, "", errors.New("path is root")
	}
	last := chunk[len(chunk)-1]
	chunk = slices.Delete(chunk, len(chunk)-1, len(chunk))
	pathWithoutLastResource := filepath.FromSlash(strings.Join(chunk, filepath.FromSlash("/")))

	return Path(pathWithoutLastResource), last, nil
}

func (pa *PathAdapter) StepForward(path Path, next string) (Path, Path, error) {
	pathStr := string(path)
	trimStr := strings.TrimRight(pathStr, filepath.FromSlash("/"))
	chunk := strings.Split(trimStr, filepath.FromSlash("/"))
	last := chunk[len(chunk)-1]

	if strings.Contains(last, ".") {
		return path, Path(""), errors.New("this path closed for grow")
	}
	newPath := filepath.FromSlash(trimStr + "/" + next)

	return path, Path(newPath), nil
}

func (pa *PathAdapter) BuildOutPathFilesMap(aliasMap map[string]string, hash2lvl hash.Hash2lvl) map[string]string {
	resultMap := make(map[string]string)

	for _, alias := range aliasMap {
		resultMap[pa.GenerateOutDirFile(hash2lvl, alias)] = alias
	}

	return resultMap
}
