package problemSorting

import "strings"

type toClassifyByTime []problemWithTypeTimeAndCorrectDB

func (list toClassifyByTime) InSameClass(i, j int) bool {
	return list[i].Time.Year() == list[j].Time.Year() && list[i].Time.YearDay() == list[j].Time.YearDay()
}

func (list toClassifyByTime) Length() int {
	return len(list)
}

type toSortByTimeForTime [][]problemWithTypeTimeAndCorrectDB

func (c toSortByTimeForTime) Len() int {
	return len(c)
}

func (c toSortByTimeForTime) Swap(i int, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c toSortByTimeForTime) Less(i int, j int) bool {
	return c[i][0].Time.Before(c[j][0].Time)
}

type toSortByTypeForTime []problemWithTypeTimeAndCorrectDB

func (c toSortByTypeForTime) Len() int {
	return len(c)
}

func (c toSortByTypeForTime) Swap(i int, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c toSortByTypeForTime) Less(i int, j int) bool {
	if c[i].TypeInfo.Chapter != c[j].TypeInfo.Chapter {
		return c[i].TypeInfo.Chapter < c[j].TypeInfo.Chapter
	}

	if c[i].TypeInfo.Section != c[j].TypeInfo.Section {
		return c[i].TypeInfo.Section < c[j].TypeInfo.Section
	}

	if c[i].TypeInfo.Priority != c[j].TypeInfo.Priority {
		return c[i].TypeInfo.Priority < c[j].TypeInfo.Priority
	}

	// 此时还没return的是同一题型的
	if strings.Compare(c[i].ProblemID, c[j].ProblemID) != 0 {
		return strings.Compare(c[i].ProblemID, c[j].ProblemID) < 0
	}

	return c[i].SubIdx < c[j].SubIdx
}
