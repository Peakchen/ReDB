package students

import (
	"net/http"

	"LearnServer/models/contentDB"
	"LearnServer/models/userDB"
	"LearnServer/services/students/validation"
	"LearnServer/utils"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

func getBooksHandler(c echo.Context) error {
	var id string
	err := validation.ValidateUser(c, &id)
	if err != nil {
		return err
	}

	stu := struct {
		SchoolID bson.ObjectId `bson:"schoolID"`
		Grade    string        `bson:"grade"`
		Class    int           `bson:"classID"`
	}{}

	err = userDB.C("students").FindId(bson.ObjectIdHex(id)).One(&stu)
	if err != nil {
		return err
	}

	classBookIDs := struct {
		BookIDs []string `bson:"books"`
	}{}
	err = userDB.C("classes").Find(bson.M{
		"schoolID": stu.SchoolID,
		"grade":    stu.Grade,
		"class":    stu.Class,
		"valid":    true,
	}).One(&classBookIDs)
	if err != nil {
		if err.Error() == "not found" {
			return utils.NotFound("this student has no books")
		}
		return err
	}

	return c.JSON(http.StatusOK, contentDB.GetBooksByBookID(classBookIDs.BookIDs))
}
