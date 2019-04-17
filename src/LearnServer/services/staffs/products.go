package staffs

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"time"

	"LearnServer/models/userDB"
	"LearnServer/utils"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

// ProductDetailType 产品信息
type ProductDetailType struct {
	// 包括所有需要配置的字段，不包括 productID 与 date、status
	ProblemCode         string    `json:"problemCode" bson:"problemCode"`               // 问题代码，错题学习为"E"
	Gradation           int       `json:"gradation" bson:"gradation"`                   // 层次， 1 2 3
	Depth               int       `json:"depth" bson:"depth"`                           // 深度， 1 2 3
	Name                string    `json:"name" bson:"name"`                             // 产品名称
	Level               string    `json:"level" bson:"level"`                           // 产品级别
	Object              string    `json:"object" bson:"object"`                         // 产品对象
	Epu                 int       `json:"epu" bson:"epu"`                               // EPU, 1 2
	ProblemMax          int       `json:"problemMax" bson:"problemMax"`                 // 题量控制
	WrongProblemStatus  int       `json:"wrongProblemStatus" bson:"wrongProblemStatus"` // 错题状态，1 现在仍错，2 曾经错过
	ProblemType         []string  `json:"problemType" bson:"problemType"`               // 题目种类
	SameTypeMax         int       `json:"sameTypeMax" bson:"sameTypeMax"`               // 同类最大题量
	SameTypeSource      []string  `json:"sameTypeSource" bson:"sameTypeSource"`         // 同类来源
	ProblemTemplateID   string    `json:"problemTemplateID" bson:"problemTemplateID"`   // 题目模板
	AnswerTemplateID    string    `json:"answerTemplateID" bson:"answerTemplateID"`     // 答案模板
	BorderControl       string    `json:"borderControl" bson:"borderControl"`           // 边界控制
	ProblemSource       []string  `json:"problemSource" bson:"problemSource"`           // 错题源， 如： ["课本", "平时试卷"]
	ServiceType         string    `json:"serviceType" bson:"serviceType"`               // 服务类型
	ServiceLauncher     string    `json:"serviceLauncher" bson:"serviceLauncher"`       // 服务发起
	ServiceStartTimeInt int64     `json:"serviceStartTime" bson:"-"`                    // 服务开始时间, unix时间戳
	ServiceStartTime    time.Time `json:"-" bson:"serviceStartTime"`                    // 服务开始时间
	ServiceEndTimeInt   int64     `json:"serviceEndTime" bson:"-"`                      // 服务结束时间, unix时间戳
	ServiceEndTime      time.Time `json:"-" bson:"serviceEndTime"`                      // 服务结束时间
	ServiceTimes        int       `json:"serviceTimes" bson:"serviceTimes"`             // 服务次数
	ServiceDuration     string    `json:"serviceDuration" bson:"serviceDuration"`       // 服务时长
	DeliverType         string    `json:"deliverType" bson:"deliverType"`               // 交付类型
	DeliverPriority     int       `json:"deliverPriority" bson:"deliverPriority"`       // 交付优先级
	DeliverTime         []struct {
		Day  int    `json:"day" bson:"day"`   // 周日0,周一到周六分别是1到6,
		Time string `json:"time" bson:"time"` // 时间，格式按照"08:00:00"
	} `json:"deliverTime" bson:"deliverTime"` // 交付节点
	DeliverExpected  int    `json:"deliverExpected" bson:"deliverExpected"`   // 交付预期，预期多少小时内
	ExceptionHandler int    `json:"exceptionHandler" bson:"exceptionHandler"` // 异常处理，发现未标记：1 全部标记为对再生成， 2 全部标记为错再生成， 3 不生成
	Price            int    `json:"price" bson:"price"`                       // 单价
	Subject          string `json:"subject" bson:"subject"`                   // 学科
	Grade            string `json:"grade" bson:"grade"`                       // 年级（全部直接用“全部”）
}

func setProductAutoStartAndEndTime(productID string, startTime time.Time, endTime time.Time) error {
	// 添加自动任务，自动启动与停止产品
	startTimeStr := startTime.Format("15:04 2006-01-02")
	startCmd := exec.Command("/bin/bash", "-c", "echo '/bin/bash ./update_product_status.sh "+productID+" on' | at "+startTimeStr)
	if err := startCmd.Start(); err != nil {
		log.Printf("failed to auto start product, err %v \n", err)
		return err
	}

	endTimeStr := endTime.Format("15:04 2006-01-02")
	endCmd := exec.Command("/bin/bash", "-c", "echo '/bin/bash ./update_product_status.sh "+productID+" off' | at "+endTimeStr)
	if err := endCmd.Start(); err != nil {
		log.Printf("failed to auto start product, err %v \n", err)
		return err
	}
	return nil
}

