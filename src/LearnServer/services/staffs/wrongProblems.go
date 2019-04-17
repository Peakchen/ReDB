package staffs

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"LearnServer/models/contentDB"
	"LearnServer/models/userDB"
	"LearnServer/services/students/problempdfs"
	"LearnServer/utils"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

type bookPageInput struct {
	BookID    string `json:"bookID"`
	StartPage int    `json:"startPage"`
	EndPage   int    `json:"endPage"`
}

func getWrongProblemsHandler(c echo.Context) error {
	// 获取学生错题

	learnID, err := strconv.Atoi(c.Param("learnID"))
	if err != nil {
		return utils.InvalidParams("learnID is invalid")
	}
	studentID, err := getStudentIDByLearnID(learnID)
	if err != nil {
		return utils.NotFound("can not find the information of this learnID")
	}

	type inputType struct {
		ProductID        string          `json:"productID"`
		WrongProblemType int             `json:"wrongProblemType"`
		Sort             int             `json:"sort"`
		BatchID          string          `json:"batchID"`
		BookPage         []bookPageInput `json:"bookPage"`
		PaperIDs         []string        `json:"paperIDs"`
	}

	var input inputType
	if err := c.Bind(&input); err != nil {
		return utils.InvalidParams("invalid inputs!, err: " + err.Error())
	}

	batchInfo, ok := batchDownloadMap[input.BatchID]
	if !ok {
		return utils.InvalidParams("this batchID doesn't exist!")
	}
	batchInfo.lastFileTime = time.Now()
	studentIndex, err := getStudentIndex(input.BatchID, learnID)
	if err != nil {
		return utils.InvalidParams(err.Error())
	}

	epu, max, err := problempdfs.GetEpuAndMaxProblems(input.ProductID)
	if err != nil {
		log.Printf("can not get parameter 'max' of this student, id: %s, err: %v\n", studentID, err)
		return utils.NotFound("can not get parameter 'max' of this student")
	}

	var problems []problempdfs.ProblemForCreatingFiles

	switch {
	case epu == 1 && input.WrongProblemType == 1:
		problems, err = getNewestWrongProblems(studentID, input.Sort, max, input.BookPage, input.PaperIDs)
	case epu == 1 && input.WrongProblemType == 2:
		problems, err = getOnceWrongProblems(studentID, input.Sort, max, input.BookPage, input.PaperIDs)
	case epu == 2:
		problems, err = autoGetWrongProblems(studentID, input.Sort, max)
	}
	if err != nil {
		log.Printf("failed to get problems, err %v\n", err)
		return err
	}

	if len(problems) == 0 {
		batchInfo.Students[studentIndex].Problems = problems
		batchInfo.Students[studentIndex].ProblemFileStatus = false
		batchInfo.Students[studentIndex].AnswerFileStatus = false
		batchInfo.Students[studentIndex].ProblemStatusCode = 400
		batchInfo.Students[studentIndex].AnswerStatusCode = 400
		batchDownloadMap[input.BatchID] = batchInfo
		return utils.NotFound("No wrong problems in the books pages selected.")
	}

	problemsBytes, err := json.Marshal(problems)
	if err != nil {
		log.Printf("%v, err: %v", problems, err)
		return err
	}
	err = userDB.C("students").UpdateId(bson.ObjectIdHex(studentID), bson.M{
		"$set": bson.M{
			"lastWrongProblemsStr": string(problemsBytes),
		},
	})
	if err != nil {
		log.Printf("save lastWrongProblemsStr of student id %s failed, err : %v", studentID, err)
	}

	batchInfo.Students[studentIndex].Problems = problems
	batchDownloadMap[input.BatchID] = batchInfo

	return c.JSON(http.StatusOK, problems)
}

// autoGetWrongProblems 自动获取用来生成纠错本的错题（EPU2）
func autoGetWrongProblems(studentID string, sort int, max int) ([]problempdfs.ProblemForCreatingFiles, error) {
	wrongProblemsPicked, err := problempdfs.PickWrongProblems(studentID, max)
	if err != nil {
		return nil, err
	}
	problemsForSelect, err := contentDB.GetAllProblems()
	if err != nil {
		return nil, err
	}

	bookIDs, err := userDB.GetStudentBookIDs(studentID)
	if err != nil {
		return nil, err
	}

	typeProblemsList, _, err := problempdfs.GetWrongProblemsAlgorithm(wrongProblemsPicked, problemsForSelect, max, sort, bookIDs)
	if err != nil {
		return nil, err
	}

	return problempdfs.ConvertToNewFormat(typeProblemsList), nil
}

