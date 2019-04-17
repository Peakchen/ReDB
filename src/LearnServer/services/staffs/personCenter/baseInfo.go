package personCenter
/*
	author: lazycos 2572915286@qq.com
	date: 20190412 19:50
	version: 1.0
	purpose: person base info update
*/

import (
	//"log"
	"net/http"
	//"strconv"

	"LearnServer/models/userDB"
	"LearnServer/services/staffs/validation"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
	"LearnServer/utils"
)

type TPersonBaseInfo struct {
	name  string  `json:"name" bson:"name"`//名字
    gender  string `json:"gender" bson:"gender"`//性别
    subject  string `json:"subject" bson:"subject"`//学科
    nickname string `json:"nickname" bson:"nickname"`//昵称
}

/*
	func: UpdatePersonInfo
	param1:  context
	purpose: update person base info. 
*/
func UpdatePersonInfo(c echo.Context)error{
	var id string
	err := validation.ValidateStaff(c, &id)
	if err != nil {
		return err
	}

	staffID := c.Param("staffID")
	if staffID != id {
		return utils.InvalidParams("UpdatePersonInfo invalid input staffID: " + staffID + ", right data: " + staffID)
	}

	var Baseinfo TPersonBaseInfo
	if err := c.Bind(&Baseinfo); err != nil {
		return utils.InvalidParams("UpdatePersonInfo invalid input! " + err.Error())
	}

	err = userDB.C("staffs").UpdateId(bson.ObjectIdHex(id), bson.M{
		"$set": bson.M{
			"name":    Baseinfo.name,
			"gender":    Baseinfo.gender,
			"subject" :    Baseinfo.subject,
			"nickname":    Baseinfo.nickname,
		},
	})

	return c.JSON(http.StatusOK, "Successfully update base info.")
}