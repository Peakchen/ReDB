package personCenter

/*
	author: lazycos 2572915286@qq.com
	date: 20190411 9:28
	version: 1.0
	purpose: use can update password or telphone bind.
*/

import (
	"log"
	"net/http"

	"LearnServer/models/userDB"
	"LearnServer/services/students/validation"
	"LearnServer/utils"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
	"regexp"
)

type TLearnUserInfo struct {
	LearnID     int64  `bson:"learnID" json:"learnID"`
	Password 	string  `bson:"password"  json:"password"`
}

type TPasswordData struct {
	OldPassword string `bson:"password" json:"password"`
	NewPassword string `bson:"newpassword" json:"newpassword"`
	NewPassword2 string `bson:"newpassword2" json:"newpassword2"`
}

type TTelphoneInfo struct {
	TelphoneNumber string `bson:"telphone" json:"telphone"`
}

const (
    regular = "^(13[0-9]|14[57]|15[0-35-9]|18[07-9])\\d{8}$"
)
 
/*
	func: CheckTelphoneNumberRight
	param1:  mobileNum  type: string
	purpose: check user input number is right Telphone phone. 
*/
func CheckTelphoneNumberRight(TelphoneNum string) bool {
    reg := regexp.MustCompile(regular)
    return reg.MatchString(TelphoneNum)
}

/*
	func: canUpdateUsePsw
	param1:  id  type: string
	param2:  password type: string
	purpose: check passwd is input true, then give vailable tips. 
*/
func canUpdateUsePsw(id, password string)error{
	useObj := TLearnUserInfo{}
	err := userDB.C("classes").FindId(bson.ObjectIdHex(id)).Select(bson.M{
		"learnID": 1,
		"password": 1,
	}).One(&useObj)

	if err != nil {
		return utils.Unauthorized("用户ID或者密码错误！")
	}

	if password != useObj.Password {
		return utils.Unauthorized("当前密码输入错误！")
	}

	return nil
}

/*
	func: ChangeUserPasswd
	param1:  context
	purpose: change user true password.  
*/
func ChangeUserPasswd(c echo.Context) error {
	var id string
	err := validation.ValidateUser(c, &id)
	if err != nil {
		return err
	}

	var password TPasswordData
	if err := c.Bind(&password); err != nil {
		return utils.InvalidParams()
	}

	if PswErr := canUpdateUsePsw(id, password.OldPassword); PswErr != nil{
		return PswErr
	}

	if password.NewPassword != password.NewPassword2{
		return utils.Unauthorized("新密码确认匹配失败，请再次输入验证！")
	}

	err = userDB.C("classes").UpdateId(bson.ObjectIdHex(id), bson.M{
		"$set": bson.M{
			"password": password.NewPassword,
		},
	})

	if err != nil {
		log.Println(err)
		return err
	}

	return c.JSON(http.StatusOK, "Successfully update passwd info.")
}

/*
	func: BindTelPhoneNumber
	params: context
	purpose: bind current user telphone number.
*/
func BindTelPhoneNumber(c echo.Context)error{
	var id string
	err := validation.ValidateUser(c, &id)
	if err != nil {
		return err
	}

	var telphoneInfo TTelphoneInfo
	if err := c.Bind(&telphoneInfo); err != nil {
		return utils.InvalidParams()
	}

	if !CheckTelphoneNumberRight(telphoneInfo.TelphoneNumber){
		return utils.Unauthorized("当前手机号输入错误，请检查在进行绑定！")
	}

	err = userDB.C("classes").UpdateId(bson.ObjectIdHex(id), bson.M{
		"$set": bson.M{
			"telphone": telphoneInfo.TelphoneNumber,
		},
	})

	if err != nil {
		log.Println(err)
		return err
	}

	return c.JSON(http.StatusOK, "Successfully bind.")
}