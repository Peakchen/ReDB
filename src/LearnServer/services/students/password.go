package students

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

func changePasswordHandler(c echo.Context) error {
	var id string
	err := validation.ValidateUser(c, &id)
	if err != nil {
		return err
	}

	type passwordType struct {
		Password string `bson:"password" json:"password"`
	}
	var password passwordType
	if err := c.Bind(&password); err != nil {
		return utils.InvalidParams()
	}
	err = userDB.C("students").UpdateId(bson.ObjectIdHex(id), bson.M{
		"$set": bson.M{
			"password": password.Password,
		},
	})
	if err != nil {
		log.Println(err)
		return err
	}
	return c.NoContent(http.StatusOK)
}
