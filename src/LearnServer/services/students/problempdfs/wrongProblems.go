package problempdfs

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	// "LearnServer/models/contentDB"
	// "LearnServer/models/userDB"
	// "LearnServer/services/students/validation"
	// "LearnServer/utils"

	"LearnServer/models/contentDB"
	"LearnServer/models/userDB"
	"LearnServer/services/students/validation"
	"LearnServer/utils"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

// GetEpuAndMaxProblems 获取用户产品EPU和设定的最大题量参数
func GetEpuAndMaxProblems(productID string) (int, int, error) {
	problemMax := struct {
		EPU int `bson:"epu"`
		Max int `bson:"problemMax"`
	}{}
	err := userDB.C("products").Find(bson.M{
		"productID": productID,
	}).Select(bson.M{
		"problemMax": 1,
		"epu":        1,
	}).One(&problemMax)
	if err != nil {
		return 0, 0, err
	}

	return problemMax.EPU, problemMax.Max, nil
}

// GetWrongProblemsForLearning 自动获取用来生成纠错本的错题
func GetWrongProblemsForLearning(c echo.Context) error {
	var id string
	if err := validation.ValidateUser(c, &id); err != nil {
		return err
	}

	productID := c.QueryParam("productID")

	_, max, err := GetEpuAndMaxProblems(productID)
	if err != nil {
		log.Printf("can not get parameter 'max' of this student, id: %s, err: %v\n", id, err)
		return utils.NotFound("can not get parameter 'max' of this student")
	}

	sort, err := strconv.Atoi(c.QueryParam("sort"))
	if err != nil {
		return utils.InvalidParams("sort is invalid.")
	}

	wrongProblemsPicked, err := PickWrongProblems(id, max)
	if err != nil {
		return err
	}
	problemsForSelect, err := contentDB.GetAllProblems()
	if err != nil {
		return err
	}

	bookIDs, err := userDB.GetStudentBookIDs(id)
	if err != nil {
		return err
	}

	typeProblemsList, totalNum, err := GetWrongProblemsAlgorithm(wrongProblemsPicked, problemsForSelect, max, sort, bookIDs)
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
		return utils.NotFound("No wrong problems in the books pages selected.")
	}

	resultBytes, err := json.Marshal(result)
	if err != nil {
		log.Printf("%v, err: %v", result, err)
		return err
	}
	err = userDB.C("students").UpdateId(bson.ObjectIdHex(id), bson.M{
		"$set": bson.M{
			"lastWrongProblemsStr": string(resultBytes),
		},
	})
	if err != nil {
		log.Printf("save lastWrongProblemsStr of id %s failed, err : %v", id, err)
	}

	return c.JSON(http.StatusOK, result)
}

// GetNewestWrongProblemsByChapSectHandler 根据章节获取最新做错的题目，并按特定方式展示
func GetNewestWrongProblemsByChapSectHandler(c echo.Context) error {

	var id string
	if err := validation.ValidateUser(c, &id); err != nil {
		return err
	}

	productID := c.QueryParam("productID")

	chapter, err := strconv.Atoi(c.QueryParam("chapter"))
	if err != nil {
		return utils.InvalidParams("filter chapter is invalid.")
	}

	section, err := strconv.Atoi(c.QueryParam("section"))
	if err != nil {
		return utils.InvalidParams("filter section is invalid.")
	}

	bookIDs, err := userDB.GetStudentBookIDs(id)
	if err != nil {
		return err
	}

	_, max, err := GetEpuAndMaxProblems(productID)
	if err != nil {
		log.Printf("can not get parameter 'max' of this student, id: %s, err: %v\n", id, err)
		return utils.NotFound("can not get parameter 'max' of this student")
	}

	// timeBefore 设置成一天之后，即获取错题范围是所有错题
	timeBefore := time.Now().AddDate(0, 0, 1)

	tmpProblems, err := getNewestWrongProblemsOfChapSect(id, chapter, section, timeBefore)
	if err != nil {
		return err
	}

	problemsOfChapAndSect, err := contentDB.GetProblemsByChapterSection(chapter, section)
	if err != nil {
		return err
	}

	typeProblemsList, totalNum, err := GetWrongProblemsAlgorithm(tmpProblems, problemsOfChapAndSect, max, 2, bookIDs)
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
		return utils.NotFound("No wrong problems in this chapter and this section.")
	}

	return c.JSON(http.StatusOK, result)
}

