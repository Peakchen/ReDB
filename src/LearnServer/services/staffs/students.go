package staffs

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"LearnServer/conf"
	"LearnServer/models/userDB"
	"LearnServer/utils"
	"github.com/labstack/echo"
	"github.com/tealeg/xlsx"
	"gopkg.in/mgo.v2/bson"
)

func getStudentsHandler(c echo.Context) error {

	schoolID := c.QueryParam("schoolID")
	grade := c.QueryParam("grade")
	classID, err := strconv.Atoi(c.QueryParam("class"))
	if err != nil {
		return utils.InvalidParams("class is invalid!")
	}
	studentName := c.QueryParam("studentName")
	epu := c.QueryParam("epu")
	serviceType := c.QueryParam("serviceType")
	productID := c.QueryParam("productID")

	students, err := userDB.GetStudents(schoolID, grade, classID, 0, epu, serviceType, studentName, productID)
	if err != nil {
		if err.Error() == "epu is invalid" {
			return utils.InvalidParams("epu is invalid!")
		}
		if err.Error() == "can't find students of this school and class" {
			return utils.NotFound("can't find students of this school and class")
		}
		return err
	}

	// 获取总人数
	findOption := bson.M{
		"schoolID": bson.ObjectIdHex(schoolID),
		"grade":    grade,
		"valid":    true,
	}
	if classID != 0 {
		findOption["classID"] = classID
	}
	total, err := userDB.C("students").Find(findOption).Count()
	if err != nil {
		return utils.NotFound("can't find students of this school and class")
	}

	result := struct {
		Total    int                  `json:"total"`
		LearnIDs []userDB.StudentType `json:"learnIDs"`
	}{
		Total:    total,
		LearnIDs: students,
	}
	return c.JSON(http.StatusOK, result)
}

func uploadStudentsHandler(c echo.Context) error {
	SAVE_FILE_DIR := conf.AppConfig.StudentFilesDir
	URL_PREFIX := conf.AppConfig.StudentFilesURL

	type studentsUploadType struct {
		SchoolID string `json:"schoolID"`
		Grade    string `json:"grade"`
		Class    int    `json:"class"`
		UID      string `json:"studentFile"`
	}

	uploadData := studentsUploadType{}
	if err := c.Bind(&uploadData); err != nil {
		return utils.InvalidParams("invalid input!" + err.Error())
	}

	schoolDetail := struct {
		Name string `bson:"name"`
	}{}
	err := userDB.C("schools").Find(bson.M{
		"_id": bson.ObjectIdHex(uploadData.SchoolID),
	}).One(&schoolDetail)
	if err != nil || schoolDetail.Name == "" {
		return utils.NotFound("can not find this school")
	}

	if err := userDB.CreateClassIfNotExist(uploadData.SchoolID, uploadData.Grade, uploadData.Class); err != nil {
		return err
	}

	studentsData, err := utils.Manager.Get(uploadData.UID)

	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Sheet1")
	if err != nil {
		return err
	}
	row := sheet.AddRow()
	row.SetHeightCM(1)
	cell := row.AddCell()
	cell.Value = "姓名"
	cell = row.AddCell()
	cell.Value = "性别"
	cell = row.AddCell()
	cell.Value = "学习号"
	cell = row.AddCell()
	cell.Value = "密码"

	failedStudents := []string{}

	defaultProductIDs := []struct {
		ProductID string `bson:"productID"`
	}{}
	err = userDB.C("products").Find(bson.M{
		"default": true,
	}).Select(bson.M{
		"productID": 1,
	}).All(&defaultProductIDs)
	if err != nil {
		log.Printf("no default product, err %v\n", err)
		return err
	}
	defaultProductStrList := make([]string, len(defaultProductIDs))
	for i, p := range defaultProductIDs {
		defaultProductStrList[i] = p.ProductID
	}

	for _, student := range studentsData.([]studentType) {
		newLearnID, err := userDB.GetNewID("students")
		if err != nil {
			failedStudents = append(failedStudents, student.Name)
			log.Printf("fail to get new learnID err %v", err)
			continue
		}
		rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
		password := fmt.Sprintf("%06v", rnd.Int31n(1000000))
		err = userDB.C("students").Insert(bson.M{
			"nickName":   "",
			"realName":   student.Name,
			"password":   password,
			"gender":     student.Gender,
			"grade":      uploadData.Grade,
			"school":     schoolDetail.Name,
			"schoolID":   bson.ObjectIdHex(uploadData.SchoolID),
			"telephone":  "",
			"learnID":    newLearnID,
			"classID":    uploadData.Class,
			"level":      -1,
			"problems":   []int{},
			"tasks":      []int{},
			"createTime": time.Now(),
			"productID":  defaultProductStrList,
			"valid":      true,
		})
		if err != nil {
			failedStudents = append(failedStudents, student.Name)
			log.Printf("add student %v failed, err %v", student, err)
			continue
		}

		row := sheet.AddRow()
		row.SetHeightCM(1)
		cell := row.AddCell()
		cell.Value = student.Name
		cell = row.AddCell()
		cell.Value = student.Gender
		cell = row.AddCell()
		cell.Value = "学习号" + strconv.FormatInt(newLearnID, 10)
		cell = row.AddCell()
		cell.Value = "密码" + password
	}

	if len(failedStudents) != 0 {
		row := sheet.AddRow()
		row.SetHeightCM(1)
		cell := row.AddCell()
		cell.Value = "以下学生创建失败"
		for _, name := range failedStudents {
			row := sheet.AddRow()
			row.SetHeightCM(1)
			cell := row.AddCell()
			cell.Value = name
		}
	}

	fileName := schoolDetail.Name + uploadData.Grade + strconv.Itoa(uploadData.Class) + "班学生账户信息表.xlsx"
	err = file.Save(SAVE_FILE_DIR + fileName)
	if err != nil {
		log.Printf("failed to create file err %v \n", err)
		return err
	}

	utils.Manager.Delete(uploadData.UID)

	return c.JSON(http.StatusOK, bson.M{
		"URL": URL_PREFIX + fileName,
	})
}

