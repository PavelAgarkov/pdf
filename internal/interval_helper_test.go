package internal

import (
	"fmt"
	"slices"
	"testing"
)

func Test_parseIntervals(t *testing.T) {
	many, intervals := ParseIntervals([]string{"1-22", "55-77"})

	find, ok := slices.BinarySearch(many, 55)

	fmt.Println(find, ok, intervals)
}
