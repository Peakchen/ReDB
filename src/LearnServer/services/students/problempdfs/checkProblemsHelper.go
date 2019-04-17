package problempdfs

import (
	"log"
	"sort"
	"strings"

	// "LearnServer/models/contentDB"
	// "LearnServer/models/userDB"

	"LearnServer/models/contentDB"
	"LearnServer/models/userDB"
)

func getHowCodeInCheckProblems(how string) int {
	if how == "选择题" {
		return 0
	}
	if how == "填空题" {
		return 1
	}
	return 2
}

func judgeFullHelper(subIdx int) bool {
	if subIdx == -1 {
		return true
	}
	return false
}

type checkProblemSlice []detailedProblem

func (c checkProblemSlice) Len() int {
	return len(c)
}

func (c checkProblemSlice) Swap(i int, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c checkProblemSlice) Less(i int, j int) bool {
	if c[i].howCode != c[j].howCode {
		return c[i].howCode < c[j].howCode
	}

	if strings.Compare(c[i].ProblemID, c[j].ProblemID) > 0 {
		return false
	}
	if strings.Compare(c[i].ProblemID, c[j].ProblemID) < 0 {
		return true
	}

	return c[i].SubIdx < c[j].SubIdx
}

type checkProblemTypeSlice []detailedProblemOfTypeInfo

func (c checkProblemTypeSlice) Len() int {
	return len(c)
}

func (c checkProblemTypeSlice) Swap(i int, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c checkProblemTypeSlice) Less(i int, j int) bool {
	if c[i].Type.Chapter != c[j].Type.Chapter {
		return c[i].Type.Chapter < c[j].Type.Chapter
	}

	if c[i].Type.Section != c[j].Type.Section {
		return c[i].Type.Section < c[j].Type.Section
	}

	return c[i].Type.Priority < c[j].Type.Priority
}

type wrongProblemForCheckSlice []detailedProblem

func (c wrongProblemForCheckSlice) Len() int {
	return len(c)
}

func (c wrongProblemForCheckSlice) Swap(i int, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c wrongProblemForCheckSlice) Less(i int, j int) bool {
	ciHow := getHowCodeInCheckProblems(c[i].how)
	cjHow := getHowCodeInCheckProblems(c[j].how)
	if ciHow != cjHow {
		// howCode大的在前面，即顺序“其它→填空→选择”
		return ciHow > cjHow
	}

	if strings.Compare(c[i].ProblemID, c[j].ProblemID) > 0 {
		return false
	}
	if strings.Compare(c[i].ProblemID, c[j].ProblemID) < 0 {
		return true
	}

	return c[i].SubIdx < c[j].SubIdx
}

func getCheckProblemsForOneType(id string, typeName string, hows [3]bool, wrongProblems wrongProblemForCheckSlice) ([]detailedProblem, error) {
	// 为一个题型获取检验题，hows: 这个类型的howCode状态，wrongProblems: 这个类型做过的错题
	problemsForOneType := []detailedProblem{}
	for code, hasProblems := range hows {
		if hasProblems {
			problems, err := getProblemsOfHowAndType(code, typeName)
			if err != nil {
				return nil, err
			}
			// 找一道没做过的题目
			if p, err := getOneProblemNotDone(id, problems); err == nil {
				problemsForOneType = append(problemsForOneType, detailedProblem{
					ProblemID: p.ProblemID,
					SubIdx:    p.SubIdx,
					Full:      judgeFullHelper(p.SubIdx),
					howCode:   code,
				})
			}
		}
	}

	// 没有找到未做过的且出题方式匹配的检验题
	if len(problemsForOneType) == 0 {
		for code, hasProblems := range hows {
			if hasProblems {
				var ps1 []problem
				var ps2 []problem
				switch code {
				case 0:
					{
						// 选择题
						var err error
						ps1, err = getProblemsOfHowAndType(1, typeName)
						if err != nil {
							return nil, err
						}
						ps2, err = getProblemsOfHowAndType(2, typeName)
						if err != nil {
							return nil, err
						}
					}
				case 1:
					{
						// 填空题
						var err error
						ps1, err = getProblemsOfHowAndType(0, typeName)
						if err != nil {
							return nil, err
						}
						ps2, err = getProblemsOfHowAndType(2, typeName)
						if err != nil {
							return nil, err
						}
					}
				default:
					{
						// 其他题
						var err error
						ps1, err = getProblemsOfHowAndType(1, typeName)
						if err != nil {
							return nil, err
						}
						ps2, err = getProblemsOfHowAndType(0, typeName)
						if err != nil {
							return nil, err
						}
					}
				}

				if p, err := getOneProblemNotDone(id, ps1); err == nil {
					problemsForOneType = append(problemsForOneType, detailedProblem{
						ProblemID: p.ProblemID,
						SubIdx:    p.SubIdx,
						Full:      judgeFullHelper(p.SubIdx),
						howCode:   code,
					})
				} else {
					if p, err = getOneProblemNotDone(id, ps2); err == nil {
						problemsForOneType = append(problemsForOneType, detailedProblem{
							ProblemID: p.ProblemID,
							SubIdx:    p.SubIdx,
							Full:      judgeFullHelper(p.SubIdx),
							howCode:   code,
						})
					}
				}
			}

			if len(problemsForOneType) != 0 {
				// 之前没有找到未做过的且出题方式匹配的检验题，现在得到一道题即可
				break
			}
		}
	}

	// 还是找不到任何其它检验题，就从错题中选择1个题目当做检验题
	if len(problemsForOneType) == 0 {
		sort.Sort(wrongProblems)
		// wrongProblems[0]必然存在
		p := wrongProblems[0]
		problemsForOneType = append(problemsForOneType, detailedProblem{
			ProblemID: p.ProblemID,
			SubIdx:    p.SubIdx,
			Full:      judgeFullHelper(p.SubIdx),
			howCode:   getHowCodeInCheckProblems(p.how),
		})
	}

	return problemsForOneType, nil
}

