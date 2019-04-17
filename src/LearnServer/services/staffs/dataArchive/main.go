package dataArchive

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo"

	"LearnServer/models/userDB"
	"LearnServer/utils"
	"gopkg.in/mgo.v2/bson"
)

func tagAllStudentsInvalid() error {
	_, err := userDB.C("students").UpdateAll(bson.M{}, bson.M{
		"$set": bson.M{
			"valid": false,
		},
	})
	return err
}

func tagAllClassesInvalid() error {
	_, err := userDB.C("classes").UpdateAll(bson.M{}, bson.M{
		"$set": bson.M{
			"valid": false,
		},
	})
	return err
}

// gradeTransform 获取变更后的年级，如果已经到最高年级，不需要新建变更数据，返回""
func gradeTransform(oldGrade string) string {
	switch oldGrade {
	case "一":
		return "二"
	case "二":
		return "三"
	case "三":
		return "四"
	case "四":
		return "五"
	case "五":
		return "六"
	case "六":
		return "七"
	case "七":
		return "八"
	case "八":
		return "九"
	case "九":
		return "高一"
	case "高一":
		return "高二"
	case "高二":
		return "高三"
	default:
		return ""
	}
}

func updateStudentsData(needGradeTransform bool) error {
	oldStudents := []struct {
		IDDB            bson.ObjectId   `bson:"_id"`
		NickName        string          `bson:"nickName"`
		RealName        string          `bson:"realName"`
		Password        string          `bson:"password"`
		Gender          string          `bson:"gender"`
		Grade           string          `bson:"grade"`
		School          string          `bson:"school"`
		SchoolID        bson.ObjectId   `bson:"schoolID,omitempty"`
		Telephone       string          `bson:"telephone"`
		LearnID         int             `bson:"learnID"`
		ClassID         int             `bson:"classID"`
		Level           int             `bson:"level"`
		CreateTime      time.Time       `bson:"createTime"`
		ProductID       []string        `bson:"productID"`
		UsedProductIDs  []string        `bson:"usedProductIDs"`
		OldStudentsIDDB []bson.ObjectId `bson:"oldStudentsIDDB,omitempty"`
	}{}
	if err := userDB.C("students").Find(bson.M{}).All(&oldStudents); err != nil {
		return err
	}

	// 新建数据
	for _, stu := range oldStudents {
		grade := stu.Grade
		if needGradeTransform {
			grade = gradeTransform(grade)
		}
		if grade != "" {
			oldStudentsIDDB := append(stu.OldStudentsIDDB, stu.IDDB)
			if err := userDB.C("students").Insert(bson.M{
				"nickName":        stu.NickName,
				"realName":        stu.RealName,
				"password":        stu.Password,
				"gender":          stu.Gender,
				"grade":           grade,
				"school":          stu.School,
				"schoolID":        stu.SchoolID,
				"telephone":       stu.Telephone,
				"learnID":         stu.LearnID,
				"classID":         stu.ClassID,
				"level":           stu.Level,
				"createTime":      stu.CreateTime,
				"productID":       stu.ProductID,
				"usedProductIDs":  stu.UsedProductIDs,
				"valid":           true,
				"oldStudentsIDDB": oldStudentsIDDB,
			}); err != nil {
				return err
			}
		}
	}

	return nil
}

func updateClassesData(needGradeTransform bool) error {
	var oldClasses []map[string]interface{}
	if err := userDB.C("classes").Find(bson.M{}).All(&oldClasses); err != nil {
		return err
	}

	// 精准匹配需要的keys
	saveKeysPrecise := map[string]bool{
		"schoolID":         true,
		"grade":            true,
		"class":            true,
		"productID":        true,
		"totalLevel":       true,
		"valid":            true,
		"oldClassesIDDB":   true,
		"examScoreRecords": true,
		"examThoughts":     true,
	}

	for _, class := range oldClasses {
		if _, ok := class["grade"]; !ok {
			class["grade"] = ""
		}
		if needGradeTransform {
			class["grade"] = gradeTransform(class["grade"].(string))
		}
		if class["grade"].(string) != "" {
			class["valid"] = true
			if _, ok := class["oldClassesIDDB"]; !ok {
				class["oldClassesIDDB"] = []bson.ObjectId{}
			}
			class["oldClassesIDDB"] = append(class["oldClassesIDDB"].([]bson.ObjectId), class["_id"].(bson.ObjectId))
			newClass := make(map[string]interface{})
			for key, value := range class {
				if _, showSave := saveKeysPrecise[key]; showSave || strings.Contains(key, "level") {
					newClass[key] = value
				}
			}
			if err := userDB.C("classes").Insert(newClass); err != nil {
				return err
			}
		}
	}

	return nil
}

// ArchiveDataHandler 封存数据，进入下学期
func ArchiveDataHandler(c echo.Context) error {
	input := struct {
		NeedGradeTransform bool `json:"needGradeTransform"`
	}{}
	if err := c.Bind(&input); err != nil {
		return utils.InvalidParams("invalid inputs, err:" + err.Error())
	}

	if err := tagAllStudentsInvalid(); err != nil {
		log.Printf("tagAllStudentsInvalid failed, err %v", err)
		return err
	}
	if err := tagAllClassesInvalid(); err != nil {
		log.Printf("tagAllClassesInvalid failed, err %v", err)
		return err
	}
	if err := updateStudentsData(input.NeedGradeTransform); err != nil {
		log.Printf("updateStudentsData failed, err %v", err)
		return err
	}
	if err := updateClassesData(input.NeedGradeTransform); err != nil {
		log.Printf("updateClassesData failed, err %v", err)
		return err
	}

	return c.JSON(http.StatusOK, "successfully archived")
}
