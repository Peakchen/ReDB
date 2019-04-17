package problempdfs

import (
	"log"
	"net/http"

	// "LearnServer/models/userDB"
	// "LearnServer/services/students/validation"
	// "LearnServer/utils"

	"LearnServer/models/userDB"
	"LearnServer/services/students/validation"
	"LearnServer/utils"

	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

// SetWrongProblemFileStateHandler 设置纠错本生成流程状态
func SetWrongProblemFileStateHandler(c echo.Context) error {
	var id string
	if err := validation.ValidateUser(c, &id); err != nil {
		return err
	}

	type inputType struct {
		State int `json:"state"`
	}

	input := inputType{}
	err := c.Bind(&input)
	if err != nil {
		return utils.InvalidParams("invalid input, err: " + err.Error())
	}

	err = userDB.C("students").UpdateId(bson.ObjectIdHex(id), bson.M{
		"$set": bson.M{
			"wrongProblemFileState": input.State,
		},
	})

	if err != nil {
		log.Printf("failed to update wrongProblemFileState of id %s, err: %v", id, err)
		return err
	}

	return c.JSON(http.StatusOK, "successfully set state")
}

// GetWrongProblemFileStateHandler 获取纠错本生成流程状态
func GetWrongProblemFileStateHandler(c echo.Context) error {
	var id string
	if err := validation.ValidateUser(c, &id); err != nil {
		return err
	}

	data := struct {
		WrongProblemFileState int `bson:"wrongProblemFileState" json:"state"`
	}{}

	err := userDB.C("students").FindId(bson.ObjectIdHex(id)).Select(bson.M{
		"wrongProblemFileState": 1,
	}).One(&data)

	if err != nil {
		log.Printf("failed to get wrongProblemFileState of id %s, err: %v", id, err)
		return err
	}

	return c.JSON(http.StatusOK, data)
}
