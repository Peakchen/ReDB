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

func uploadSemesterHandler(c echo.Context) error {
	// 提交班级学期起始标记


	uploadedData := struct {
		SchoolID  string `json:"schoolID"`  // 学校识别码
		Grade     string `json:"grade"`     // 年级
		Class     int    `json:"class"`     // 班级, 0 代表全部
		Semester  string `json:"semester"`  // 学期，值为 "上" "下" "未定"
		StartTime int64  `json:"startTime"` // 学期开始时间，unix时间戳
		EndTime   int64  `json:"endTime"`   // 学期结束时间，unix时间戳
	}{}
	if err := c.Bind(&uploadedData); err != nil {
		return utils.InvalidParams("invalid inputs, err:" + err.Error())
	}

	_, err := userDB.C("classes").Upsert(bson.M{
		"schoolID": bson.ObjectIdHex(uploadedData.SchoolID),
		"grade":    uploadedData.Grade,
		"class":    uploadedData.Class,
		"valid":    true,
	}, bson.M{
		"$set": bson.M{
			"semester":  uploadedData.Semester,
			"startTime": time.Unix(uploadedData.StartTime, 0),
			"endTime":   time.Unix(uploadedData.EndTime, 0),
		},
	})
	if err != nil {
		log.Printf("failed to save semester, err %v\n", err)
		return err
	}

	return c.JSON(http.StatusOK, "successfully uploaded")
}

func getSemesterHandler(c echo.Context) error {
	// 获取班级学期起始标记结果


	schoolID := c.QueryParam("schoolID")
	grade := c.QueryParam("grade")
	class, err := strconv.Atoi(c.QueryParam("class"))
	if err != nil {
		return utils.InvalidParams("class is not a number!")
	}

	result := struct {
		Semester    string    `json:"semester" bson:"semester"` // 学期，值为 "上" "下" "未定"
		StartTime   int64     `json:"startTime" bson:"-"`       // 学期开始时间，unix时间戳
		EndTime     int64     `json:"endTime" bson:"-"`         // 学期结束时间，unix时间戳
		StartTimeDB time.Time `json:"-" bson:"startTime"`       // 学期开始时间，unix时间戳
		EndTimeDB   time.Time `json:"-" bson:"endTime"`         // 学期结束时间，unix时间戳
	}{}

	err = userDB.C("classes").Find(bson.M{
		"schoolID": bson.ObjectIdHex(schoolID),
		"grade":    grade,
		"class":    class,
		"valid":    true,
	}).One(&result)
	if err != nil {
		log.Printf("failed to get semester, err %v\n", err)
		return err
	}
	result.StartTime = result.StartTimeDB.Unix()
	result.EndTime = result.EndTimeDB.Unix()

	return c.JSON(http.StatusOK, result)
}
