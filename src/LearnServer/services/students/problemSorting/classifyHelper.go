package problemSorting

func constructProblemsByIndex(allProblems []problemWithTypeTimeAndCorrectDB, indexes [][]int) [][]problemWithTypeTimeAndCorrectDB {
	result := make([][]problemWithTypeTimeAndCorrectDB, len(indexes))
	for ri, is := range indexes {
		result[ri] = make([]problemWithTypeTimeAndCorrectDB, len(is))
		for rj, j := range is {
			result[ri][rj] = allProblems[j]
		}
	}
	return result
}
