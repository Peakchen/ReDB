package userDB

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// StudentType 学生类型
type StudentType struct {
	ID           bson.ObjectId `bson:"_id" json:"-"`
	LearnID      int64         `bson:"learnID" json:"learnID"`
	Name         string        `bson:"realName" json:"name"`
	Time         int64         `json:"createTime"`
	CreateTimeDB time.Time     `bson:"createTime" json:"-"`
	Gender       string        `bson:"gender" json:"gender"`
	Grade        string        `bson:"grade" json:"grade"`
	Class        int           `bson:"classID" json:"class"`
	Level        int           `bson:"level" json:"level"`
}

// GetStudents 获取学生
func GetStudents(schoolID string, grade string, classID int, level int, epuStr string, serviceType string, studentName string, productID string) ([]StudentType, error) {
	// level 学生层级，-1或者0都代表全部
	// classID 班级，0代表全部
	// epuStr、serviceType、studentName、productID 为 "" 意味着不对该参数做限制
	// productID 优先级高于 epuStr、serviceType，即设定了查询含有某 productID 的学生，则 epuStr、serviceType 无效，以避免 productID 与 epuStr、serviceType 冲突问题

	// 获取符合条件的productID
	var productIDs []string
	if productID == "" {
		productQuery := bson.M{}
		if epuStr != "" {
			epuNum, err := strconv.Atoi(epuStr)
			if err != nil {
				return nil, fmt.Errorf("epu is invalid")
			}
			productQuery["epu"] = epuNum
		}
		if serviceType != "" {
			productQuery["serviceType"] = serviceType
		}

		products := []struct {
			ProductID string `bson:"productID"`
		}{}
		err := C("products").Find(productQuery).Select(bson.M{"productID": 1}).All(&products)
		if err != nil {
			log.Printf("fail to get products, query: %v, err %v\n", productQuery, err)
		}

		productIDs = make([]string, len(products))
		for i, p := range products {
			productIDs[i] = p.ProductID
		}
	} else {
		productIDs = []string{productID}
	}

	findOption := bson.M{
		"schoolID": bson.ObjectIdHex(schoolID),
		"grade":    grade,
		"productID": bson.M{
			"$in": productIDs,
		},
		"valid": true,
	}
	if classID != 0 {
		findOption["classID"] = classID
	}
	if studentName != "" {
		findOption["realName"] = bson.M{
			"$regex": studentName,
		}
	}
	if level != -1 && level != 0 {
		findOption["level"] = level
	}

	students := []StudentType{}
	err := C("students").Find(findOption).All(&students)
	if err != nil {
		return nil, fmt.Errorf("can't find students of this school and class")
	}

	for i, stu := range students {
		students[i].Time = stu.CreateTimeDB.Unix()
	}

	return students, nil
}
