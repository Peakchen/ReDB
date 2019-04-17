package problemSorting

import (
	"log"
	"net/http"
	"sort"
	"strconv"

	// "LearnServer/models/contentDB"
	// "LearnServer/models/userDB"

	// "LearnServer/services/students/validation"
	// "LearnServer/utils"
	// "LearnServer/utils/classify"

	"LearnServer/models/contentDB"
	"LearnServer/models/userDB"

	"LearnServer/services/students/validation"
	"LearnServer/utils"
	"LearnServer/utils/classify"

	"github.com/labstack/echo"
)

// GetProblemsSortByTimeHandler 获取按作业布置时间归类排序的做过的题目
func GetProblemsSortByTimeHandler(c echo.Context) error {
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
	var problemsDone toClassifyByTime
	problemsDone, err = getProblemsDoneOfChapSect(id, chapter, section)
	if err != nil {
		return err
	}

	// 根据时间分类
	var classifyByTimeResult toSortByTimeForTime
	classifyByTimeResult = constructProblemsByIndex(problemsDone, classify.Classify(problemsDone))

	// 按时间排序
	sort.Sort(classifyByTimeResult)

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
		Category   string `json:"category"`
	}

	type oneDayProblemsType struct {
		Time     int64                       `json:"time"`
		Problems []detailedProblemResultType `json:"problems"`
	}
	result := make([]oneDayProblemsType, 0, len(classifyByTimeResult))

	for _, ps := range classifyByTimeResult {
		for pi, p := range ps {
			// 获取题型信息
			if err := contentDB.GetTypeInfo(p.ProblemID, p.SubIdx, &ps[pi].TypeInfo); err != nil {
				log.Printf("getting type info for ProblemID %s SubIdx %d failed. err: %v\n", p.ProblemID, p.SubIdx, err)
			}
		}

		// 每天内部的题目排序
		var problemsOfOneDay toSortByTypeForTime
		problemsOfOneDay = ps
		sort.Sort(problemsOfOneDay)

		// psOneDayNew是修改过correct的，将被添加到result中
		psOneDayNew := oneDayProblemsType{
			Time:     problemsOfOneDay[0].Time.Unix(),
			Problems: make([]detailedProblemResultType, 0, len(problemsOfOneDay)),
		}

		for _, p := range problemsOfOneDay {
			detailedP := detailedProblemResultType{
				ProblemID: p.ProblemID,
				SubIdx:    p.SubIdx,
				Category:  p.TypeInfo.Category,
			}
			if err := contentDB.ScanDetailedProblem(p.ProblemID, p.SubIdx, &detailedP, bookIDsUsed); err != nil {
				log.Printf("scanning details of problemID %s, subIdx %d failed, err: %v \n", p.ProblemID, p.SubIdx, err)
				// 不展示这道题
				continue
			}
			// tmps 只是为了类型转换
			var tmps []problemWithTypeTimeAndCorrectDB
			tmps = problemsDone
			detailedP.Correct = getCorrectStatus(tmps, p.ProblemID, p.SubIdx, p.Time)
			psOneDayNew.Problems = append(psOneDayNew.Problems, detailedP)
		}

		if len(psOneDayNew.Problems) != 0 {
			result = append(result, psOneDayNew)
		}
	}

	return c.JSON(http.StatusOK, result)
}
