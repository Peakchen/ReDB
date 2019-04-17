package students

import (
	"log"
	"net/http"
	"time"

	"LearnServer/models/userDB"
	"LearnServer/services/students/validation"
	"LearnServer/utils"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

func getProductsHandler(c echo.Context) error {
	// 获取学生所有在运行的产品
	var id string
	err := validation.ValidateUser(c, &id)
	if err != nil {
		return err
	}

	stu := struct {
		ProductID []string `bson:"productID"`
	}{}

	err = userDB.C("students").FindId(bson.ObjectIdHex(id)).One(&stu)
	if err != nil {
		return err
	}

	type productType struct {
		ProductID           string    `json:"productID" bson:"productID"` // 产品编号
		DateUnix            int64     `json:"date" bson:"-"`              // 设计日期
		Date                time.Time `json:"-" bson:"date"`
		Status              bool      `json:"status" bson:"status"`                         // 服务状态
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

	productsResult := []productType{}

	for _, productID := range stu.ProductID {
		product := productType{}
		err := userDB.C("products").Find(bson.M{
			"productID": productID,
		}).One(&product)
		if err != nil {
			log.Printf("can not find this product, err: %v\n", err)
			return utils.NotFound("can not find this product")
		}

		if product.Status {
			// 不获取停用的产品
			product.DateUnix = product.Date.Unix()
			product.ServiceStartTimeInt = product.ServiceStartTime.Unix()
			product.ServiceEndTimeInt = product.ServiceEndTime.Unix()
			productsResult = append(productsResult, product)
		}
	}

	return c.JSON(http.StatusOK, productsResult)
}
