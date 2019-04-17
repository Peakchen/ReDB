package problemSorting

import (
	"sort"
	"time"
)

type toSortAllProblemsByTime []problemWithTypeTimeAndCorrectDB

func (list toSortAllProblemsByTime) Len() int {
	return len(list)
}

func (list toSortAllProblemsByTime) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list toSortAllProblemsByTime) Less(i, j int) bool {
	return list[i].Time.Before(list[j].Time)
}

func getCorrectStatus(allProblems toSortAllProblemsByTime, problemID string, subIdx int, timeConstraint time.Time) int {
	// 寻找在timeConstraint时间节点下，某个题目的correct状态
	// 1：最新做对了且以前没错过
	// 2：最新做对了，以前错过
	// 3：最新做错了
	// -1: 找不到这道题

	sort.Sort(allProblems)
	correct := -1
	for _, p := range allProblems {
		if timeConstraint.Before(p.Time) {
			break
		}
		if p.ProblemID == problemID && p.SubIdx == subIdx {
			if !p.Correct {
				// 3：最新做错了
				correct = 3
			} else {
				if correct == 3 {
					// 3说明之前的题目中已经做错了（题目是按时间顺序遍历的），所以改为2：最新做对了，以前错过
					correct = 2
				} else if correct == -1 {
					// 以前没出现过这道题，而现在最新是做对的
					correct = 1
				}
				// 如果原来为2或者1，不用变
			}
		}
	}
	return correct
}