// getNewestWrongProblems 获取最新做错的题目
func getNewestWrongProblems(studentID string, sort int, max int, bookPage []bookPageInput, paperIDs []string) ([]problempdfs.ProblemForCreatingFiles, error) {

	// timeBefore 设置成一天之后，即获取错题范围是所有错题
	timeBefore := time.Now().AddDate(0, 0, 1)

	bookIDs, err := userDB.GetStudentBookIDs(studentID)
	if err != nil {
		return nil, err
	}

	wrongProblems := []problempdfs.DetailedProblem{}
	problemsForSelect := []contentDB.DetailedProblem{}
	for _, bp := range bookPage {
		wrongProblemsTmp, err := problempdfs.GetNewestWrongProblemsOfBookPage(studentID, bp.BookID, bp.StartPage, bp.EndPage, timeBefore)
		if err != nil {
			return nil, err
		}

		problemsOfBookPage, err := contentDB.GetProblemsByBookPage(bp.BookID, bp.StartPage, bp.EndPage)
		if err != nil {
			return nil, err
		}
		wrongProblems = append(wrongProblems, wrongProblemsTmp...)
		problemsForSelect = append(problemsForSelect, problemsOfBookPage...)
	}

	for _, paperID := range paperIDs {
		paperWrongProblemsTmp, err := problempdfs.GetNewestWrongPaperProblems(studentID, paperID, timeBefore)
		if err != nil {
			return nil, err
		}

		problemsOfPaper, err := contentDB.GetProblemsByPaper(paperID)
		if err != nil {
			return nil, err
		}
		wrongProblems = append(wrongProblems, paperWrongProblemsTmp...)
		problemsForSelect = append(problemsForSelect, problemsOfPaper...)
	}

	typeProblemsList, _, err := problempdfs.GetWrongProblemsAlgorithm(wrongProblems, problemsForSelect, max, sort, bookIDs)
	if err != nil {
		return nil, err
	}

	return problempdfs.ConvertToNewFormat(typeProblemsList), nil
}

// getOnceWrongProblems 获取曾经错过的所有题目
func getOnceWrongProblems(studentID string, sort int, max int, bookPage []bookPageInput, paperIDs []string) ([]problempdfs.ProblemForCreatingFiles, error) {

	// timeBefore 设置成一天之后，即获取错题范围是所有错题
	timeBefore := time.Now().AddDate(0, 0, 1)

	bookIDs, err := userDB.GetStudentBookIDs(studentID)
	if err != nil {
		return nil, err
	}

	wrongProblems := []problempdfs.DetailedProblem{}
	problemsForSelect := []contentDB.DetailedProblem{}
	for _, bp := range bookPage {
		wrongProblemsTmp, err := problempdfs.GetOnceWrongProblemsOfBookPage(studentID, bp.BookID, bp.StartPage, bp.EndPage, timeBefore)
		if err != nil {
			return nil, err
		}

		problemsOfBookPage, err := contentDB.GetProblemsByBookPage(bp.BookID, bp.StartPage, bp.EndPage)
		if err != nil {
			return nil, err
		}
		wrongProblems = append(wrongProblems, wrongProblemsTmp...)
		problemsForSelect = append(problemsForSelect, problemsOfBookPage...)
	}

	for _, paperID := range paperIDs {
		paperWrongProblemsTmp, err := problempdfs.GetOncePaperWrongProblems(studentID, paperID, timeBefore)
		if err != nil {
			return nil, err
		}

		problemsOfPaper, err := contentDB.GetProblemsByPaper(paperID)
		if err != nil {
			return nil, err
		}
		wrongProblems = append(wrongProblems, paperWrongProblemsTmp...)
		problemsForSelect = append(problemsForSelect, problemsOfPaper...)
	}

	typeProblemsList, _, err := problempdfs.GetWrongProblemsAlgorithm(wrongProblems, problemsForSelect, max, sort, bookIDs)
	if err != nil {
		return nil, err
	}

	return problempdfs.ConvertToNewFormat(typeProblemsList), nil
}
