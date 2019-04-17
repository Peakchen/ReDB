package userDB

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// GetNewID user this atomic operation to make sure no documents have duplicated id.
func GetNewID(collectionName string) (int64, error) {
	doc := struct {
		ID int64 `bson:"id"`
	}{}
	_, err := C("ids").Find(bson.M{
		"collectionName": collectionName,
	}).Apply(mgo.Change{
		Update:    bson.M{"$inc": bson.M{"id": 1}},
		ReturnNew: true,
	}, &doc)
	if err != nil {
		return -1, err
	}
	return doc.ID, nil
}