func addOneNewStudentHandler(c echo.Context) error {
	// 添加一个新学生
	type uploadType struct {
		SchoolID string `json:"schoolID" bson:"schoolID"` // 学校识别码
		Grade    string `json:"grade" bson:"grade"`       // 年级
		Class    int    `json:"class" bson:"classID"`
		Name     string `json:"name" bson:"name"`
		Gender   string `json:"gender" bson:"gender"`
	}

	uploadData := uploadType{}
	if err := c.Bind(&uploadData); err != nil {
		return utils.InvalidParams("invalid input, error: " + err.Error())
	}

	err := userDB.AddNewStudent(uploadData.SchoolID, uploadData.Grade, uploadData.Class, uploadData.Name, uploadData.Gender)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "successfully added")
}

func deleteStudentHandler(c echo.Context) error {
	learnID, err := strconv.Atoi(c.Param("learnID"))
	if err != nil {
		return utils.InvalidParams("wrong learnID!")
	}

	err = userDB.C("students").Remove(bson.M{
		"learnID": learnID,
		"valid":   true,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "successfully deleted!")
}

func retriveStudentDetailHandler(c echo.Context) error {
	// 获取某个学生的个人信息

	learnID, err := strconv.Atoi(c.Param("learnID"))
	if err != nil {
		return utils.InvalidParams("wrong learnID!")
	}

	type studentType struct {
		Name      string `bson:"realName" json:"name"`
		Gender    string `bson:"gender" json:"gender"`
		School    string `bson:"school" json:"school"`
		Grade     string `bson:"grade" json:"grade"`
		Class     int    `bson:"classID" json:"class"`
		ProductID string `bson:"productID" json:"productID"`
	}

	student := studentType{}
	err = userDB.C("students").Find(bson.M{
		"learnID": learnID,
		"valid":   true,
	}).One(&student)
	if err != nil {
		log.Printf("finding student %d failed, err %v\n", learnID, err)
		return utils.NotFound("can not find this student")
	}

	return c.JSON(http.StatusOK, student)
}

func updateStudentInfoHandler(c echo.Context) error {
	// 更新某个学生

	learnID, err := strconv.Atoi(c.Param("learnID"))
	if err != nil {
		return utils.InvalidParams("wrong learnID!")
	}

	type uploadType struct {
		Name   string `json:"name" bson:"realName"`
		Gender string `json:"gender" bson:"gender"`
		Class  int    `json:"class" bson:"classID"`
	}

	uploadData := uploadType{}
	if err := c.Bind(&uploadData); err != nil {
		return utils.InvalidParams("invalid input, error: " + err.Error())
	}

	err = userDB.C("students").Update(bson.M{
		"learnID": learnID,
		"valid":   true,
	}, bson.M{
		"$set": bson.M{
			"realName": uploadData.Name,
			"gender":   uploadData.Gender,
			"classID":  uploadData.Class,
		},
	})
	if err != nil {
		log.Printf("can not update this students, err: %v\n", err)
		return err
	}

	return c.JSON(http.StatusOK, "successfully updated this student")
}

func updateStudentProductIDHandler(c echo.Context) error {
	// 更新某个学生的在用产品

	learnID, err := strconv.Atoi(c.Param("learnID"))
	if err != nil {
		return utils.InvalidParams("wrong learnID!")
	}

	type uploadType struct {
		ProductID string `json:"productID" bson:"productID"`
	}

	uploadData := uploadType{}
	if err := c.Bind(&uploadData); err != nil {
		return utils.InvalidParams("invalid input, error: " + err.Error())
	}

	err = userDB.UpdateStudentProductID(learnID, uploadData.ProductID)
	if err != nil {
		if err.Error() == "invalid productID" {
			return utils.InvalidParams("invalid productID!")
		}
		log.Printf("can not update productID of this students, err: %v\n", err)
		return err
	}

	return c.JSON(http.StatusOK, "Successfully updated productID of this student")
}

func updateStudentLevelHandler(c echo.Context) error {
	type uploadType struct {
		LearnID int `json:"learnID" bson:"learnID"`
		Level   int `json:"level" bson:"level"`
	}

	uploadData := []uploadType{}
	if err := c.Bind(&uploadData); err != nil {
		return utils.InvalidParams("invalid input, error: " + err.Error())
	}

	for _, stu := range uploadData {
		err := userDB.C("students").Update(bson.M{
			"learnID": stu.LearnID,
			"valid":   true,
		}, bson.M{
			"$set": bson.M{
				"level": stu.Level,
			},
		})
		if err != nil {
			log.Printf("fail to update level of learnID %d, err: %v\n", stu.LearnID, err)
		}
	}

	return c.JSON(http.StatusOK, "Successfully updated students' levels.")
}
