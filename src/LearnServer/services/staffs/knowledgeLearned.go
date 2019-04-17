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

func uploadKnowledgeLearnedHandler(c echo.Context) error {
	// 上传班级知识讲解标记结果


	uploadedData := struct {
		SchoolID   string    `json:"schoolID" bson:"-"` // 学校识别码
		Grade      string    `json:"grade" bson:"-"`    // 年级
		Class      int       `json:"class" bson:"-"`    // 班级, 0 代表全部
		Time       int64     `json:"time" bson:"-"`     // 讲课时间，unix时间戳（注意因为界面只选择了日期，用0时0分0秒）
		TimeDB     time.Time `json:"-" bson:"time"`
		Knowledges []struct {
			Chapter      int   `json:"chapter" bson:"chapter"`           // 章
			Section      int   `json:"section" bson:"section"`           // 节
			KnowledgeNum []int `json:"knowledgeNum" bson:"knowledgeNum"` // 知识点序号构成的数组
		} `json:"knowledges" bson:"knowledges"`
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
			"knowledgeLearned": uploadedData,
		},
	})
	if err != nil {
		log.Printf("failed to save knowledges learned, err %v\n", err)
		return err
	}

	return c.JSON(http.StatusOK, "successfully uploaded")
}
