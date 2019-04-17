package userDB

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// AddNewStudent 添加一个学生
func AddNewStudent(schoolID string, grade string, class int, name string, gender string) error {
	schoolDetail := struct {
		Name string `bson:"name"`
	}{}
	err := C("schools").Find(bson.M{
		"_id": bson.ObjectIdHex(schoolID),
	}).One(&schoolDetail)
	if err != nil || schoolDetail.Name == "" {
		return err
	}

	if err := CreateClassIfNotExist(schoolID, grade, class); err != nil {
		return err
	}

	defaultProductIDs := []struct {
		ProductID string `bson:"productID"`
	}{}
	err = C("products").Find(bson.M{
		"default": true,
	}).Select(bson.M{
		"productID": 1,
	}).All(&defaultProductIDs)
	if err != nil {
		log.Printf("no default product, err %v\n", err)
		return err
	}
	defaultProductStrList := make([]string, len(defaultProductIDs))
	for i, p := range defaultProductIDs {
		defaultProductStrList[i] = p.ProductID
	}

	newLearnID, err := GetNewID("students")
	if err != nil {
		log.Printf("fail to get new learnID err %v", err)
		return err
	}
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	password := fmt.Sprintf("%06v", rnd.Int31n(1000000))
	err = C("students").Insert(bson.M{
		"nickName":   "",
		"realName":   name,
		"password":   password,
		"gender":     gender,
		"grade":      grade,
		"school":     schoolDetail.Name,
		"schoolID":   bson.ObjectIdHex(schoolID),
		"telephone":  "",
		"learnID":    newLearnID,
		"classID":    class,
		"level":      -1,
		"problems":   []int{},
		"tasks":      []int{},
		"createTime": time.Now(),
		"productID":  defaultProductStrList,
		"valid":      true,
	})
	if err != nil {
		log.Printf("add student failed, err %v", err)
		return err
	}

	return nil
}
