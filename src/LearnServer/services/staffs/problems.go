package staffs

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"LearnServer/models/userDB"
	"LearnServer/utils"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

func uploadProblemResultHandler(c echo.Context) error {
	type problemResult struct {
		IsCorrect  bool   `json:"isCorrect"`
		ProblemID  string `json:"problemID"`
		SubIdx     int    `json:"subIdx"`
		SourceID   string `json:"sourceID"`
		SourceType int    `json:"sourceType"`
	}

	type uploadType struct {
		Time     int64           `json:"time"`
		Problems []problemResult `json:"problems"`
	}

	learnID, err := strconv.Atoi(c.Param("learnID"))
	if err != nil {
		return utils.InvalidParams("learnID is invalid")
	}
	studentID, err := getStudentIDByLearnID(learnID)
	if err != nil {
		return utils.NotFound("can not find the information of this learnID")
	}

	var uploadData uploadType
	if err := c.Bind(&uploadData); err != nil {
		return utils.InvalidParams("invalid inputs, err:" + err.Error())
	}

	for _, p := range uploadData.Problems {
		err := userDB.C("students").UpdateId(bson.ObjectIdHex(studentID), bson.M{
			"$push": bson.M{
				"problems": bson.M{
					"assignDate": time.Unix(uploadData.Time, 0),
					"problemID":  p.ProblemID,
					"subIdx":     p.SubIdx,
					"correct":    p.IsCorrect,
					"sourceID":   p.SourceID,
					"sourceType": p.SourceType,
				},
			},
		})
		if err != nil {
			log.Println(err)
		}
	}
	return c.JSON(http.StatusOK, "successfully uploaded")
}
