package personCenter
/*
	author: lazycos 2572915286@qq.com
	date: 20190410 9:19
	version: 1.0
	purpose: use can update password or telphone bind.
*/

import (
	"log"
	"net/http"
	"strconv"

	"LearnServer/models/userDB"
	"LearnServer/services/staffs/validation"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
	"LearnServer/utils"
)

type ManageClassesInfo struct{
	province string `json:"province" bson:"province"` // 省
	city string `json:"city" bson:"city"`  // 市
	district string `json:"district" bson:"district"`  // 区

	schoolID bson.ObjectId `json:"schoolID" bson:"schoolID"`   // 学校识别码
	grade string `json:"grade" bson:"grade"` // 年级
	class int `json:"class" bson:"class"` // 班级
	role string `json:"role" bson:"role"` //角色
}

type TSchoolInfo struct {
	Province string `json:"province" bson:"province"`
	City     string `json:"city" bson:"city"`
	District string `json:"district" bson:"district"`
}

type TClassInfo struct {
	SchoolID bson.ObjectId `json:"schoolID" bson:"schoolID"`   // 学校识别码
	Grade string `json:"grade" bson:"grade"` // 年级
	Class int `json:"class" bson:"class"` // 班级
	Role string `json:"role" bson:"role"` //角色
	ClassIndex int `json:"classIndex" bson:"classIndex"` //角色
}

/*
	func: UpdateClassInfo
	param1:  context
	purpose: update class info data by class index. 
*/
func UpdateClassInfo(c echo.Context)error{
	var id string
	err := validation.ValidateStaff(c, &id)
	if err != nil {
		return err
	}

	classIndex, err := strconv.Atoi(c.Param("classIndex"))
	if err != nil {
		return utils.InvalidParams("classIndex is invalid.")
	}

	var ClassesInfo ManageClassesInfo
	if err = c.Bind(&ClassesInfo); err != nil {
		return utils.InvalidParams("AddClassInfo invalid input! " + err.Error())
	}

	StaffClassInfo := struct {
		ManageClasses        userDB.ClassInfoType `json:"manageClasses" bson:"-"` // 管理的班级
		ManageClassObjectID  bson.ObjectId        `json:"-" bson:"manageClasses"`
	}{}

	err = userDB.C("staffs").FindId(bson.ObjectIdHex(id)).Select(bson.M{
		"manageClasses": 1,
	}).One(&StaffClassInfo)

	var class userDB.ClassInfoType
	err = userDB.C("classes").FindId(StaffClassInfo.ManageClassObjectID).Select(bson.M{
		"schoolID": 1,
		"grade":    1,
		"class":    1,
		"role" :    1,
		"classIndex": 1,
	}).One(&class)

	if err != nil {
		log.Printf("failed to get class info of class objectID %v, error %v\n", StaffClassInfo.ManageClassObjectID, err)
		return err
	}

	if class.ClassIndex != classIndex{
		return utils.Unauthorized("查询classIndex错误.")
	}

	err = userDB.C("classes").UpdateId(StaffClassInfo.ManageClassObjectID, bson.M{
		"$set": bson.M{
			"grade":    ClassesInfo.grade,
			"class":    ClassesInfo.class,
			"role" :    ClassesInfo.role,
		},
	})

	var school userDB.ClassInfoType
	err = userDB.C("schools").FindId(class.SchoolObjectID).Select(bson.M{
		"province": 1,
		"city": 1,
		"district": 1,

	}).One(&school)

	if err != nil {
		log.Printf("failed to get school name of school objectID %v, error %v\n", class.SchoolID, err)
		return err
	}

	err = userDB.C("schools").UpdateId(class.SchoolObjectID, bson.M{
		"$set": bson.M{
			"province":    	ClassesInfo.province,
			"city":    		ClassesInfo.city,
			"district" :    ClassesInfo.district,
		},
	})

	if err != nil {
		log.Printf("failed to get school name of school objectID %v, error %v\n", class.SchoolID, err)
		return err
	}
	
	return c.JSON(http.StatusOK, "Successfully update class info.")
}

/*
	func: AddClassInfo
	param1:  context
	purpose: Add class info data by class index. 
*/

func AddClassInfo(c echo.Context)error{
	var id string
	if err := validation.ValidateStaff(c, &id); err != nil {
		return err
	}

	var ClassesInfo ManageClassesInfo
	if err := c.Bind(&ClassesInfo); err != nil {
		return utils.InvalidParams("AddClassInfo invalid input! " + err.Error())
	}

	SchoolInfo := TSchoolInfo {
		Province: ClassesInfo.province,
		City: ClassesInfo.city,
		District: ClassesInfo.district,
	}

	err := userDB.C("schools").Insert(SchoolInfo)
	if err != nil {
		log.Printf("unable to insert school! err: %v", err)
		return err
	}

	objID := struct {
		ID bson.ObjectId `bson:"_id"`
	}{}

	err = userDB.C("schools").Find(bson.M{
		"province": SchoolInfo.Province,
		"city":     SchoolInfo.City,
		"district": SchoolInfo.District,
	}).One(&objID)

	class := TClassInfo{
		SchoolID: objID.ID, //objID.ID
		Grade: ClassesInfo.grade,
		Class: ClassesInfo.class,
		Role: ClassesInfo.role,
	}

	err = userDB.C("classes").Insert(class)
	if err != nil {
		log.Printf("unable to insert classes! err: %v", err)
		return err
	}

	return c.JSON(http.StatusOK, "Successfully add class info.")
}

/*
	func: DeleteClassInfo
	param1:  context
	purpose: Delete class info data by class index. 
*/

func DeleteClassInfo(c echo.Context)error{
	var id string
	err := validation.ValidateStaff(c, &id)
	if err != nil {
		return err
	}

	classIndex, err := strconv.Atoi(c.Param("classIndex"))
	if err != nil {
		return utils.InvalidParams("classIndex is invalid.")
	}

	StaffClassInfo := struct {
		ManageClasses        userDB.ClassInfoType `json:"manageClasses" bson:"-"` // 管理的班级
		ManageClassObjectID  bson.ObjectId        `json:"-" bson:"manageClasses"`
	}{}

	err = userDB.C("staffs").FindId(bson.ObjectIdHex(id)).Select(bson.M{
		"manageClasses": 1,
	}).One(&StaffClassInfo)

	if err != nil {
		log.Printf("failed to get class info of class objectID %v, error %v\n", bson.ObjectIdHex(id), err)
		return err
	}

	var class userDB.ClassInfoType
	err = userDB.C("classes").FindId(StaffClassInfo.ManageClassObjectID).Select(bson.M{
		"schoolID": 1,
		"grade":    1,
		"class":    1,
		"role" :    1,
		"classIndex": 1,
	}).One(&class)

	if err != nil {
		log.Printf("failed to get class info of class objectID %v, error %v\n", StaffClassInfo.ManageClassObjectID, err)
		return err
	}

	if class.ClassIndex != classIndex{
		return utils.Unauthorized("查询classIndex错误.")
	}

	if err := userDB.C("schools").RemoveId(class.SchoolObjectID); err != nil {
		return utils.Unauthorized("删除school 失败.")
	}


	if err := userDB.C("classes").RemoveId(StaffClassInfo.ManageClassObjectID); err != nil {
		return utils.Unauthorized("classes 失败.")
	}

	return c.JSON(http.StatusOK, "Successfully delete class info.")
}