package problempdfs

import (
	"net/http"
	"time"

	// "LearnServer/models/contentDB"
	// "LearnServer/services/students/validation"
	// "LearnServer/utils"

	"LearnServer/models/contentDB"
	"LearnServer/services/students/validation"
	"LearnServer/utils"
	"github.com/labstack/echo"
)

// GetNewestWrongPaperProblemsHandler 根据试卷识别码获取最新做错的题目，并按特定方式展示
func GetNewestWrongPaperProblemsHandler(c echo.Context) error {

	var id string
	if err := validation.ValidateUser(c, &id); err != nil {
		return err
	}

	type inputType struct {
		Sort     int      `json:"sort"`
		Paper    int      `json:"paper"`
		Max      int      `json:"max"`
		PaperIDs []string `json:"paperIDs"`
	}

	var input inputType
	if err := c.Bind(&input); err != nil {
		return utils.InvalidParams("invalid inputs!")
	}

	// timeBefore 设置成一天之后，即获取错题范围是所有错题
	timeBefore := time.Now().AddDate(0, 0, 1)

	wrongProblems := []detailedProblem{}
	problemsForSelect := []contentDB.DetailedProblem{}
	for _, paperID := range input.PaperIDs {
		wrongProblemsTmp, err := GetNewestWrongPaperProblems(id, paperID, timeBefore)
		if err != nil {
			return err
		}

		problemsOfPaper, err := contentDB.GetProblemsByPaper(paperID)
		if err != nil {
			return err
		}
		wrongProblems = append(wrongProblems, wrongProblemsTmp...)
		problemsForSelect = append(problemsForSelect, problemsOfPaper...)
	}

	typeProblemsList, totalNum, err := GetWrongProblemsAlgorithm(wrongProblems, problemsForSelect, input.Max, input.Sort, []string{})
	if err != nil {
		return err
	}

	result := struct {
		TotalNum int                    `json:"totalNum"`
		Problems []detailedProblemsType `json:"wrongProblems"`
	}{
		TotalNum: totalNum,
		Problems: typeProblemsList,
	}

	if result.TotalNum == 0 {
		return utils.NotFound("No wrong problems in the paper selected.")
	}

	return c.JSON(http.StatusOK, result)
}

// GetOnceWrongPaperProblemsHandler 根据试卷识别码获取曾经错过的所有题目，并按特定方式展示
func GetOnceWrongPaperProblemsHandler(c echo.Context) error {

	var id string
	if err := validation.ValidateUser(c, &id); err != nil {
		return err
	}

	type inputType struct {
		Sort     int      `json:"sort"`
		Paper    int      `json:"paper"`
		Max      int      `json:"max"`
		PaperIDs []string `json:"paperIDs"`
	}

	var input inputType
	if err := c.Bind(&input); err != nil {
		return utils.InvalidParams("invalid inputs!")
	}

	// timeBefore 设置成一天之后，即获取错题范围是所有错题
	timeBefore := time.Now().AddDate(0, 0, 1)

	wrongProblems := []detailedProblem{}
	problemsForSelect := []contentDB.DetailedProblem{}
	for _, paperID := range input.PaperIDs {
		wrongProblemsTmp, err := GetOncePaperWrongProblems(id, paperID, timeBefore)
		if err != nil {
			return err
		}

		problemsOfPaper, err := contentDB.GetProblemsByPaper(paperID)
		if err != nil {
			return err
		}

		wrongProblems = append(wrongProblems, wrongProblemsTmp...)
		problemsForSelect = append(problemsForSelect, problemsOfPaper...)
	}

	typeProblemsList, totalNum, err := GetWrongProblemsAlgorithm(wrongProblems, problemsForSelect, input.Max, input.Sort, []string{})
	if err != nil {
		return err
	}

	result := struct {
		TotalNum int                    `json:"totalNum"`
		Problems []detailedProblemsType `json:"wrongProblems"`
	}{
		TotalNum: totalNum,
		Problems: typeProblemsList,
	}

	if result.TotalNum == 0 {
		return utils.NotFound("No wrong problems in the paper selected.")
	}

	return c.JSON(http.StatusOK, result)
}
