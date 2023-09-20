package adapter

import (
	"fmt"
	"pdf/internal/storage"
	"testing"
)

func Test_GetAllEntriesInDir(t *testing.T) {
	pa := NewPathAdapter()
	secondLevelHash := "cc39e488d1c810a2640d348985fde3fe0bde24d4ca580a329421fff7773bd5a4"
	splitDir := pa.GenerateDirPathToSplitFiles(storage.Hash2lvl(secondLevelHash))
	fa := NewFileAdapter()

	e, err := fa.GetAllEntriesFromDir(string(splitDir), ".pdf")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(e[13])
}
