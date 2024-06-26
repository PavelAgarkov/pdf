package adapter

import (
	"errors"
	"fmt"
	"path/filepath"
	"pdf/internal"
	"slices"
	"strings"
)

const (
	PathAlias = "path"
)

type PathAdapter struct{}

func GenerateFrontendDist() string {
	return filepath.FromSlash(internal.FrontendDist)
}

func NewPathAdapter() *PathAdapter {
	return &PathAdapter{}
}

func (pa *PathAdapter) GetAlias() string {
	return PathAlias
}

func (pa *PathAdapter) GenerateDirPathToSplitFiles(hash2lvl internal.Hash2lvl) internal.SplitDir {
	return internal.SplitDir(filepath.FromSlash(fmt.Sprintf("./files/%s/split/", string(hash2lvl))))
}

func (pa *PathAdapter) GenerateRootDir(hash2lvl internal.Hash2lvl) internal.RootDir {
	return internal.RootDir(filepath.FromSlash(fmt.Sprintf("./files/%s/", string(hash2lvl))))
}

func (pa *PathAdapter) GenerateOutDirPath(hash2lvl internal.Hash2lvl) internal.OutDir {
	return internal.OutDir(filepath.FromSlash(fmt.Sprintf("./files/%s/out/", string(hash2lvl))))
}

func (pa *PathAdapter) GenerateOutDirFile(hash2lvl internal.Hash2lvl, file string) string {
	return filepath.FromSlash(fmt.Sprintf("./files/%s/out/%s", string(hash2lvl), file))
}

func (pa *PathAdapter) GenerateArchiveDirPath(hash2lvl internal.Hash2lvl) internal.ArchiveDir {
	return internal.ArchiveDir(filepath.FromSlash(fmt.Sprintf("./files/%s/archive/", string(hash2lvl))))
}

func (pa *PathAdapter) GenerateInDirPath(hash2lvl internal.Hash2lvl) internal.InDir {
	return internal.InDir(filepath.FromSlash(fmt.Sprintf("./files/%s/in/", string(hash2lvl))))
}

func (pa *PathAdapter) StepBack(path internal.Path) (internal.Path, string, error) {
	pathStr := string(path)
	chunk := strings.Split(pathStr, filepath.FromSlash("/"))
	if len(chunk) < 1 || (len(chunk) == 1 && (chunk[0] == "." || chunk[0] == "..")) || pathStr == "" {
		return path, "", errors.New("path is root")
	}
	last := chunk[len(chunk)-1]
	chunk = slices.Delete(chunk, len(chunk)-1, len(chunk))
	pathWithoutLastResource := filepath.FromSlash(strings.Join(chunk, filepath.FromSlash("/")))

	return internal.Path(pathWithoutLastResource), last, nil
}

func (pa *PathAdapter) StepForward(path internal.Path, next string) (internal.Path, internal.Path, error) {
	pathStr := string(path)
	trimStr := strings.TrimRight(pathStr, filepath.FromSlash("/"))
	chunk := strings.Split(trimStr, filepath.FromSlash("/"))
	last := chunk[len(chunk)-1]

	if strings.Contains(last, ".") {
		return path, "", errors.New("this path closed for grow")
	}
	newPath := filepath.FromSlash(trimStr + "/" + next)

	return path, internal.Path(newPath), nil
}

func (pa *PathAdapter) BuildOutPathFilesMap(aliasMap map[string]string, hash2lvl internal.Hash2lvl) map[string]string {
	resultMap := make(map[string]string)

	for _, alias := range aliasMap {
		resultMap[pa.GenerateOutDirFile(hash2lvl, alias)] = alias
	}

	return resultMap
}
