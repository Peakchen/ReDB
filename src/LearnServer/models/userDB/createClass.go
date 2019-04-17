package userDB

import (
	"log"

	"gopkg.in/mgo.v2/bson"
)

// CreateClassIfNotExist 如果班级不存在则创建该班级，存在则不做任何操作
func CreateClassIfNotExist(schoolID string, grade string, class int) error {
	_, err := C("classes").Upsert(bson.M{
		"schoolID":   bson.ObjectIdHex(schoolID),
		"grade":      grade,
		"class":      class,
		"totalLevel": 5,
		"valid":      true,
	}, bson.M{
		"$setOnInsert": bson.M{
			"schoolID":   bson.ObjectIdHex(schoolID),
			"grade":      grade,
			"class":      class,
			"totalLevel": 5,
			"valid":      true,
		},
	})

	if err != nil {
		log.Printf("add class failed, err: %v\n", err)
	}

	return err
}
