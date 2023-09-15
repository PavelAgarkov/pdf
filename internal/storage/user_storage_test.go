package storage

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_user_storage_hash(t *testing.T) {
	us := NewUserStorage()

	firstLevelHash := us.GenerateFirstLevelHash()
	secondLevelHash := us.GenerateNextLevelHashByPrevious(firstLevelHash, true)

	if firstLevelHash != secondLevelHash {
		fmt.Println(firstLevelHash, secondLevelHash)
	}

	assert.NotEqual(t, firstLevelHash, secondLevelHash)

}
