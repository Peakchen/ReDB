package problempdfs

// ConvertToNewFormat 将旧的 API 中错题类型转换为新的 API 中的错题类型
// TODO: 修改获取错题方法，不转换
func ConvertToNewFormat(oldProblems []DetailedProblemsType) []ProblemForCreatingFiles {
	newProblems := []ProblemForCreatingFiles{}
	for _, pt := range oldProblems {
		if pt.Type != "" {
			var newPro ProblemForCreatingFiles
			var subIdxs []int
			var lastProblemID string
			for _, p := range pt.Problems {
				if lastProblemID != p.ProblemID && lastProblemID != "" {
					// 进入新题目
					// 判断 lastProblemID 避免刚进入循环时插入一个空题目
					newPro.SubIdxs = subIdxs
					newProblems = append(newProblems, newPro)
					subIdxs = []int{}

				}
				newPro = ProblemForCreatingFiles{
					BasicProblemForCreatingFiles: BasicProblemForCreatingFiles{
						Type:      p.Type,
						Book:      p.Book,
						Page:      p.Page,
						Column:    p.Column,
						Idx:       p.Idx,
						ProblemID: p.ProblemID,
						Full:      p.Full,
						How:       p.how,
						Reason:    p.Reason,
					},
				}
				subIdxs = append(subIdxs, p.SubIdx)
			}
			if len(pt.Problems) != 0 {
				// 最后一个题目因为还没触发 ”进入新题目“，手动加入
				newPro.SubIdxs = subIdxs
				newProblems = append(newProblems, newPro)
			}
		} else {
			if len(pt.Problems) != 0 {
				newPro := ProblemForCreatingFiles{
					BasicProblemForCreatingFiles: BasicProblemForCreatingFiles{
						Type:      pt.Problems[0].Type,
						Book:      pt.Problems[0].Book,
						Page:      pt.Problems[0].Page,
						Column:    pt.Problems[0].Column,
						Idx:       pt.Problems[0].Idx,
						ProblemID: pt.Problems[0].ProblemID,
						Full:      pt.Problems[0].Full,
						How:       pt.Problems[0].how,
						Reason:    pt.Problems[0].Reason,
					},
				}
				subIdxs := make([]int, len(pt.Problems))
				for i, p := range pt.Problems {
					subIdxs[i] = p.SubIdx
				}
				newPro.SubIdxs = subIdxs
				newProblems = append(newProblems, newPro)
			}
		}
	}
	return newProblems
}
