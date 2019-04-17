package students

import (
	"log"
	"net/http"

	"LearnServer/models/contentDB"
	"LearnServer/models/userDB"
	"LearnServer/services/students/validation"
	"LearnServer/utils"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

func getProblemsByPaperIDHandler(c echo.Context) error {
	// 获取某一本书某一页的未做过的题目
	var id string
	if err := validation.ValidateUser(c, &id); err != nil {
		return err
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

	allProblemsDone := struct {
		Problems []problem `bson:"problems"`
	}{}

	err = userDB.C("students").FindId(bson.ObjectIdHex(id)).Select(bson.M{
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
