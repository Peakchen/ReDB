package staffs

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"LearnServer/models/userDB"
	
	"LearnServer/utils"
	"github.com/labstack/echo"
	"github.com/tealeg/xlsx"
	"gopkg.in/mgo.v2/bson"
)

func uploadScoreFileHandler(c echo.Context) error {
	// 上传班级成绩excel表获取数据

	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	xls, err := xlsx.OpenReaderAt(src, file.Size)
	if err != nil {
		return err
	}

	if len(xls.Sheets) > 0 && len(xls.Sheets[0].Rows) > 1 && len(xls.Sheets[0].Rows[0].Cells) > 1 &&
		xls.Sheets[0].Rows[0].Cells[0].String() != "姓名" || xls.Sheets[0].Rows[0].Cells[1].String() != "成绩" {
		return utils.Forbidden("wrong format")
	}

	type resultType struct {
		Name  string  `xlsx:"0" json:"name"`  // 姓名（暂不考虑班级内出现重名的情况）
		Score float64 `xlsx:"1" json:"score"` // 成绩
	}
	results := []resultType{}
	for _, row := range xls.Sheets[0].Rows[1:] {
		s := resultType{}
		err = row.ReadStruct(&s)
		if err != nil {
			log.Println(err)
			continue
		}
		results = append(results, s)
	}
	return c.JSON(http.StatusOK, results)
}

func uploadExamScoresHandler(c echo.Context) error {
	// 上传学生考试成绩


	type examScoreType struct {
		SchoolID string `json:"schoolID"` // 学校识别码
		Grade    string `json:"grade"`    // 年级
		Class    int    `json:"class"`    // 班级, 0 代表全部
		Time     int64  `json:"time"`     // 考试时间，unix时间戳（注意因为界面只选择了日期，用0时0分0秒）
		PaperID  string `json:"paperID"`  // 考试试卷ID
		Scores   []struct {
			LearnID int     `json:"learnID" bson:"learnID"` // 学生学习号
			Name    string  `json:"name" bson:"name"`       // 学生姓名
			Score   float32 `json:"score" bson:"score"`     // 成绩
		} `json:"scores" bson:"scores"` // 成绩
	}

	var uploadedData examScoreType
	if err := c.Bind(&uploadedData); err != nil {
		return utils.InvalidParams("invalid input!" + err.Error())
	}

	for _, data := range uploadedData.Scores {
		// examScores ：一个 Map ，每个 key 为 paperID ， value 为对应数据
		err := userDB.C("students").Update(bson.M{
			"learnID": data.LearnID,
			"valid":   true,
		}, bson.M{
			"$set": bson.M{
				"examScores.paperID" + uploadedData.PaperID: bson.M{
					"time":  time.Unix(uploadedData.Time, 0),
					"score": data.Score,
				},
			},
		})
		if err != nil {
			log.Printf("failed to save score of learnID %d, error %v\n", data.LearnID, err)
			return err
		}
	}

	_, err := userDB.C("classes").Upsert(bson.M{
		"schoolID": bson.ObjectIdHex(uploadedData.SchoolID),
		"grade":    uploadedData.Grade,
		"class":    uploadedData.Class,
		"valid":    true,
	}, bson.M{
		"$set": bson.M{
			"examScoreRecords.paperID" + uploadedData.PaperID: bson.M{
				"time":   time.Unix(uploadedData.Time, 0),
				"scores": uploadedData.Scores,
			},
		},
	})
	if err != nil {
		log.Printf("failed to save score in classes, error %v\n", err)
		return err
	}
	return c.JSON(http.StatusOK, "successfully uploaded")
}

func getExamScoresHandler(c echo.Context) error {
	// 获取学生考试成绩


	schoolID := c.QueryParam("schoolID")
	grade := c.QueryParam("grade")
	class, err := strconv.Atoi(c.QueryParam("class"))
	if err != nil {
		return utils.InvalidParams("class is not a number!")
	}
	paperID := c.QueryParam("paperID")

	type examScoreType struct {
		Time   int64     `json:"time" bson:"-"` // 考试时间，unix时间戳（注意因为界面只选择了日期，用0时0分0秒）
		TimeDB time.Time `json:"-" bson:"time"`
		Scores []struct {
			LearnID int     `json:"learnID" bson:"learnID"` // 学生学习号
			Name    string  `json:"name" bson:"name"`       // 学生姓名
			Score   float32 `json:"score" bson:"score"`     // 成绩
		} `json:"scores" bson:"scores"` // 成绩
	}

	classExamRecord := make(map[string]map[string]examScoreType)

	err = userDB.C("classes").Find(bson.M{
		"schoolID": bson.ObjectIdHex(schoolID),
		"grade":    grade,
		"class":    class,
		"valid":    true,
	}).Select(bson.M{
		"examScoreRecords": 1,
	}).One(&classExamRecord)
	if err != nil {
		log.Printf("getting class examRecord failed")
		return err
	}

	examRecord, ok := classExamRecord["examScoreRecords"]
	if !ok {
		return utils.NotFound("can not find this exam record!")
	}

	result, ok := examRecord["paperID"+paperID]
	if !ok {
		return utils.NotFound("can not find this exam record!")
	}
	result.Time = result.TimeDB.Unix()

	return c.JSON(http.StatusOK, result)
}
