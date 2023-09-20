package adapter

import (
	"fmt"
	"pdf/internal/storage"
	"strconv"
	"testing"
)

func Test_GetAllEntriesInDir(t *testing.T) {
	pa := NewPathAdapter()
	secondLevelHash := "b4116a5731610794d2e50216b7be02c16c1407c14bf60808ea8c1276f1f11491"
	splitDir := pa.GenerateDirPathToSplitFiles(storage.Hash2lvl(secondLevelHash))
	fa := NewFileAdapter()

	e, err := fa.GetAllEntriesFromDir(string(splitDir), ".pdf")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(e[strconv.Itoa(13)])
}
