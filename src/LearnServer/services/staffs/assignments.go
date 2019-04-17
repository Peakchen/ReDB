package staffs

import (
	"log"
	"net/http"
	"strconv"
	"time"

	// "LearnServer/models/contentDB"
	// "LearnServer/models/userDB"
	// "LearnServer/utils"

	"LearnServer/models/contentDB"
	"LearnServer/models/userDB"
	"LearnServer/utils"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

func getUnassignedBookProblemsHandler(c echo.Context) error {
	// 获取某一本书某一页的还没布置的题目

	bookID := c.QueryParam("bookID")
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		return utils.InvalidParams("filter page is invalid")
	}

	schoolID := c.QueryParam("schoolID")
	grade := c.QueryParam("grade")
	classID, err := strconv.Atoi(c.QueryParam("class"))
	if err != nil {
		return utils.InvalidParams("class is not a number!")
	}

	problems, err := contentDB.GetNonExampleProblemsByBookAndPage(bookID, page)
	if err != nil {
		return err
	}
	if len(problems) == 0 {
		return utils.NotFound("no problems in this page in the book")
	}

	type problemType struct {
		ProblemID string `json:"problemID" bson:"problemID"`
		SubIdx    int    `json:"subIdx" bson:"subIdx"`
	}

	allProblemsAssigned := struct {
		Problems []problemType `bson:"assignments"`
	}{}

	err = userDB.C("classes").Find(bson.M{
		"schoolID": bson.ObjectIdHex(schoolID),
		"grade":    grade,
		"class":    classID,
		"valid":    true,
	}).Select(bson.M{
		"assignments": 1,
	}).One(&allProblemsAssigned)
	if err != nil {
		return err
	}

	result := []contentDB.DetailedProblem{}
	// 去除已经做过的题目
	for _, p := range problems {
		hasDone := false
		for _, pDone := range allProblemsAssigned.Problems {
			if pDone.ProblemID == p.ProblemID && pDone.SubIdx == p.SubIdx {
				hasDone = true
				break
			}
		}

		if !hasDone {
			result = append(result, p)
		}
	}

	return c.JSON(http.StatusOK, result)
}

func uploadAssignmentsHandler(c echo.Context) error {
	type assignmentType struct {
		Time      time.Time `json:"time" bson:"time"`
		ProblemID string    `json:"problemID" bson:"problemID"`
		SubIdx    int       `json:"subIdx" bson:"subIdx"`
	}

	type uploadType struct {
		SchoolID string           `json:"schoolID"`
		Grade    string           `json:"grade"`
		ClassID  int              `json:"class"`
		Time     int64            `json:"time"`
		Problems []assignmentType `json:"problems" bson:"problems"`
	}

	var uploadData uploadType
	if err := c.Bind(&uploadData); err != nil {
		return utils.InvalidParams("invalid inputs, err:" + err.Error())
	}

	assignments := uploadData.Problems
	time := time.Unix(uploadData.Time, 0)
	for i := range assignments {
		assignments[i].Time = time
	}

	updateFunc := userDB.C("classes").Upsert
	selector := bson.M{
		"schoolID": bson.ObjectIdHex(uploadData.SchoolID),
		"grade":    uploadData.Grade,
		"class":    uploadData.ClassID,
		"valid":    true,
	}
	if uploadData.ClassID == 0 {
		// 更新全部班级
		updateFunc = userDB.C("classes").UpdateAll
		selector = bson.M{
			"schoolID": bson.ObjectIdHex(uploadData.SchoolID),
			"grade":    uploadData.Grade,
			"valid":    true,
		}
	}

	_, err := updateFunc(selector, bson.M{
		"$addToSet": bson.M{
			"assignments": bson.M{
				"$each": assignments,
			},
		},
	})
	if err != nil {
		log.Printf("failed to add assignments, err %v\n", err)
		return err
	}
	return c.JSON(http.StatusOK, "successfully uploaded")
}
