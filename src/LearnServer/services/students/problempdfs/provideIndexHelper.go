package problempdfs

func addIndex(typeProblemsRaw []detailedProblemOfTypeInfo) ([]detailedProblemsType, int) {
	// 添加index

	typeProblemsList := []detailedProblemsType{}
	index := 1
	for _, pt := range typeProblemsRaw {
		for i, p := range pt.Problems {
			found := false
			for j, pj := range pt.Problems {
				if j >= i {
					break
				}

				if p.ProblemID == pj.ProblemID {
					// 这里遍历的是已经处理过index的题目
					// 同属于一道大题，index应该一样
					pt.Problems[i].Index = pj.Index
					found = true
					// 保证一道大题只有一个index，故找到一个即可
					break
				}
			}

			if !found {
				pt.Problems[i].Index = index
				index++
			}
		}

		typeProblemsList = append(typeProblemsList, detailedProblemsType{
			Type:     pt.Type.Type,
			Problems: pt.Problems,
		})
	}

	return typeProblemsList, index - 1
}
