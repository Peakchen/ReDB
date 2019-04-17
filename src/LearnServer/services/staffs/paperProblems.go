package staffs

import (
	"log"
	"net/http"
	"strconv"

	"LearnServer/models/contentDB"
	"LearnServer/models/userDB"
	"LearnServer/utils"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

func getPaperProblemsHandler(c echo.Context) error {
	// 获取某个试卷的未做过的题目

	learnID, err := strconv.Atoi(c.Param("learnID"))
	if err != nil {
		return utils.InvalidParams("learnID is invalid")
	}
	studentID, err := getStudentIDByLearnID(learnID)
	if err != nil {
		return utils.NotFound("can not find the information of this learnID")
	}

	type DetailedProblem struct {
		LessonName string `json:"lessonName"`
		Column     string `json:"column"`
		ProblemID  string `json:"problemID" db:"problemID"`
		Idx        int    `json:"idx,omitempty" db:"problemIndex"`
		SubIdx     int    `json:"subIdx" db:"subIdx"`
	}
	paperID := c.QueryParam("paperID")

	db := contentDB.GetDB()
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.Preparex("SELECT m.problemID as problemID, m.problemIndex, t.subIdx FROM examproblem as m, probtypes as t WHERE m.examPaperID = ? and t.problemID = m.problemID;")
	if err != nil {
		return err
	}
	rows, err := stmt.Queryx(paperID)
	if err != nil {
		return err
	}

	problems := []DetailedProblem{}
	for rows.Next() {
		var p DetailedProblem
		err := rows.StructScan(&p)
		if err != nil {
			log.Println(err)
		} else {
			problems = append(problems, p)
		}
	}
	if len(problems) == 0 {
		return utils.NotFound()
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

	var result = make([]DetailedProblem, len(problems))
	copy(result, problems[:])
	// 复制，避免在迭代过程中删除元素
	// 不能使用var result = problems[:]，不是深复制

	// 去除已经做过的题目
	for _, pResp := range problems {
		hasDone := false
		for _, pDone := range allProblemsDone.Problems {
			if pDone.SourceType == 2 && pDone.SourceID == paperID && pDone.ProblemID == pResp.ProblemID && pDone.SubIdx == pResp.SubIdx {
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