func getCheckProblemsAlgorithm(wrongProblemsTmp []detailedProblem, id string) ([]detailedProblemsType, int, error) {
	// 获取检验题所发，wrongProblemsTmp是所有的错题，返回检验题和检验题数量

	// wrongProblems最后真正的错题（去除了内容库没有的题目）
	wrongProblems := []detailedProblem{}

	// 学生正在使用的bookID
	var bookIDsUsed []string
	bookIDsUsed, err := userDB.GetStudentBookIDs(id)
	if err != nil {
		bookIDsUsed = []string{}
	}

	// 补充题型信息与出题方式信息
	for i, p := range wrongProblemsTmp {
		if err := contentDB.ScanDetailedProblem(p.ProblemID, p.SubIdx, &wrongProblemsTmp[i], bookIDsUsed); err != nil {
			// 不对这道题寻找检验题
			log.Printf("Scaning details of problemID %s subIdx %d failed, err: %v\n", p.ProblemID, p.SubIdx, err)
			continue
		}
		var err error
		if wrongProblemsTmp[i].how, err = contentDB.GetHow(p.ProblemID); err != nil {
			// 不对这道题寻找检验题
			log.Printf("Getting how of problemID %s failed, err: %v\n", p.ProblemID, err)
			continue
		}
		wrongProblemsTmp[i].howCode = getHowCodeInCheckProblems(wrongProblemsTmp[i].how)
		wrongProblems = append(wrongProblems, wrongProblemsTmp[i])
	}

	// typeHow[0], [1], [2]分别表示选择，填空，其它，为true代表这个出题方式有题目
	type typeHow [3]bool
	typeHowMap := make(map[string]typeHow)
	typeProblemsMap := make(map[string]wrongProblemForCheckSlice)

	for _, p := range wrongProblems {
		if _, ok := typeHowMap[p.Type]; !ok {
			typeHowMap[p.Type] = typeHow{false, false, false}
		}
		th := typeHowMap[p.Type]
		th[p.howCode] = true
		typeHowMap[p.Type] = th

		if _, ok := typeProblemsMap[p.Type]; !ok {
			typeProblemsMap[p.Type] = wrongProblemForCheckSlice{}
		}
		typeProblemsMap[p.Type] = append(typeProblemsMap[p.Type], p)
	}

	// checkTypeProblems 存放每一个类型与对应找出来的检验题
	checkTypeProblems := checkProblemTypeSlice{}
	for _, ps := range typeProblemsMap {
		typeDetail := typeInfo{}
		if err := contentDB.GetTypeInfo(ps[0].ProblemID, ps[0].SubIdx, &typeDetail); err != nil {
			log.Printf("can not find type info of problemID %s subIdx %d, err: %v\n", ps[0].ProblemID, ps[0].SubIdx, err)
			continue
		}
		checkTypeProblems = append(checkTypeProblems, detailedProblemOfTypeInfo{
			Type:     typeDetail,
			Problems: []detailedProblem{},
		})
	}

	// 得到检验题目
	for index, pt := range checkTypeProblems {
		typeName := pt.Type.Type
		var problemsForOneType checkProblemSlice

		// 获取某一类型的检验题
		problemsForOneType, err := getCheckProblemsForOneType(id, typeName, typeHowMap[typeName], typeProblemsMap[typeName])
		if err != nil {
			return nil, 0, err
		}

		sort.Sort(problemsForOneType)

		// 补充具体题目信息
		for i, p := range problemsForOneType {
			if err := contentDB.ScanDetailedProblem(p.ProblemID, p.SubIdx, &problemsForOneType[i], bookIDsUsed); err != nil {
				log.Printf("scanning details of problemID %s, subIdx %d failed, err: %v\n", p.ProblemID, p.SubIdx, err)
				continue
			}
		}

		checkTypeProblems[index].Problems = problemsForOneType
	}

	sort.Sort(checkTypeProblems)

	// 添加序号
	typeProblemsList, totalNum := addIndex(checkTypeProblems)

	return typeProblemsList, totalNum, nil
}
