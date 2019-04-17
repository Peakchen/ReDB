package userDB

import (
	"log"

	"gopkg.in/mgo.v2/bson"
)

// ClassInfoType 班级信息类型
type ClassInfoType struct {
	Province	   string  		 `json:"province" bson:"province"` // 省
	City	 	   string  		 `json:"city" bson:"city"`// 市
	District	   string  		 `json:"district" bson:"district"`// 区
	SchoolName     string        `json:"schoolName" bson:"name"` // 学校名称
	SchoolObjectID bson.ObjectId `json:"-" bson:"schoolID"`
	SchoolID       string        `json:"schoolID" bson:"-"`  // 学校识别码
	Grade          string        `json:"grade" bson:"grade"` // 年级 （一、二、三、四...）
	Class          int           `json:"class" bson:"class"` // 班级
	Role 		   string 		 `json:"role" bson:"role"` //角色 
	ClassIndex	   int  		 `json:"index" bson:"classIndex"`// 标识教学信息的classIndex
}

// GetClassesByObjectID 根据 ObjectID 获取班级信息
func GetClassesByObjectID(objectIDs []bson.ObjectId) []ClassInfoType {
	classes := []ClassInfoType{}
	for _, classObjectID := range objectIDs {
		var class ClassInfoType
		err := C("classes").FindId(classObjectID).Select(bson.M{
			"schoolID": 1,
			"grade":    1,
			"class":    1,
			"role" :    1,
			"classIndex": 1,
		}).One(&class)

		if err != nil {
			log.Printf("failed to get class info of class objectID %v, error %v\n", classObjectID, err)
			continue
		}

		var school ClassInfoType
		err = C("schools").FindId(class.SchoolObjectID).Select(bson.M{
			"name": 1,
			"province": 1,
			"city": 1,
			"district": 1,

		}).One(&school)

		if err != nil {
			log.Printf("failed to get school name of school objectID %v, error %v\n", class.SchoolID, err)
			continue
		}

		class.SchoolID = class.SchoolObjectID.Hex()
		class.SchoolName = school.SchoolName
		class.Province = school.Province
		class.City = school.City
		class.District = school.District
		classes = append(classes, class)
	}
	return classes
}
