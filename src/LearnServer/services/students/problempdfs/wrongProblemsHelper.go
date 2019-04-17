package problempdfs

import (
	"log"
	"math"
	"sort"
	"strings"

	//"LearnServer/models/contentDB"

	"LearnServer/models/contentDB"
)

func getHowCodeInWrongProblems(how string) int {
	if how == "选择题" {
		return 1
	}
	if how == "填空题" {
		return 2
	}
	return 3
}

type wrongProblemSlice []detailedProblem

func (c wrongProblemSlice) Len() int {
	return len(c)
}

func (c wrongProblemSlice) Swap(i int, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c wrongProblemSlice) Less(i int, j int) bool {
	if c[i].sortType != 1 {
		// 不是根据出题方式排序
		ciHow := getHowCodeInWrongProblems(c[i].how)
		cjHow := getHowCodeInWrongProblems(c[j].how)
		if ciHow != cjHow {
			return ciHow < cjHow
		}

		if c[i].Page != c[j].Page {
			return c[i].Page < c[j].Page
		}
	}

	if c[i].Idx != c[j].Idx {
		return c[i].Idx < c[j].Idx
	}

	if strings.Compare(c[i].ProblemID, c[j].ProblemID) != 0 {
		return strings.Compare(c[i].ProblemID, c[j].ProblemID) < 0
	}

	return c[i].SubIdx < c[j].SubIdx
}

type wrongProblemTypeSlice []detailedProblemOfTypeInfo

func (c wrongProblemTypeSlice) Len() int {
	return len(c)
}

func (c wrongProblemTypeSlice) Swap(i int, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c wrongProblemTypeSlice) Less(i int, j int) bool {
	if c[i].Problems[0].sortType == 1 {
		// sortType 是 1 按出题方式时，type是空字符串，problems是同一道题的不同小问
		ciHow := getHowCodeInWrongProblems(c[i].Problems[0].how)
		cjHow := getHowCodeInWrongProblems(c[i].Problems[0].how)
		return ciHow < cjHow
	}

	if c[i].Type.Chapter != c[j].Type.Chapter {
		return c[i].Type.Chapter < c[j].Type.Chapter
	}

	if c[i].Type.Section != c[j].Type.Section {
		return c[i].Type.Section < c[j].Type.Section
	}

	return c[i].Type.Priority < c[j].Type.Priority
}

func createProblemTypeInfoList(problems []detailedProblem) ([]detailedProblemOfTypeInfo, error) {
	// 利用detailedProblem创建detailedProblemOfTypeInfo，将问题归类到不同类型下
	typeProblems := []detailedProblemOfTypeInfo{}

	for _, p := range problems {
		foundTypeInResult := false
		for j, pt := range typeProblems {
			if pt.Type.Type == p.Type {
				foundTypeInResult = true
				typeProblems[j].Problems = append(typeProblems[j].Problems, p)
				break
			}
		}
		if !foundTypeInResult {
			t := typeInfo{}
			err := contentDB.GetTypeInfo(p.ProblemID, p.SubIdx, &t)
			if err != nil {
				return []detailedProblemOfTypeInfo{}, err
			}
			typeProblems = append(typeProblems, detailedProblemOfTypeInfo{
				Type:     t,
				Problems: []detailedProblem{p},
			})
		}
	}

	return typeProblems, nil
}

