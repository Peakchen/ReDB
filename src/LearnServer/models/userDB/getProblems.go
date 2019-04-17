package userDB

import "gopkg.in/mgo.v2/bson"

// GetAllProblemsDone 获取所有做过的题目，data: 存放获取到的数据的指针
func GetAllProblemsDone(id string, data interface{}) error {
	err := C("students").FindId(bson.ObjectIdHex(id)).Select(bson.M{
		"problems": 1,
	}).One(data)
	return err
}
