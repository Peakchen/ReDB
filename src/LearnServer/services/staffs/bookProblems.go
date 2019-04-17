package staffs

import (
	"net/http"
	"strconv"

	// "LearnServer/models/contentDB"
	// "LearnServer/models/userDB"
	// "LearnServer/utils"

	"LearnServer/models/contentDB"
	"LearnServer/models/userDB"
	"LearnServer/utils"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

func getBookProblemsHandler(c echo.Context) error {
	// 获取某一本书某一页的未做过的题目
	bookID := c.QueryParam("book")
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		return utils.InvalidParams("filter page is invalid")
	}

	learnID, err := strconv.Atoi(c.Param("learnID"))
	if err != nil {
		return utils.InvalidParams("learnID is invalid")
	}
	studentID, err := getStudentIDByLearnID(learnID)
	if err != nil {
		return utils.NotFound("can not find the information of this learnID")
	}

	problems, err := contentDB.GetNonExampleProblemsByBookAndPage(bookID, page)
	if err != nil {
		return err
	}
	if len(problems) == 0 {
		return utils.NotFound("no problems in this page in the book")
	}

	type problemType struct {
		ProblemID  string `json:"problemID" bson:"problemID"`
		SubIdx     int    `json:"subIdx" bson:"subIdx"`
		SourceID   string `json:"-" bson:"sourceID"`
		SourceType int    `json:"-" bson:"sourceType"`
	}

	allProblemsDone := struct {
		Problems []problemType `bson:"problems"`
	}{}

	err = userDB.C("students").FindId(bson.ObjectIdHex(studentID)).Select(bson.M{
		"problems": 1,
	}).One(&allProblemsDone)
	if err != nil {
		return err
	}

	var result = make([]contentDB.DetailedProblem, len(problems))
	copy(result, problems[:])
	// 复制，避免在迭代过程中删除元素
	// 不能使用var result = problems[:]，不是深复制

	// 去除已经做过的题目
	for _, pResp := range problems {
		hasDone := false
		for _, pDone := range allProblemsDone.Problems {
			if pDone.SourceType == 1 && pDone.SourceID == bookID && pDone.ProblemID == pResp.ProblemID && pDone.SubIdx == pResp.SubIdx {
				hasDone = true
				break
			}
		}

		if hasDone {
			// 删除做过的题目
			for i, pR := range result {
				if pR.ProblemID == pResp.ProblemID && pR.SubIdx == pResp.SubIdx {
					result = append(result[:i], result[i+1:]...)
					break
				}
			}
		}
	}

	return c.JSON(http.StatusOK, result)
}
