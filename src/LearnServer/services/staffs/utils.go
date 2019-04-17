package staffs

import (
	"LearnServer/models/userDB"
	"gopkg.in/mgo.v2/bson"
)

func getStudentIDByLearnID(learnID int) (string, error) {
	id := struct {
		ID bson.ObjectId `bson:"_id"`
	}{}
	err := userDB.C("students").Find(bson.M{
		"learnID": learnID,
		"valid":   true,
	}).Select(bson.M{
		"_id": 1,
	}).One(&id)
	if err != nil {
		return "", err
	}
	return id.ID.Hex(), nil
}