// GetNewestWrongProblemsByBookPageHandler 根据书本页码获取最新做错的题目，并按特定方式展示
func GetNewestWrongProblemsByBookPageHandler(c echo.Context) error {

	var id string
	if err := validation.ValidateUser(c, &id); err != nil {
		return err
	}

	type bookPageInput struct {
		BookID    string `json:"bookID"`
		StartPage int    `json:"startPage"`
		EndPage   int    `json:"endPage"`
	}

	type inputType struct {
		Sort      int             `json:"sort"`
		ProductID string          `json:"productID"`
		Paper     int             `json:"paper"`
		BookPage  []bookPageInput `json:"bookPage"`
	}

	var input inputType
	if err := c.Bind(&input); err != nil {
		return utils.InvalidParams("invalid inputs!")
	}

	_, max, err := GetEpuAndMaxProblems(input.ProductID)
	if err != nil {
		log.Printf("can not get parameter 'max' of this student, id: %s, err: %v\n", id, err)
		return utils.NotFound("can not get parameter 'max' of this student")
	}

	bookIDs, err := userDB.GetStudentBookIDs(id)
	if err != nil {
		return err
	}

	// timeBefore 设置成一天之后，即获取错题范围是所有错题
	timeBefore := time.Now().AddDate(0, 0, 1)

	wrongProblems := []detailedProblem{}
	problemsForSelect := []contentDB.DetailedProblem{}
	for _, bp := range input.BookPage {
		wrongProblemsTmp, err := GetNewestWrongProblemsOfBookPage(id, bp.BookID, bp.StartPage, bp.EndPage, timeBefore)
		if err != nil {
			return err
		}

		problemsOfBookPage, err := contentDB.GetProblemsByBookPage(bp.BookID, bp.StartPage, bp.EndPage)
		if err != nil {
			return err
		}
		wrongProblems = append(wrongProblems, wrongProblemsTmp...)
		problemsForSelect = append(problemsForSelect, problemsOfBookPage...)
	}

	typeProblemsList, totalNum, err := GetWrongProblemsAlgorithm(wrongProblems, problemsForSelect, max, input.Sort, bookIDs)
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
		return utils.NotFound("No wrong problems in the books pages selected.")
	}

	return c.JSON(http.StatusOK, result)
}

// GetOnceWrongProblemsByBookPageHandler 根据书本页码获取曾经错过的所有题目，并按特定方式展示
func GetOnceWrongProblemsByBookPageHandler(c echo.Context) error {

	var id string
	if err := validation.ValidateUser(c, &id); err != nil {
		return err
	}

	type bookPageInput struct {
		BookID    string `json:"bookID"`
		StartPage int    `json:"startPage"`
		EndPage   int    `json:"endPage"`
	}

	type inputType struct {
		Sort      int             `json:"sort"`
		ProductID string          `json:"productID"`
		Paper     int             `json:"paper"`
		BookPage  []bookPageInput `json:"bookPage"`
	}

	var input inputType
	if err := c.Bind(&input); err != nil {
		return utils.InvalidParams("invalid inputs!")
	}

	_, max, err := GetEpuAndMaxProblems(input.ProductID)
	if err != nil {
		log.Printf("can not get parameter 'max' of this student, id: %s, err: %v\n", id, err)
		return utils.NotFound("can not get parameter 'max' of this student")
	}

	bookIDs, err := userDB.GetStudentBookIDs(id)
	if err != nil {
		return err
	}

	// timeBefore 设置成一天之后，即获取错题范围是所有错题
	timeBefore := time.Now().AddDate(0, 0, 1)

	wrongProblems := []detailedProblem{}
	problemsForSelect := []contentDB.DetailedProblem{}
	for _, bp := range input.BookPage {
		wrongProblemsTmp, err := GetOnceWrongProblemsOfBookPage(id, bp.BookID, bp.StartPage, bp.EndPage, timeBefore)
		if err != nil {
			return err
		}

		problemsOfBookPage, err := contentDB.GetProblemsByBookPage(bp.BookID, bp.StartPage, bp.EndPage)
		if err != nil {
			return err
		}
		wrongProblems = append(wrongProblems, wrongProblemsTmp...)
		problemsForSelect = append(problemsForSelect, problemsOfBookPage...)
	}

	typeProblemsList, totalNum, err := GetWrongProblemsAlgorithm(wrongProblems, problemsForSelect, max, input.Sort, bookIDs)
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
		return utils.NotFound("No wrong problems in the books pages selected.")
	}

	return c.JSON(http.StatusOK, result)
}

// UploadProblemRevisedResultHandler 上传错题复习结果
func UploadProblemRevisedResultHandler(c echo.Context) error {
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

// GetLastWrongProblemsHandler 获取上一次请求的用来生成纠错本的错题数据
func GetLastWrongProblemsHandler(c echo.Context) error {
	var id string
	if err := validation.ValidateUser(c, &id); err != nil {
		return err
	}

	data := struct {
		LastWrongProblemsStr string `bson:"lastWrongProblemsStr"`
	}{}

	err := userDB.C("students").FindId(bson.ObjectIdHex(id)).Select(bson.M{
		"lastWrongProblemsStr": 1,
	}).One(&data)

	if err != nil {
		log.Printf("failed to get lastWrongProblemsStr of id %s, err: %v", id, err)
		return err
	}

	lastWrongProblems := []byte(data.LastWrongProblemsStr)
	return c.JSONBlob(http.StatusOK, lastWrongProblems)
}
