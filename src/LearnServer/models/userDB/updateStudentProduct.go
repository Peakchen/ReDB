package userDB

import (
	"fmt"
	"log"

	"gopkg.in/mgo.v2/bson"
)

// UpdateStudentProductID update the product of a student
func UpdateStudentProductID(learnID int, productID string) error {
	product := struct {
		ProductID string `bson:"productID"`
	}{""}
	if err := C("products").Find(bson.M{
		"productID": productID,
	}).One(&product); err != nil || product.ProductID == "" {
		return fmt.Errorf("invalid productID")
	}

	err := C("students").Find(bson.M{
		"learnID": learnID,
		"valid":   true,
	}).One(&product)
	if err != nil {
		return err
	}

	err = C("students").Update(bson.M{
		"learnID": learnID,
		"valid":   true,
	}, bson.M{
		"$set": bson.M{
			"productID": productID,
		},
		"$push": bson.M{
			"usedProductIDs": product.ProductID,
		},
	})
	if err != nil {
		log.Printf("can not update productID of this students, err: %v\n", err)
		return err
	}

	return nil
}
