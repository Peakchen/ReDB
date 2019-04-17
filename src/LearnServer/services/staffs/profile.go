package staffs

import (
	"log"
	"net/http"
	"strconv"

	"LearnServer/models/userDB"
	"LearnServer/services/staffs/validation"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

func retriveStaffProfileHandler(c echo.Context) error {
	var id string
	if err := validation.ValidateStaff(c, &id); err != nil {
		return err
	}

	// fix by stefan 20190413 9:21
	profile := struct {
		StaffID              string                 `json:"staffID" bson:"-"` // 工作人员号码
		StaffIDInt           int64                  `json:"-" bson:"staffID"` 
		Name				 string 				`json:"name" bson:"name"` //名字
		Gender				 string 				`json:"gender" bson:"gender"` //性别
		Subject			     string 				`json:"subject" bson:"subject"`//学科
    	Nickname			 string 				`json:"nickname" bson:"nickname"` //昵称
		ManageClasses        []userDB.ClassInfoType `json:"manageClasses" bson:"-"` // 管理的班级
		ManageClassObjectIDs []bson.ObjectId        `json:"-" bson:"manageClasses"`
	}{}

	err := userDB.C("staffs").FindId(bson.ObjectIdHex(id)).Select(bson.M{
		"manageClasses": 1,
		"staffID":       1,
	}).One(&profile)

	if err != nil {
		log.Printf("fail to get manageClasses of staff %v, err %v\n", id, err)
		return err
	}
	
	profile.StaffID = strconv.FormatInt(profile.StaffIDInt, 10)
	profile.ManageClasses = userDB.GetClassesByObjectID(profile.ManageClassObjectIDs)

	return c.JSON(http.StatusOK, profile)
}
