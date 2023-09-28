package adapter

import (
	"fmt"
	"pdf/internal"
	"testing"
)

func Test_GetAllEntriesInDir(t *testing.T) {
	pa := NewPathAdapter()
	secondLevelHash := "189f100c16a45d22e9b5145621521ad879ad16120e716a6148b39abdc7b71c35"
	splitDir := pa.GenerateDirPathToSplitFiles(internal.Hash2lvl(secondLevelHash))
	fa := NewFileAdapter()

	e, err := fa.GetAllEntriesFromDir(string(splitDir), ".pdf")

	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println(e[strconv.Itoa(13)])
	fmt.Println(e)
}
