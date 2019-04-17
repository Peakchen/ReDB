package students

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"LearnServer/models/contentDB"
	"LearnServer/models/userDB"
	"LearnServer/services/students/validation"
	"LearnServer/utils"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

func uploadProblemResultHandler(c echo.Context) error {
	type problemResult struct {
		IsCorrect bool   `json:"isCorrect"`
		ProblemID string `json:"problemID"`
		SubIdx    int    `json:"subIdx"`
	}

	type uploadType struct {
		Time     int64           `json:"time"`
		Type     int             `json:"type"`
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
					"type":       uploadData.Type,
				},
			},
		})
		if err != nil {
			log.Println(err)
		}
	}
	return c.NoContent(http.StatusOK)
}

func getProblemsByPosHandler(c echo.Context) error {
	// 获取某一本书某一页的未做过的题目
	var id string
	if err := validation.ValidateUser(c, &id); err != nil {
		return err
	}

	book := c.QueryParam("book")
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		return utils.InvalidParams("filter page is invalid")
	}

	problems, err := contentDB.GetNonExampleProblemsByBookAndPage(book, page)
	if err != nil {
		return err
	}
	if len(problems) == 0 {
		return utils.NotFound()
	}

	allProblemsDone := struct {
		Problems []problem `bson:"problems"`
	}{}

	err = userDB.C("students").FindId(bson.ObjectIdHex(id)).Select(bson.M{
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
			if pDone.ProblemID == pResp.ProblemID && pDone.SubIdx == pResp.SubIdx {
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
