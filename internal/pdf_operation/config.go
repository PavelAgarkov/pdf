package pdf_operation

import (
	"strconv"
	"strings"
)

type OperationConfiguration struct {
	splitIntervals       []string
	removePagesIntervals []string
	mergeOrder           []string
}

func NewConfiguration(
	splitIntervals,
	removePagesIntervals,
	mergeOrder []string,
) *OperationConfiguration {
	return &OperationConfiguration{
		splitIntervals:       splitIntervals,
		removePagesIntervals: removePagesIntervals,
		mergeOrder:           mergeOrder,
	}
}

func (oc *OperationConfiguration) GetSplitIntervals() []string {
	return oc.splitIntervals
}

func (oc *OperationConfiguration) GetRemovePagesIntervals() []string {
	return oc.removePagesIntervals
}

func (oc *OperationConfiguration) GetMergeOrder() []string {
	return oc.mergeOrder
}

func (oc *OperationConfiguration) parseIntervals(intervals []string) ([]int, [][]int) {
	filling := make([]int, 0)
	intervalInt := make([][]int, len(intervals), len(intervals))

	for k, v := range intervals {
		interval := strings.Split(v, "-")
		rangeInt := make([]int, 2, 2)
		if len(interval) == 2 {
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
