package staffs

import (
	"log"
	"net/http"
	"strconv"

	"LearnServer/models/contentDB"
	"LearnServer/models/userDB"
	"LearnServer/utils"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

func getClassProductsHandler(c echo.Context) error {
	// 获取班级配置或其分层配置信息
	schoolID := c.QueryParam("schoolID")
	grade := c.QueryParam("grade")
	classID, err := strconv.Atoi(c.QueryParam("class"))
	if err != nil {
		return utils.InvalidParams("class is invalid!")
	}
	levelStr := c.QueryParam("level")

	selectFieldName := ""
	level, err := strconv.Atoi(levelStr)
	if err == nil && level > 0 {
		selectFieldName = "level" + levelStr
	} else {
		selectFieldName = "productID"
	}

	productMap := make(map[string][]string)

	err = userDB.C("classes").Find(bson.M{
		"schoolID": bson.ObjectIdHex(schoolID),
		"grade":    grade,
		"class":    classID,
		"valid":    true,
	}).Select(bson.M{
		selectFieldName: 1,
	}).One(&productMap)
	if err != nil {
		log.Printf("fail to get product of this class, err: %v\n", err)
		return utils.NotFound("fail to get product of this class")
	}

	productID := []string{}
	productID, _ = productMap[selectFieldName]
	return c.JSON(http.StatusOK, echo.Map{
		"productID": productID,
	})
}

func updateClassProductsHandler(c echo.Context) error {
	// 更新班级配置或其分层配置信息

	type uploadType struct {
		SchoolID  string   `json:"schoolID"`
		Grade     string   `json:"grade"`
		ClassID   int      `json:"class"`
		Level     int      `json:"level"`
		ProductID []string `json:"productID"`
	}

	uploadData := uploadType{}
	if err := c.Bind(&uploadData); err != nil {
		return utils.InvalidParams("invalid input!" + err.Error())
	}

	// 确保所有的 productID 有效
	product := struct {
		ProductID string `bson:"productID"`
	}{""}
	for _, pID := range uploadData.ProductID {
		if err := userDB.C("products").Find(bson.M{
			"productID": pID,
		}).One(&product); err != nil || pID == "" {
			return utils.InvalidParams("invalid productID")
		}
	}

	updateFieldName := ""
	if uploadData.Level == -1 || uploadData.Level == 0 {
		// level 0 是 没提供 level 信息的情况，-1 或者 0 都是指修改整个班级整体产品信息
		updateFieldName = "productID"
	} else {
		updateFieldName = "level" + strconv.Itoa(uploadData.Level)
	}

	_, err := userDB.C("classes").Upsert(bson.M{
		"schoolID": bson.ObjectIdHex(uploadData.SchoolID),
		"grade":    uploadData.Grade,
		"class":    uploadData.ClassID,
		"valid":    true,
	}, bson.M{
		"$set": bson.M{
			updateFieldName: uploadData.ProductID,
		},
	})
	if err != nil {
		log.Printf("fail to update productID of class %v, err %v\n", uploadData, err)
		return err
	}

	query := bson.M{
		"schoolID": bson.ObjectIdHex(uploadData.SchoolID),
		"grade":    uploadData.Grade,
		"classID":  uploadData.ClassID,
		"valid":    true,
	}
	if uploadData.Level != -1 && uploadData.Level != 0 {
		query["level"] = uploadData.Level
	}
	studentsData := bson.M{
		"productID": uploadData.ProductID,
	}
	// if uploadData.Level != -1 && uploadData.Level != 0 {
	// 	// 清除掉学生配置的层级
	// 	studentsData["level"] = -1
	// }
	_, err = userDB.C("students").UpdateAll(query, bson.M{
		"$set": studentsData,
		"$push": bson.M{
			"usedProductIDs": uploadData.ProductID,
		},
	})
	if err != nil {
		log.Printf("fail to update productID of students, err %v \n", err)
		return err
	}

	return c.JSON(http.StatusOK, "Successfully updated productID of this class")
}

func getClassTotalLevelHandler(c echo.Context) error {
	schoolID := c.QueryParam("schoolID")
	grade := c.QueryParam("grade")
	classID, err := strconv.Atoi(c.QueryParam("class"))
	if err != nil {
		return utils.InvalidParams("class is invalid!")
	}

	result := struct {
		TotalLevel int `json:"totalLevel" bson:"totalLevel"`
	}{}

	err = userDB.C("classes").Find(bson.M{
		"schoolID": bson.ObjectIdHex(schoolID),
		"grade":    grade,
		"class":    classID,
		"valid":    true,
	}).Select(bson.M{
		"totalLevel": 1,
	}).One(&result)
	if err != nil {
		log.Printf("fail to get totalLevel of this class, err: %v\n", err)
		return utils.NotFound("can not get totalLevel of this class")
	}

	return c.JSON(http.StatusOK, result)
}

func updateClassTotalLevelHandler(c echo.Context) error {

	type uploadType struct {
		SchoolID   string `json:"schoolID"`
		Grade      string `json:"grade"`
		ClassID    int    `json:"class"`
		TotalLevel int    `json:"totalLevel"`
	}

	uploadData := uploadType{}
	if err := c.Bind(&uploadData); err != nil {
		return utils.InvalidParams("invalid input!" + err.Error())
	}

	// 获取现有的 totalLevel 方便修改新增层次的默认目标规划
	currTotalLevel := struct {
		TotalLevel int `bson:"totalLevel"`
	}{}
	err := userDB.C("classes").Find(bson.M{
		"schoolID": bson.ObjectIdHex(uploadData.SchoolID),
		"grade":    uploadData.Grade,
		"class":    uploadData.ClassID,
		"valid":    true,
	}).Select(bson.M{
		"totalLevel": 1,
	}).One(&currTotalLevel)
	if err != nil {
		log.Printf("fail to get totalLevel of this class, err: %v\n", err)
		return utils.NotFound("can not get totalLevel of this class")
	}

	_, err = userDB.C("classes").Upsert(bson.M{
		"schoolID": bson.ObjectIdHex(uploadData.SchoolID),
		"grade":    uploadData.Grade,
		"class":    uploadData.ClassID,
		"valid":    true,
	}, bson.M{
		"$set": bson.M{
			"totalLevel": uploadData.TotalLevel,
		},
	})
	if err != nil {
		log.Printf("fail to update totalLevel of class %v, err %v\n", uploadData, err)
		return err
	}

	// 设置默认目标为所有的题型
	var allTargets targetSlice
	err = contentDB.GetDB().Select(&allTargets, "SELECT chapNum, sectNum, name FROM typenames")
	if err != nil {
		log.Printf("failed to get allTargets, err %v\n", err)
		return err
	}
	exams := [4]string{
		"单元考试",
		"期中考试",
		"期末考试",
		"中考",
	}
	if currTotalLevel.TotalLevel < uploadData.TotalLevel {
		for level := currTotalLevel.TotalLevel + 1; level <= uploadData.TotalLevel; level++ {
			for _, exam := range exams {
				_, err = userDB.C("classes").Upsert(bson.M{
					"schoolID": bson.ObjectIdHex(uploadData.SchoolID),
					"grade":    uploadData.Grade,
					"class":    uploadData.ClassID,
					"valid":    true,
				}, bson.M{
					"$set": bson.M{
						"level" + strconv.Itoa(level) + "targets" + exam: allTargets,
					},
				})
				if err != nil {
					log.Printf("failed to add targets, err %v\n", err)
					continue
				}
			}
		}
	}

	// 清除减少的层级的产品和目标信息
	if currTotalLevel.TotalLevel > uploadData.TotalLevel {
		for level := uploadData.TotalLevel + 1; level <= currTotalLevel.TotalLevel; level++ {
			_, err = userDB.C("classes").Upsert(bson.M{
				"schoolID": bson.ObjectIdHex(uploadData.SchoolID),
				"grade":    uploadData.Grade,
				"class":    uploadData.ClassID,
				"valid":    true,
			}, bson.M{
				"$unset": bson.M{
					"level" + strconv.Itoa(level) + "targets": 1,
					"level" + strconv.Itoa(level):             1,
				},
			})
			if err != nil {
				log.Printf("failed to add targets, err %v\n", err)
				continue
			}
		}
	}

	return c.JSON(http.StatusOK, "successfully update total level of this class")
}
