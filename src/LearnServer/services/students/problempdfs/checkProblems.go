package problempdfs

import (
	"log"
	"net/http"
	"strconv"
	"time"

	// "LearnServer/services/students/validation"

	// "LearnServer/models/userDB"
	// "LearnServer/utils"

	"LearnServer/services/students/validation"
	"LearnServer/models/userDB"
	"LearnServer/utils"
	
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

// GetCheckProblemsByChapSectHandler 根据章节获取检验题
func GetCheckProblemsByChapSectHandler(c echo.Context) error {
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

	// timeBefore 设置成一天之后，即获取错题范围是所有错题
	timeBefore := time.Now().AddDate(0, 0, 1)

	// 获取最新的错题
	tmpProblems, err := getNewestWrongProblemsOfChapSect(id, chapter, section, timeBefore)
	if err != nil {
		return err
	}

	typeProblemsList, totalNum, err := getCheckProblemsAlgorithm(tmpProblems, id)
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
		return utils.NotFound("No problems for checking in this chapter and this section.")
	}

	return c.JSON(http.StatusOK, result)
}

// GetAllCheckProblemsByBookPageHandler 根据书本和页码获取所有错过的题目的检验题目
func GetAllCheckProblemsByBookPageHandler(c echo.Context) error {
	var id string
	if err := validation.ValidateUser(c, &id); err != nil {
		return err
	}

	type bookPageInput struct {
		Book      string `json:"book"`
		StartPage int    `json:"startPage"`
		EndPage   int    `json:"endPage"`
	}

	var bookpages []bookPageInput
	if err := c.Bind(&bookpages); err != nil {
		return utils.InvalidParams("invalid inputs!")
	}

	// timeBefore 设置成一天之后，即获取错题范围是所有错题
	timeBefore := time.Now().AddDate(0, 0, 1)

	wrongProblems := []detailedProblem{}
	for _, bp := range bookpages {
		wrongProblemsTmp, err := GetOnceWrongProblemsOfBookPage(id, bp.Book, bp.StartPage, bp.EndPage, timeBefore)
		if err != nil {
			return err
		}
		wrongProblems = append(wrongProblems, wrongProblemsTmp...)
	}

	typeProblemsList, totalNum, err := getCheckProblemsAlgorithm(wrongProblems, id)
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
		return utils.NotFound("No problems for checking for the books and pages.")
	}

	return c.JSON(http.StatusOK, result)
}

// GetKnownCheckProblemsByBookPageHandler 根据书本和页码获取曾经错过现在会做的题目的检验题目
func GetKnownCheckProblemsByBookPageHandler(c echo.Context) error {
	var id string
	if err := validation.ValidateUser(c, &id); err != nil {
		return err
	}

	type bookPageInput struct {
		Book      string `json:"book"`
		StartPage int    `json:"startPage"`
		EndPage   int    `json:"endPage"`
	}

	var bookpages []bookPageInput
	if err := c.Bind(&bookpages); err != nil {
		return utils.InvalidParams("invalid inputs!")
	}

	wrongProblems := []detailedProblem{}
	for _, bp := range bookpages {
		wrongProblemsTmp, err := getKnownWrongProblemsOfBookPage(id, bp.Book, bp.StartPage, bp.EndPage)
		if err != nil {
			return err
		}
		wrongProblems = append(wrongProblems, wrongProblemsTmp...)
	}

	typeProblemsList, totalNum, err := getCheckProblemsAlgorithm(wrongProblems, id)
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
		return utils.NotFound("No problems for checking for the books and pages.")
	}

	return c.JSON(http.StatusOK, result)
}

// GetNewestCheckProblemsByBookPageHandler 根据书本和页码获取最新依然是错的题目的检验题目
func GetNewestCheckProblemsByBookPageHandler(c echo.Context) error {
	var id string
	if err := validation.ValidateUser(c, &id); err != nil {
		return err
	}

	type bookPageInput struct {
		Book      string `json:"book"`
		StartPage int    `json:"startPage"`
		EndPage   int    `json:"endPage"`
	}

	var bookpages []bookPageInput
	if err := c.Bind(&bookpages); err != nil {
		return utils.InvalidParams("invalid inputs!")
	}

	// timeBefore 设置成一天之后，即获取错题范围是所有错题
	timeBefore := time.Now().AddDate(0, 0, 1)

	wrongProblems := []detailedProblem{}
	for _, bp := range bookpages {
		wrongProblemsTmp, err := GetNewestWrongProblemsOfBookPage(id, bp.Book, bp.StartPage, bp.EndPage, timeBefore)
		if err != nil {
			return err
		}
		wrongProblems = append(wrongProblems, wrongProblemsTmp...)
	}

	typeProblemsList, totalNum, err := getCheckProblemsAlgorithm(wrongProblems, id)
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
		return utils.NotFound("No problems for checking for the books and pages.")
	}

	return c.JSON(http.StatusOK, result)
}

// UploadProblemCheckedResultHandler 上传检验题结果
func UploadProblemCheckedResultHandler(c echo.Context) error {
	type problemResult struct {
		IsCorrect  bool   `json:"isCorrect"`
		ProblemID  string `json:"problemID"`
		SubIdx     int    `json:"subIdx"`
		Smooth     int    `json:"smooth"`
		Understood int    `json:"understood"`
	}

	type uploadType struct {
		Time     int64           `json:"time"`
		Problems []problemResult `json:"problems"`
	}

	var id string
	if err := validation.ValidateUser(c, &id); err != nil {
		return err
	}

	var uploadData uploadType
	if err := c.Bind(&uploadData); err != nil {
		return utils.InvalidParams()
	}

	for _, p := range uploadData.Problems {
		err := userDB.C("students").UpdateId(bson.ObjectIdHex(id), bson.M{
			"$push": bson.M{
				"problems": bson.M{
					"assignDate": time.Unix(uploadData.Time, 0),
					"problemID":  p.ProblemID,
					"subIdx":     p.SubIdx,
					"correct":    p.IsCorrect,
					"smooth":     p.Smooth,
					"understood": p.Understood,
				},
			},
		})
		if err != nil {
			log.Println(err)
		}
	}
	return c.JSON(http.StatusOK, "Succeeded.")
}
