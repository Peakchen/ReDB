package staffs

import (
	"log"
	"net/http"
	"time"

	"LearnServer/models/userDB"
	
	"LearnServer/utils"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

func uploadProblemsLearnedMethodOneHandler(c echo.Context) error {
	// 上传班级题目讲解标记结果(方式1)


	uploadedData := struct {
		SchoolID string    `json:"schoolID" bson:"-"` // 学校识别码
		Grade    string    `json:"grade" bson:"-"`    // 年级
		Class    int       `json:"class" bson:"-"`    // 班级, 0 代表全部
		Time     int64     `json:"time" bson:"-"`     // 讲课时间，unix时间戳（注意因为界面只选择了日期，用0时0分0秒）
		TimeDB   time.Time `json:"-" bson:"time"`
		Problems []struct {
			ProblemHow string `json:"problemHow" bson:"problemHow"` // 出题方式（选择题、填空题、解答题）
			Source     string `json:"source" bson:"source"`         // 题目来源
		} `json:"problems" bson:"problems"`
	}{}
	if err := c.Bind(&uploadedData); err != nil {
		return utils.InvalidParams("invalid inputs, err:" + err.Error())
	}
	uploadedData.TimeDB = time.Unix(uploadedData.Time, 0)

	_, err := userDB.C("classes").Upsert(bson.M{
		"schoolID": bson.ObjectIdHex(uploadedData.SchoolID),
		"grade":    uploadedData.Grade,
		"class":    uploadedData.Class,
		"valid":    true,
	}, bson.M{
		"$push": bson.M{
			"problemsLearnedMethodOne": uploadedData,
		},
	})
	if err != nil {
		log.Printf("failed to save problems learned, err %v\n", err)
		return err
	}

	return c.JSON(http.StatusOK, "successfully uploaded")
}

func uploadProblemsLearnedMethodTwoHandler(c echo.Context) error {
	// 上传班级题目讲解标记结果(方式2)


	uploadedData := struct {
		SchoolID string    `json:"schoolID" bson:"-"` // 学校识别码
		Grade    string    `json:"grade" bson:"-"`    // 年级
		Class    int       `json:"class" bson:"-"`    // 班级, 0 代表全部
		Time     int64     `json:"time" bson:"-"`     // 讲课时间，unix时间戳（注意因为界面只选择了日期，用0时0分0秒）
		TimeDB   time.Time `json:"-" bson:"time"`
		Problems []struct {
			ProblemID string `json:"problemID" bson:"problemID"`
			SubIdx    int `json:"subIdx" bson:"subIdx"`
		} `json:"problems" bson:"problems"`
	}{}
	if err := c.Bind(&uploadedData); err != nil {
		return utils.InvalidParams("invalid inputs, err:" + err.Error())
	}
	uploadedData.TimeDB = time.Unix(uploadedData.Time, 0)

	_, err := userDB.C("classes").Upsert(bson.M{
		"schoolID": bson.ObjectIdHex(uploadedData.SchoolID),
		"grade":    uploadedData.Grade,
		"class":    uploadedData.Class,
		"valid":    true,
	}, bson.M{
		"$push": bson.M{
			"problemsLearnedMethodTwo": uploadedData,
		},
	})
	if err != nil {
		log.Printf("failed to save problems learned, err %v\n", err)
		return err
	}

	return c.JSON(http.StatusOK, "successfully uploaded")
}
