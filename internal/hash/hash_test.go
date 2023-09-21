package hash

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_user_storage_hash(t *testing.T) {
	firstLevelHash := GenerateFirstLevelHash()
	secondLevelHash := GenerateNextLevelHashByPrevious(firstLevelHash, true)

	if string(firstLevelHash) != string(secondLevelHash) {
		fmt.Println(firstLevelHash, secondLevelHash)
	}

	assert.NotEqual(t, firstLevelHash, secondLevelHash)
}
