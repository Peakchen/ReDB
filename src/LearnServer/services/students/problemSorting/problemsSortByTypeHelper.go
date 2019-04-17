package problemSorting

import "strings"

type toClassifyByType []problemWithTypeTimeAndCorrectDB

func (list toClassifyByType) InSameClass(i, j int) bool {
	return list[i].TypeInfo.Type == list[j].TypeInfo.Type
}

func (list toClassifyByType) Length() int {
	return len(list)
}

type toSortByTypeForType [][]problemWithTypeTimeAndCorrectDB

func (c toSortByTypeForType) Len() int {
	return len(c)
}

func (c toSortByTypeForType) Swap(i int, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c toSortByTypeForType) Less(i int, j int) bool {
	if c[i][0].TypeInfo.Chapter != c[j][0].TypeInfo.Chapter {
		return c[i][0].TypeInfo.Chapter < c[j][0].TypeInfo.Chapter
	}

	if c[i][0].TypeInfo.Section != c[j][0].TypeInfo.Section {
		return c[i][0].TypeInfo.Section < c[j][0].TypeInfo.Section
	}

	return c[i][0].TypeInfo.Priority < c[j][0].TypeInfo.Priority
}

type toSortByTimeForType []problemWithTypeTimeAndCorrectDB

func (c toSortByTimeForType) Len() int {
	return len(c)
}

func (c toSortByTimeForType) Swap(i int, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c toSortByTimeForType) Less(i int, j int) bool {
	if !c[i].Time.Equal(c[j].Time) {
		return c[i].Time.Before(c[j].Time)
	}

	if strings.Compare(c[i].ProblemID, c[j].ProblemID) != 0 {
		return strings.Compare(c[i].ProblemID, c[j].ProblemID) < 0
	}

	return c[i].SubIdx < c[j].SubIdx
}

func removeDuplicateProblems(problems []problemWithTypeTimeAndCorrectDB) []problemWithTypeTimeAndCorrectDB {
	// 去除重复的题目，并将留下的题目的time更新到最新（只是时间更新，不保证和最新的完全一致）
	result := []problemWithTypeTimeAndCorrectDB{}
	for _, p := range problems {
		found := false
		for j, pj := range result {
			if p.ProblemID == pj.ProblemID && p.SubIdx == pj.SubIdx {
				// 如果已有，更新时间
				if pj.Time.Before(p.Time) {
					result[j].Time = p.Time
				}
				found = true
				break
			}
		}
		if !found {
			result = append(result, p)
		}
	}
	return result
}