func getPartOfProblems(problemTypes []detailedProblemOfTypeInfo, max int, sortType int) []detailedProblemOfTypeInfo {
	// 根据题量最大值获取一部分题目，problems应该是已经排序过的，sortType排序方式(1按出题方式，2按题目类型)
	if sortType == 1 {
		if max < len(problemTypes) {
			problemTypes = problemTypes[0:max]
		}
		return problemTypes
	}

	// probCount 统计每节对应有多少个problems，同一个problemID算一个
	probCount := []int{}
	total := 0
	chapterFormer := -1
	sectionFormer := -1
	// chapSectCount 统计一个节有多少题目
	chapSectCount := 0
	for _, pt := range problemTypes {
		probIDFormer := ""
		if pt.Type.Chapter != chapterFormer || pt.Type.Section != sectionFormer {
			chapterFormer = pt.Type.Chapter
			sectionFormer = pt.Type.Section
			if chapterFormer != -1 {
				probCount = append(probCount, chapSectCount)
				total += chapSectCount
			}
			chapSectCount = 0
		}

		for _, p := range pt.Problems {
			if p.ProblemID != probIDFormer {
				chapSectCount++
			}
			probIDFormer = p.ProblemID
		}
	}
	probCount = append(probCount, chapSectCount)
	total += chapSectCount

	if total <= max {
		return problemTypes
	}

	result := []detailedProblemOfTypeInfo{}

	chapterFormer = -1
	sectionFormer = -1
	// 总节数
	totalSect := len(probCount)
	// 第几节
	sectIndex := 0

	// remaining 剩下的节对应的题目总数量
	remaining := float64(total)
	// toGetRemaining 剩余需要获取的题目数量
	toGetRemaining := float64(max)

	// getProbsThisTime 这一节需要获取的题目数量
	var getProbsThisTime float64
	// 这一节已经获取的题目数量
	var probsGot float64

	// log.Print(remaining)
	// log.Print(total)
	for _, pt := range problemTypes {
		if pt.Type.Chapter != chapterFormer || pt.Type.Section != sectionFormer {
			chapterFormer = pt.Type.Chapter
			sectionFormer = pt.Type.Section
			if sectIndex != 0 {
				remaining -= float64(probCount[sectIndex])
				toGetRemaining -= probsGot
			}
			sectIndex++
			probsGot = 0
		}

		// log.Print("sectIndex ", sectIndex)
		// log.Print("probsGot ", probsGot)
		// log.Print("toGetRemaining", toGetRemaining)
		// log.Print("remaining", remaining)

		if sectIndex == totalSect {
			getProbsThisTime = toGetRemaining
		} else {
			getProbsThisTime = math.Trunc(toGetRemaining * float64(probCount[sectIndex]) / remaining)
		}
		if getProbsThisTime <= 0 {
			continue
		}

		// log.Print("getProbsThisTime", getProbsThisTime)
		newPt := detailedProblemOfTypeInfo{
			Type:     pt.Type,
			Problems: []detailedProblem{},
		}
		probIDFormer := ""
		for _, p := range pt.Problems {
			if probIDFormer != p.ProblemID {
				if probsGot >= getProbsThisTime {
					break
				}
				probsGot++
			}
			probIDFormer = p.ProblemID
			newPt.Problems = append(newPt.Problems, p)
		}
		if len(newPt.Problems) != 0 {
			result = append(result, newPt)
		}
	}
	return result
}

