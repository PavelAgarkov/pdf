package pdf_operation

import (
	"fmt"
	"slices"
	"testing"
)

func Test_parseIntervals(t *testing.T) {
	cnf := NewConfiguration(nil, nil)

	many, intervals := cnf.parseIntervals([]string{"1-22", "55-77"})

	find, ok := slices.BinarySearch(many, 55)

	fmt.Println(find, ok, intervals)
}
