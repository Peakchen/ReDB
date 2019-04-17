package students

import (
	"log"
	"net/http"

	"LearnServer/models/userDB"
	"LearnServer/services/students/validation"
	"LearnServer/utils"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

// setLearningPackageHandler 设置用户学习包
func setLearningPackageHandler(c echo.Context) error {
	var id string
	if err := validation.ValidateUser(c, &id); err != nil {
		return err
	}

	type inputType struct {
		Package int `json:"package"`
	}

	input := inputType{}
	err := c.Bind(&input)
	if err != nil {
		return utils.InvalidParams("invalid input, err: " + err.Error())
	}

	err = userDB.C("students").UpdateId(bson.ObjectIdHex(id), bson.M{
		"$set": bson.M{
			"learningPackage": input.Package,
		},
	})

	if err != nil {
		log.Printf("failed to update learningPackage of id %s, err: %v", id, err)
		return err
	}

	return c.JSON(http.StatusOK, "successfully set learning package")
}

// getLearningPackageHandler 获取用户学习包
func getLearningPackageHandler(c echo.Context) error {
	var id string
	if err := validation.ValidateUser(c, &id); err != nil {
		return err
	}

	data := struct {
		LearningPackage int `bson:"learningPackage" json:"package"`
	}{}

	err := userDB.C("students").FindId(bson.ObjectIdHex(id)).Select(bson.M{
		"learningPackage": 1,
	}).One(&data)

	if err != nil {
		log.Printf("failed to get learningPackage of id %s, err: %v", id, err)
		return err
	}

	return c.JSON(http.StatusOK, data)
}
