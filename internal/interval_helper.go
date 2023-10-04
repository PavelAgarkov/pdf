package internal

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseIntervals(intervals []string) ([]int, [][]int) {
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

func ParseIntIntervalsToString(toParse [][]int) []string {
	result := make([]string, 0)

	for _, v := range toParse {
		left := v[0]
		right := v[1]
		interval := fmt.Sprintf("%s-%s", strconv.Itoa(left), strconv.Itoa(right))
		result = append(result, interval)
	}

	return result
}
