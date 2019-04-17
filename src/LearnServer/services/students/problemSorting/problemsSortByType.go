package problemSorting

import (
	"log"
	"net/http"
	"sort"
	"strconv"

	// "LearnServer/models/contentDB"
	// "LearnServer/models/userDB"
	// "LearnServer/utils/classify"

	// "LearnServer/services/students/validation"
	// "LearnServer/utils"

	"LearnServer/models/contentDB"
	"LearnServer/models/userDB"
	"LearnServer/utils/classify"

	"LearnServer/services/students/validation"
	"LearnServer/utils"

	"github.com/labstack/echo"
)

// GetProblemsSortByTypeHandler 获取学生做过的题目按类型归类和排序
func GetProblemsSortByTypeHandler(c echo.Context) error {
	var id string
	if err := validation.ValidateUser(c, &id); err != nil {
		return err
	}

	chapter, err := strconv.Atoi(c.QueryParam("chapter"))
	if err != nil {
		return utils.InvalidParams("filter chapter is invalid.")
	}

	section, err := strconv.Atoi(c.QueryParam("section"))
	if err != nil {
		return utils.InvalidParams("filter section is invalid.")
	}

	// 学生正在使用的bookID
	var bookIDsUsed []string
	bookIDsUsed, err = userDB.GetStudentBookIDs(id)
	if err != nil {
		bookIDsUsed = []string{}
	}

	// 获取某个章节所有做过的题目
	var problemsDone toClassifyByType
	problemsDone, err = getProblemsDoneOfChapSect(id, chapter, section)
	if err != nil {
		return err
	}

	for i, p := range problemsDone {
		// 获取题型信息
		if err := contentDB.GetTypeInfo(p.ProblemID, p.SubIdx, &problemsDone[i].TypeInfo); err != nil {
			log.Printf("getting type info for ProblemID %s SubIdx %d failed. err: %v\n", p.ProblemID, p.SubIdx, err)
		}
	}
	// 按题型分类，然后对题型排序
	var classifyByTypeResult toSortByTypeForType
	classifyByTypeResult = constructProblemsByIndex(problemsDone, classify.Classify(problemsDone))
	sort.Sort(classifyByTypeResult)

	type detailedProblemResultType struct {
		Book       string `json:"book" db:"book"`
		Page       int    `json:"page" db:"page"`
		LessonName string `json:"lessonName" db:"lessonName"`
		Column     string `json:"column" db:"column"`
		ProblemID  string `json:"problemID" db:"problemID"`
		Idx        int    `json:"idx" db:"idx"`
		SubIdx     int    `json:"subIdx" db:"subIdx"`
		Correct    int    `json:"isCorrect"`
		Type       string `json:"-" db:"type"`
		Time       int64  `json:"assignDate"`
	}

	type oneTypeProblem struct {
		Type     string                      `json:"type"`
		Category string                      `json:"category"`
		Problems []detailedProblemResultType `json:"problems"`
	}

	result := make([]oneTypeProblem, 0, len(classifyByTypeResult))
	for _, ps := range classifyByTypeResult {

		// 移除同一个题型下重复的题目，只保留最新的，然后进行排序
		var psSorted toSortByTimeForType
		psSorted = removeDuplicateProblems(ps)
		sort.Sort(psSorted)

		// problemsForOneTypeNew是修改过correct的，将被添加到resultType中
		problemsForOneTypeNew := oneTypeProblem{
			Type:     ps[0].TypeInfo.Type,
			Category: ps[0].TypeInfo.Category,
			Problems: make([]detailedProblemResultType, 0, len(psSorted)),
		}

		for _, p := range psSorted {
			detailedP := detailedProblemResultType{
				ProblemID: p.ProblemID,
				SubIdx:    p.SubIdx,
				Time:      p.Time.Unix(),
				Correct:   getCorrectStatus(ps, p.ProblemID, p.SubIdx, p.Time),
			}
			if err := contentDB.ScanDetailedProblem(p.ProblemID, p.SubIdx, &detailedP, bookIDsUsed); err != nil {
				log.Printf("scanning details of problemID %s, subIdx %d failed, err: %v \n", p.ProblemID, p.SubIdx, err)
				// 不展示这道题
				continue
			}

			problemsForOneTypeNew.Problems = append(problemsForOneTypeNew.Problems, detailedP)
		}

		if len(problemsForOneTypeNew.Problems) != 0 {
			result = append(result, problemsForOneTypeNew)
		}
	}

	// 只选择后10道题目
	// for index, tps := range result {
	// 	if len(tps.Problems) > 10 {
	// 		// 此处赋值给result[index]而非tps，修改tps并不会改变result
	// 		result[index].Problems = tps.Problems[len(tps.Problems)-10:]
	// 	}
	// }

	return c.JSON(http.StatusOK, result)
}
