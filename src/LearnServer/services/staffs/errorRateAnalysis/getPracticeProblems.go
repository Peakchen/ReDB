package errorRateAnalysis

import (
	"net/http"

	"LearnServer/models/contentDB"
	
	"LearnServer/services/students/problempdfs"
	"LearnServer/utils"
	"github.com/labstack/echo"
)

// GetPracticeProblemsHandler 根据预选择的题目获取真正的训练题目
func GetPracticeProblemsHandler(c echo.Context) error {



	type bookPageInput struct {
		BookID    string `json:"bookID"`
		StartPage int    `json:"startPage"`
		EndPage   int    `json:"endPage"`
	}

	type inputType struct {
		BookPage []bookPageInput `json:"bookPage"`
		PaperIDs []string        `json:"paperIDs"`
		Problems []struct {
			ProblemID string `json:"problemID"`
			SubIdx    int    `json:"subIdx"`
		} `json:"problems"`
	}

	var input inputType
	if err := c.Bind(&input); err != nil {
		return utils.InvalidParams("invalid inputs!")
	}

	const max int = 8
	const sortType int = 1

	wrongProblems := make([]problempdfs.DetailedProblem, len(input.Problems))
	for i, p := range input.Problems {
		wrongProblems[i] = problempdfs.DetailedProblem{
			ProblemID: p.ProblemID,
			SubIdx:    p.SubIdx,
		}
	}

	problemsForSelect := []contentDB.DetailedProblem{}
	for _, bp := range input.BookPage {
		problemsOfBookPage, err := contentDB.GetProblemsByBookPage(bp.BookID, bp.StartPage, bp.EndPage)
		if err != nil {
			return err
		}
		problemsForSelect = append(problemsForSelect, problemsOfBookPage...)
	}
	for _, paperID := range input.PaperIDs {
		problemsOfPaper, err := contentDB.GetProblemsByPaper(paperID)
		if err != nil {
			return err
		}
		problemsForSelect = append(problemsForSelect, problemsOfPaper...)
	}

	bookIDs := make([]string, len(input.BookPage))
	for i, b := range input.BookPage {
		bookIDs[i] = b.BookID
	}

	typeProblemsList, totalNum, err := problempdfs.GetWrongProblemsAlgorithm(wrongProblems, problemsForSelect, max, sortType, bookIDs)
	if err != nil {
		return err
	}

	result := struct {
		TotalNum int                                `json:"totalNum"`
		Problems []problempdfs.DetailedProblemsType `json:"wrongProblems"`
	}{
		TotalNum: totalNum,
		Problems: typeProblemsList,
	}

	if result.TotalNum == 0 {
		return utils.NotFound("can not find practice problems")
	}

	return c.JSON(http.StatusOK, result)
}
