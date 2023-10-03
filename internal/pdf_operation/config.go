package pdf_operation

import (
	"strconv"
	"strings"
)

type OperationConfiguration struct {
	splitIntervals       []string
	removePagesIntervals []string
}

func NewConfiguration(
	splitIntervals,
	removePagesIntervals []string,
) *OperationConfiguration {
	return &OperationConfiguration{
		splitIntervals:       splitIntervals,
		removePagesIntervals: removePagesIntervals,
	}
}

func (oc *OperationConfiguration) GetSplitIntervals() []string {
	return oc.splitIntervals
}

func (oc *OperationConfiguration) GetRemovePagesIntervals() []string {
	return oc.removePagesIntervals
}

func (oc *OperationConfiguration) parseIntervals(intervals []string) ([]int, [][]int) {
	filling := make([]int, 0)
	intervalInt := make([][]int, len(intervals), len(intervals))

	for k, v := range intervals {
		interval := strings.Split(v, "-")
		if len(interval) == 1 {
			rangeInt := make([]int, 2, 2)
			left, _ := strconv.ParseInt(interval[0], 10, 32)
			right, _ := strconv.ParseInt(interval[0], 10, 32)
			rangeInt[0] = int(left)
			rangeInt[1] = int(right)
			intervalInt[k] = rangeInt
			filling = append(filling, int(left))
		}
		if len(interval) == 2 {
			rangeInt := make([]int, 2, 2)
			left, _ := strconv.ParseInt(interval[0], 10, 32)
			right, _ := strconv.ParseInt(interval[1], 10, 32)
			rangeInt[0] = int(left)
			rangeInt[1] = int(right)
			intervalInt[k] = rangeInt
			for i := left; i <= right; i++ {
				filling = append(filling, int(i))
			}
		}
	}
	return filling, intervalInt
}
