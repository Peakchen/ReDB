package userDB

import (
	"log"

	"gopkg.in/mgo.v2/bson"
)

// GetClassBookIDs 获取某个班级的书本 BookID
func GetClassBookIDs(schoolIDBson bson.ObjectId, grade string, class int) ([]string, error) {
	books := struct {
		Books []string `bson:"books"`
	}{}
	err := C("classes").Find(bson.M{
		"schoolID": schoolIDBson,
		"grade":    grade,
		"class":    class,
		"valid":    true,
	}).Select(bson.M{
		"books": 1,
	}).One(&books)
	if err != nil {
		log.Printf("getting class books failed")
		return nil, err
	}
	return books.Books, nil
}

// GetStudentBookIDs 获取某个学生的书本 BookID, studentID 学生在数据库中的_id
func GetStudentBookIDs(studentID string) ([]string, error) {
	class := struct {
		SchoolID bson.ObjectId `bson:"schoolID"`
		Grade    string        `bson:"grade"`
		Class    int           `bson:"classID"`
	}{}
	err := C("students").FindId(bson.ObjectIdHex(studentID)).Select(bson.M{
		"schoolID": 1,
		"grade":    1,
		"classID":  1,
	}).One(&class)
	if err != nil {
		log.Printf("getting student's class failed")
		return nil, err
	}
	return GetClassBookIDs(class.SchoolID, class.Grade, class.Class)
}