func uploadProductHandler(c echo.Context) error {
	// 上传新的产品
	type productType struct {
		ProductDetailType `bson:",inline"`
		ProductID         string    `bson:"productID"`
		Date              time.Time `bson:"date"`
		Status            bool      `bson:"status"`
	}

	uploadData := productType{}
	if err := c.Bind(&uploadData); err != nil {
		return utils.InvalidParams("invalid input, error: " + err.Error())
	}

	newProductNumberInt64, err := userDB.GetNewID("products")
	if err != nil {
		log.Printf("failed to get new product number, err %v\n", err)
		return err
	}
	newProductNumber := fmt.Sprintf("%03d", newProductNumberInt64)

	uploadData.ProductID = uploadData.ProblemCode + strconv.Itoa(uploadData.Gradation) + strconv.Itoa(uploadData.Depth) + "-" + newProductNumber
	uploadData.Date = time.Now()
	uploadData.ServiceStartTime = time.Unix(uploadData.ServiceStartTimeInt, 0)
	uploadData.ServiceEndTime = time.Unix(uploadData.ServiceEndTimeInt, 0)
	// 在服务时段内则设置状态为启用
	uploadData.Status = uploadData.ServiceStartTime.Before(time.Now()) && time.Now().Before(uploadData.ServiceEndTime)
	if err := setProductAutoStartAndEndTime(uploadData.ProductID, uploadData.ServiceStartTime, uploadData.ServiceEndTime); err != nil {
		return err
	}

	err = userDB.C("products").Insert(uploadData)
	if err != nil {
		log.Printf("cannot save new product, err: %v\n", err)
		return err
	}

	return c.JSON(http.StatusOK, "successfully uploaded a product")
}

func retriveProductHandler(c echo.Context) error {
	// 根据ID获取一个产品信息
	type productType struct {
		ProductDetailType `bson:",inline"`
		ProductID         string    `json:"productID" bson:"productID"` // 产品编号
		DateUnix          int64     `json:"date" bson:"-"`              // 设计日期
		Date              time.Time `json:"-" bson:"date"`
		Status            bool      `json:"status" bson:"status"` // 服务状态
	}

	result := productType{}
	productID := c.Param("productID")

	err := userDB.C("products").Find(bson.M{
		"productID": productID,
	}).One(&result)
	if err != nil {
		log.Printf("can not find this product, err: %v\n", err)
		return utils.NotFound("can not find this product")
	}

	result.DateUnix = result.Date.Unix()
	result.ServiceStartTimeInt = result.ServiceStartTime.Unix()
	result.ServiceEndTimeInt = result.ServiceEndTime.Unix()

	return c.JSON(http.StatusOK, result)
}

func listProductHandler(c echo.Context) error {
	// 获取符合筛选条件的产品信息
	epu, err := strconv.Atoi(c.QueryParam("epu"))
	if err != nil {
		return utils.InvalidParams("query param epu is invalid")
	}

	object := c.QueryParam("object")

	type productType struct {
		ProductDetailType `bson:",inline"`
		ProductID         string    `json:"productID" bson:"productID"` // 产品编号
		DateUnix          int64     `json:"date" bson:"-"`              // 设计日期
		Date              time.Time `json:"-" bson:"date"`
		Status            bool      `json:"status" bson:"status"` // 服务状态
	}

	products := []productType{}

	query := bson.M{}
	if epu != -1 {
		query["epu"] = epu
	}
	if object != "all" {
		query["object"] = object
	}

	err = userDB.C("products").Find(query).All(&products)
	if err != nil {
		log.Printf("finding products failed, err: %v\n", err)
		return utils.NotFound("finding products failed")
	}

	for i := range products {
		products[i].DateUnix = products[i].Date.Unix()
		products[i].ServiceStartTimeInt = products[i].ServiceStartTime.Unix()
		products[i].ServiceEndTimeInt = products[i].ServiceEndTime.Unix()
	}

	return c.JSON(http.StatusOK, products)
}

func updateProductHandler(c echo.Context) error {
	// 修改一个产品
	type productType struct {
		ProductDetailType `bson:",inline"`
		Date              time.Time `bson:"date"`
		Status            bool      `bson:"status"`
	}

	productID := c.Param("productID")

	uploadData := productType{}
	if err := c.Bind(&uploadData); err != nil {
		return utils.InvalidParams("invalid input, error: " + err.Error())
	}

	uploadData.Date = time.Now()
	uploadData.ServiceStartTime = time.Unix(uploadData.ServiceStartTimeInt, 0)
	uploadData.ServiceEndTime = time.Unix(uploadData.ServiceEndTimeInt, 0)
	// 在服务时段内则设置状态为启用
	uploadData.Status = uploadData.ServiceStartTime.Before(time.Now()) && time.Now().Before(uploadData.ServiceEndTime)
	if err := setProductAutoStartAndEndTime(productID, uploadData.ServiceStartTime, uploadData.ServiceEndTime); err != nil {
		return err
	}

	err := userDB.C("products").Update(bson.M{
		"productID": productID,
	}, bson.M{
		"$set": uploadData,
	})
	if err != nil {
		log.Printf("can not update this product, err: %v\n", err)
		return err
	}

	return c.JSON(http.StatusOK, "Successfully updated this product")
}

func updateProductStatusHandler(c echo.Context) error {
	// 修改一个产品的状态信息
	status := struct {
		Status bool `json:"status" bson:"status"`
	}{}

	productID := c.Param("productID")
	if err := c.Bind(&status); err != nil {
		return utils.InvalidParams("invalid input, error: " + err.Error())
	}

	err := userDB.C("products").Update(bson.M{
		"productID": productID,
	}, bson.M{
		"$set": bson.M{
			"status": status.Status,
		},
	})
	if err != nil {
		log.Printf("can not update status of this product, err: %v\n", err)
		return err
	}

	return c.JSON(http.StatusOK, "Successfully updated the status of this product")
}
