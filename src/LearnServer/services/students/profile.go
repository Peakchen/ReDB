package students

import (
	"net/http"

	"LearnServer/models/userDB"
	"LearnServer/services/students/validation"
	"LearnServer/utils"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

type studentProfile struct {
	LearnID        int           `json:"learnID" bson:"learnID"`
	RealName       string        `json:"realName" bson:"realName"`
	School         string        `json:"school" bson:"school"`
	SchoolID       string        `json:"schoolID"`
	SchoolObjectID bson.ObjectId `json:"-" bson:"schoolID"`
	Grade          string        `json:"grade" bson:"grade"`
	ClassID        int           `json:"classID" bson:"classID"`
	Gender         string        `json:"gender" bson:"gender"`
	Telephone      string        `json:"telephone" bson:"telephone"`
}

func getProfileHandler(c echo.Context) error {
	var id string
	err := validation.ValidateUser(c, &id)
	if err != nil {
		return err
	}
	var profile studentProfile
	err = userDB.C("students").FindId(bson.ObjectIdHex(id)).Select(bson.M{
		"realName":  1,
		"gender":    1,
		"grade":     1,
		"school":    1,
		"schoolID":  1,
		"telephone": 1,
		"learnID":   1,
		"classID":   1,
	}).One(&profile)
	if err != nil {
		return utils.NotFound("can not find this student")
	}
	profile.SchoolID = profile.SchoolObjectID.Hex()
	return c.JSON(http.StatusOK, profile)
}

func updateProfileHandler(c echo.Context) error {
	var id string
	err := validation.ValidateUser(c, &id)
	if err != nil {
		return err
	}

	var profile studentProfile
	if err := c.Bind(&profile); err != nil {
		return utils.InvalidParams()
	}

	school := struct {
		Name string `bson:"name"`
	}{}

	err = userDB.C("schools").FindId(bson.ObjectIdHex(profile.SchoolID)).One(&school)
	if err != nil {
		return err
	}

	if err := userDB.CreateClassIfNotExist(profile.SchoolID, profile.Grade, profile.ClassID); err != nil {
		return err
	}

	err = userDB.C("students").UpdateId(bson.ObjectIdHex(id), bson.M{
		"$set": bson.M{
			"gender":    profile.Gender,
			"telephone": profile.Telephone,
			"realName":  profile.RealName,
			"school":    school.Name,
			"grade":     profile.Grade,
			"schoolID":  bson.ObjectIdHex(profile.SchoolID),
			"classID":   profile.ClassID,
		},
	})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, "Successfully updated")
}