// GetWrongProblemsAlgorithm 根据初步挑选出来的错题获取纠错本题目的算法
func GetWrongProblemsAlgorithm(wrongProblems []detailedProblem, problemsForSelect []contentDB.DetailedProblem, max int, sortType int, bookIDs []string) ([]detailedProblemsType, int, error) {
	// 获取错题算法，wrongProblems是错题，problemsForSelect是包含所有错题在内的用来补充完整一道错题的题目，max是展示的最大题目数量，sortType排序方式(1按出题方式，2按题目类型), bookIDs 显示题目来源优先选择出现在其中的书本ID
	// 注意 wrongProblems 每个题目 full 应该初始状态是 true
	// 返回得到的错题和错题数量

	// problems最后真正的错题（去除了内容库没有的题目）
	problems := wrongProblemSlice{}

	// full 字段的含义：如果 full 是 true，代表这道题需要做完整个题目所有小问并提交，如果 full 是 false，代表只需要做 problems 中有的那些小问即可
	// full 应该初始状态是 true !! 只需要做特定小问的是特殊的

	// 补充how信息和具体题目信息
	for i, p := range wrongProblems {
		// isCal 是否为计算题，计算题不需要补全小问 (当没找到出题方式的时候默认补全)
		how, err := contentDB.GetHow(p.ProblemID)
		if err != nil {
			// 不展示这道题
			log.Printf("finding how of problemID %s failed, err: %v\n", p.ProblemID, err)
			continue
		}

		wrongProblems[i].how = how
		if wrongProblems[i].SubIdx != -1 && (how == "计算题" || how == "填空题") {
			wrongProblems[i].Full = false
		}

		// 补充具体题目信息
		if err := contentDB.ScanDetailedProblem(p.ProblemID, p.SubIdx, &wrongProblems[i], bookIDs); err != nil {
			// 不展示这道题
			log.Printf("scanning details of problemID %s, subIdx %d failed, err: %v\n", p.ProblemID, p.SubIdx, err)
			continue
		}

		wrongProblems[i].sortType = sortType
		problems = append(problems, wrongProblems[i])
	}

	// 当 full 为 false 的题目的所有小问都已经在 problems 中了，将 full 改为 true
	for i, p := range problems {
		if p.Full == true {
			continue
		}

		// findOtherProbs 是否能找到一道同 problemID 不同小问，而又不在 problems 中的题目
		findOtherProbs := false
		for _, ps := range problemsForSelect {
			if ps.ProblemID == p.ProblemID && p.SubIdx != ps.SubIdx {
				exist := false
				for _, pt := range problems {
					if pt.ProblemID == ps.ProblemID && pt.SubIdx == ps.SubIdx {
						exist = true
						break
					}
				}
				if !exist {
					findOtherProbs = true
					break
				}
			}
		}

		if !findOtherProbs {
			problems[i].Full = true
		}
	}

	// full 是 true 有小问的在problems中补全所有小问
	// 补充小问的原因：提交做题结果的时候需要，当文档中是一整道题的时候，做题结果需要提交所有小问
	for _, p := range problems {
		if p.SubIdx == -1 || p.Full == false {
			// 没必要补充小问
			continue
		}

		for _, pCS := range problemsForSelect {
			if p.ProblemID == pCS.ProblemID && p.SubIdx != pCS.SubIdx {
				// 是否已经有了，不用补充
				exist := false
				for _, ptmp := range problems {
					if ptmp.ProblemID == pCS.ProblemID && ptmp.SubIdx == pCS.SubIdx {
						exist = true
						break
					}
				}
				if !exist {
					tmpP := detailedProblem{
						ProblemID: pCS.ProblemID,
						SubIdx:    pCS.SubIdx,
						Full:      true,
						Reason:    "该小问与其它小问关系紧密",
						how:       p.how,
						sortType:  sortType,
					}
					if err := contentDB.ScanDetailedProblem(tmpP.ProblemID, tmpP.SubIdx, &tmpP, bookIDs); err != nil {
						log.Printf("scanning details of problemID %s, subIdx %d failed, err: %v\n", tmpP.ProblemID, tmpP.SubIdx, err)
						continue
					}
					problems = append(problems, tmpP)
				}
			}
		}
	}

	// 先排序，此时相同problemID都聚集在一起
	sort.Sort(problems)

	// 按照类型归类，虽然createProblemTypeInfoList返回值是[]detailedProblemOfTypeInfo，但赋值时有隐式类型转换
	// （必须声明typeProblemsRaw类型，否则sort.Sort()会出错（sort.Sort()接收的是接口，并没有类型转换，而[]detailedProblemOfTypeInfo没有实现该接口）
	var typeProblemsRaw wrongProblemTypeSlice

	if sortType == 1 {
		// 根据出题方式排序

		for i := range problems {
			// remove type information
			problems[i].Type = ""
		}

		problemIDFormer := ""
		for _, p := range problems {
			if p.ProblemID != problemIDFormer {
				tpTmp := detailedProblemOfTypeInfo{
					Type: typeInfo{
						Type: "",
					},
					Problems: []detailedProblem{p},
				}
				typeProblemsRaw = append(typeProblemsRaw, tpTmp)
				problemIDFormer = p.ProblemID
			} else {
				typeProblemsRaw[len(typeProblemsRaw)-1].Problems = append(typeProblemsRaw[len(typeProblemsRaw)-1].Problems, p)
			}
		}

	} else {
		// 此时problems中，同一个problemID的不同小问必定相连
		// pLast 上一个problem
		pLast := detailedProblem{}
		for index, p := range problems {
			if p.ProblemID == pLast.ProblemID {
				// 把不同小问的类型修正为相同类型，因为之后一道题是作为一个整体进行归类和排序的
				if p.Type != pLast.Type {
					problems[index].Type = pLast.Type
					p.Type = pLast.Type
				}
			}
			pLast = p
		}

		var err error
		typeProblemsRaw, err = createProblemTypeInfoList(problems)
		if err != nil {
			return nil, 0, err
		}
	}

	// typeProblemsRaw每个元素内部的problems已经是有序的
	sort.Sort(typeProblemsRaw)

	typeProblemsRaw = getPartOfProblems(typeProblemsRaw, max, sortType)

	// 添加序号
	typeProblemsList, totalNum := addIndex(typeProblemsRaw)

	return typeProblemsList, totalNum, nil
}
