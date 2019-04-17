package errorRateAnalysis

import (
	"log"
	"net/http"
	"time"

	"LearnServer/models/contentDB"
	"LearnServer/models/userDB"
	"LearnServer/services/students/problempdfs"
	"LearnServer/utils"
	"github.com/labstack/echo"
)

func findProblemInProblems(problems []problempdfs.DetailedProblem, problemID string, subIdx int) bool {
	// 在错题列表中寻找某道题的状态，true找到，false没找到
	for _, p := range problems {
		if p.ProblemID == problemID && p.SubIdx == subIdx {
			return true
		}
	}
	return false
}

// GetErrorRateHandler 获取某个班级某层级的错误率分析
func GetErrorRateHandler(c echo.Context) error {

	type bookPageInput struct {
		BookID    string `json:"bookID"`
		StartPage int    `json:"startPage"`
		EndPage   int    `json:"endPage"`
	}

	type inputType struct {
		WrongProblemStatus int             `json:"wrongProblemStatus"` // 错题状态，1现在仍错的题目，2曾经错过的
		BookPage           []bookPageInput `json:"bookPage"`
		PaperIDs           []string        `json:"paperIDs"`
		SchoolID           string          `json:"schoolID"`
		Grade              string          `json:"grade"`
		Class              int             `json:"class"`
		Level              int             `json:"level"`
		Exam               string          `json:"exam"`
		DateBefore         int64           `json:"dateBefore"`
	}

	var input inputType
	if err := c.Bind(&input); err != nil {
		return utils.InvalidParams("invalid inputs! error: " + err.Error())
	}

	type studentType struct {
		StudentDetail userDB.StudentType
		Problems      []problempdfs.DetailedProblem
	}

	studentsTmp, err := userDB.GetStudents(input.SchoolID, input.Grade, input.Class, input.Level, "", "", "", "")
	if err != nil {
		log.Printf("failed to get students, input %v, err %v\n", input, err)
		return err
	}

	students := make([]studentType, len(studentsTmp))
	for i, stu := range studentsTmp {
		students[i].StudentDetail = stu
	}

	timeBefore := time.Unix(input.DateBefore, 0)

	for i, stu := range students {
		wrongProblems := []problempdfs.DetailedProblem{}
		for _, bp := range input.BookPage {
			wrongProblemsTmp, err := problempdfs.GetNewestWrongProblemsOfBookPage(stu.StudentDetail.ID.Hex(), bp.BookID, bp.StartPage, bp.EndPage, timeBefore)
			if err != nil {
				return err
			}
			wrongProblems = append(wrongProblems, wrongProblemsTmp...)
		}

		for _, paperID := range input.PaperIDs {
			paperWrongProblemsTmp, err := problempdfs.GetNewestWrongPaperProblems(stu.StudentDetail.ID.Hex(), paperID, timeBefore)
			if err != nil {
				return err
			}
			wrongProblems = append(wrongProblems, paperWrongProblemsTmp...)
		}
		students[i].Problems = wrongProblems
	}

	allProblems := []contentDB.DetailedProblem{}
	for _, bp := range input.BookPage {
		problemsOfBookPage, err := contentDB.GetProblemsByBookPage(bp.BookID, bp.StartPage, bp.EndPage)
		if err != nil {
			return err
		}
		for i := range problemsOfBookPage {
			problemsOfBookPage[i].SourceID = bp.BookID
		}
		allProblems = append(allProblems, problemsOfBookPage...)
	}
	for _, paperID := range input.PaperIDs {
		problemsOfPaper, err := contentDB.GetProblemsByPaper(paperID)
		if err != nil {
			return err
		}
		allProblems = append(allProblems, problemsOfPaper...)
	}

	type resultType struct {
		Source        string   `json:"source"` // 来源， 书本名称或者试卷名称
		Page          int      `json:"page"`
		Column        string   `json:"column"`        // 栏目名称
		Idx           int      `json:"idx"`           // 题目序号
		ProblemID     string   `json:"problemID"`     // 题目识别码
		SubIdx        int      `json:"subIdx"`        // 小问序号
		Probability   float32  `json:"probability"`   // 考试概率
		ErrorRate     float32  `json:"errorRate"`     // 错误率
		WrongStudents []string `json:"wrongStudents"` // 错误的学生名单
		TotalStudents int      `json:"totalStudents"` // 分析的学生总数（因为选择了分析某个层级的学生，所以这里可能不等于班级学生总数）
	}

	results := []resultType{}

	for _, p := range allProblems {
		err := contentDB.ScanDetailedProblem(p.ProblemID, p.SubIdx, &p, []string{p.SourceID})
		if err != nil {
			log.Printf("failed to get detail of problem %v, error %v\n", p, err)
			continue
		}
		var probability float32
		var fieldName string
		switch input.Exam {
		case "单元考试":
			fieldName = "unitExamProb"
		case "期中考试":
			fieldName = "midtermProb"
		case "期末考试":
			fieldName = "finalProb"
		case "中考":
			fieldName = "seniorEntranceProb"
		default:
			fieldName = "finalProb"
		}
		sql := "SELECT " + fieldName + " FROM typeExamProbability WHERE typeName = ?;"
		if err := contentDB.GetDB().Get(&probability, sql, p.Type); err != nil {
			probability = 0.01
		}
		var wrongProblemResult resultType
		wrongProblemResult = resultType{
			Source:        p.Book,
			Page:          p.Page,
			Column:        p.Column,
			Idx:           p.Idx,
			ProblemID:     p.ProblemID,
			SubIdx:        p.SubIdx,
			Probability:   probability,
			TotalStudents: len(students),
			WrongStudents: []string{},
		}

		for _, stu := range students {
			if findProblemInProblems(stu.Problems, p.ProblemID, p.SubIdx) {
				wrongProblemResult.WrongStudents = append(wrongProblemResult.WrongStudents, stu.StudentDetail.Name)
			}
		}

		wrongProblemResult.ErrorRate = float32(len(wrongProblemResult.WrongStudents)) / float32(wrongProblemResult.TotalStudents)
		results = append(results, wrongProblemResult)
	}
	return c.JSON(http.StatusOK, results)
}
