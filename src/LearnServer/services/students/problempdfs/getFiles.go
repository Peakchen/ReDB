package problempdfs

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	// "LearnServer/tools/documents"

	// "LearnServer/conf"
	// "LearnServer/models/userDB"
	// "LearnServer/services/students/validation"
	// "LearnServer/utils"

	"LearnServer/tools/documents"
	"LearnServer/conf"
	"LearnServer/models/userDB"
	"LearnServer/services/students/validation"
	"LearnServer/utils"

	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

// templateDetailType 模板信息
type templateDetailType struct {
	Name          string  `json:"name" bson:"name"`                   // 模板名称
	Info          string  `json:"info" bson:"info"`                   // 模板说明
	Type          string  `json:"type" bson:"type"`                   // 模板类型
	PageType      string  `json:"pageType" bson:"pageType"`           // 纸张大小，"A3"或者"A4"
	PageDirection string  `json:"pageDirection" bson:"pageDirection"` // 纸张方向
	ColumnCount   int     `json:"columnCount" bson:"columnCount"`     // 分栏数
	MarginTop     float64 `json:"marginTop" bson:"marginTop"`         // 上下页边距
	MarginLeft    float64 `json:"marginLeft" bson:"marginLeft"`       // 左右页边距
	Operations    []struct {
		Type       int      `json:"type" bson:"type"`             // 1 变更格式 2 添加内容
		Font       string   `json:"font" bson:"font"`             // 字体（仅当type为1变更格式时有效）
		FontSize   int      `json:"fontSize" bson:"fontSize"`     // 字号（仅当type为1变更格式时有效）
		FontEffect []string `json:"fontEffect" bson:"fontEffect"` // 字体效果 bold 加粗， underlined 下划线 italic 倾斜（仅当type为1变更格式时有效）
		Alignment  int      `json:"alignment" bson:"alignment"`   // 对齐，0 左对齐 1 居中对齐 2 右对齐（仅当type为1变更格式时有效）
		RowSpacing float64  `json:"rowSpacing" bson:"rowSpacing"` // 行距（仅当type为1变更格式时有效）
		Content    string   `json:"content" bson:"content"`       // 内容（仅当type为2添加内容时有效）
	} `json:"operations" bson:"operations"` // 文档开头部分的操作
}

func hasMarkTasks(id string) bool {
	// 判断是否有未标记的纠错本
	type taskInfoType struct {
		Time time.Time `bson:"time"`
	}

	tasksInfo := struct {
		Tasks []taskInfoType `bson:"tasks"`
	}{}
	err := userDB.C("students").FindId(bson.ObjectIdHex(id)).Select(bson.M{
		"tasks": 1,
	}).One(&tasksInfo)
	if err != nil {
		log.Printf("fail to get tasks: id: %s, err: %v", id, err)
		return false
	}

	return len(tasksInfo.Tasks) > 0
}

func sendGettingFileRequestCallBackCreator(columnName string, id string) func(documents.DocumentURLType, int) {
	return func(urls documents.DocumentURLType, statusCode int) {
		// 更新用户数据库中的lastURL
		if statusCode != 200 {
			// failed
			return
		}
		contentServer := conf.AppConfig.FilesServer
		err := userDB.C("students").UpdateId(bson.ObjectIdHex(id), bson.M{
			"$set": bson.M{
				columnName: contentServer + urls.PdfURL,
			},
		})
		if err != nil {
			log.Printf("set "+columnName+" of id %s failed, err : %v", id, err)
		}
	}
}

func setLastFileURLToBlank(columnName string, id string) error {
	// 将 lastProblemPDF lastAnswerPDF 这些 url 设置为""，代表新文件已经加入生成队列，但尚未生成完成
	err := userDB.C("students").UpdateId(bson.ObjectIdHex(id), bson.M{
		"$set": bson.M{
			columnName: "",
		},
	})
	if err != nil {
		log.Printf("set "+columnName+" of id %s failed, err : %v", id, err)
	}
	return err
}

func getTemplate(productID string, templateType string) (templateDetailType, error) {
	// 获取用户产品中设定的模板
	template := templateDetailType{}

	templateID := struct {
		ProblemTemplateID string `bson:"problemTemplateID"` // 题目模板
		AnswerTemplateID  string `bson:"answerTemplateID"`  // 答案模板
	}{}
	err := userDB.C("products").Find(bson.M{
		"productID": productID,
	}).Select(bson.M{
		"problemTemplateID": 1,
		"answerTemplateID":  1,
	}).One(&templateID)
	if err != nil {
		return template, err
	}

	templateDBID := bson.ObjectIdHex(templateID.ProblemTemplateID)
	if templateType == "answer" {
		templateDBID = bson.ObjectIdHex(templateID.AnswerTemplateID)
	}
	if err := userDB.C("templates").FindId(templateDBID).One(&template); err != nil {
		return template, err
	}

	return template, nil
}

type timeSlice []time.Time

func (t timeSlice) Len() int           { return len(t) }
func (t timeSlice) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t timeSlice) Less(i, j int) bool { return t[i].Before(t[j]) }

func getDeliverMessage(productID string) (int, string, error) {
	// 获取交付的有关信息，包括优先级，和对用户的提示信息
	product := struct {
		DeliverType     string `json:"deliverType" bson:"deliverType"`         // 交付类型
		DeliverPriority int    `json:"deliverPriority" bson:"deliverPriority"` // 交付优先级
		DeliverTime     []struct {
			Day  time.Weekday `json:"day" bson:"day"`   // 周日0,周一到周六分别是1到6,
			Time string       `json:"time" bson:"time"` // 时间，格式按照"08:00:00"
		} `json:"deliverTime" bson:"deliverTime"` // 交付节点
		DeliverExpected int `json:"deliverExpected" bson:"deliverExpected"` // 交付预期，预期多少小时内
	}{}

	err := userDB.C("products").Find(bson.M{
		"productID": productID,
	}).One(&product)
	if err != nil {
		return 0, "", err
	}

	switch product.DeliverType {
	case "立即交付":
		const timePerPerson float64 = 0.5
		// 立即交付 优先级为0
		time := float64(documents.GetChannelWaitingCount(0)) * timePerPerson
		return product.DeliverPriority, "大约需要等待" + strconv.FormatFloat(time, 'f', -1, 64) + "分钟", nil

	case "预期交付":
		expectedTime := time.Now().Add(time.Duration(product.DeliverExpected) * time.Hour).Format("01-02 15:04")
		return product.DeliverPriority, expectedTime + "前可以获得", nil

	case "节点交付":
		// 节点list
		timeList := timeSlice{}
		todayWeekDay := time.Now().Weekday()
		for _, dt := range product.DeliverTime {
			// durationDays 为 dt 节点在今天的后多少天
			durationDays := int((dt.Day + 7 - todayWeekDay) % 7)
			deliverTime, err := time.Parse("2006-01-02 15:04:05", time.Now().AddDate(0, 0, durationDays).Format("2006-01-02 ")+dt.Time)
			if err != nil {
				return 0, "", err
			}
			if deliverTime.Before(time.Now()) {
				// 如果节点与今天星期相同，则要结合时间来比较，确保节点在今天之后
				deliverTime = deliverTime.AddDate(0, 0, 7)
			}
			timeList = append(timeList, deliverTime)
		}

		if len(timeList) <= 0 {
			return 0, "", fmt.Errorf("DeliverTime not set")
		}

		if len(timeList) <= 1 {
			// 保证有至少两个节点可以选
			timeList = append(timeList, timeList[0].AddDate(0, 0, 7))
		}
		sort.Sort(timeList)

		var selectedDelivertime time.Time
		if timeList[0].Sub(time.Now()).Hours() >= 8 {
			// 节点离现在8小时之后
			selectedDelivertime = timeList[0]
		} else {
			selectedDelivertime = timeList[1]
		}
		return product.DeliverPriority, selectedDelivertime.Format("01-02 15:04") + "前可以获得", nil
	}

	return 0, "", fmt.Errorf("wrong deliverType in productID %s", productID)
}

// GetProblemsFileHandler 获取题目文件
func GetProblemsFileHandler(c echo.Context) error {
	var id string
	if err := validation.ValidateUser(c, &id); err != nil {
		return err
	}

	if hasMarkTasks(id) {
		return utils.Forbidden("this student hasn't finished his upload tasks.")
	}

	type recvProbType struct {
		Type     string `json:"type"`
		Problems []struct {
			ProblemID string `json:"problemID"`
			SubIdx    int    `json:"subIdx"`
			Index     int    `json:"index"`
			Full      bool   `json:"full"`
		} `json:"problems"`
	}

	type recvType struct {
		ProductID string         `json:"productID"`
		Problems  []recvProbType `json:"problems"`
	}

	recvData := recvType{}

	if err := c.Bind(&recvData); err != nil {
		return utils.InvalidParams("invalid input, err:" + err.Error())
	}

	template, err := getTemplate(recvData.ProductID, "problem")
	if err != nil {
		log.Printf("can not get template of this product %s, err: %v\n", recvData.ProductID, err)
		return utils.NotFound("can not get template of this product")
	}

	type postType struct {
		LearnID  int64              `json:"learnID" bson:"learnID"`
		School   string             `json:"school" bson:"school"`
		Grade    string             `json:"grade" bson:"grade"`
		ClassID  int                `json:"classID" bson:"classID"`
		Name     string             `json:"name" bson:"realName"`
		Template templateDetailType `json:"template"`
		Problems []recvProbType     `json:"problems"`
	}

	postData := postType{}
	err = userDB.C("students").FindId(bson.ObjectIdHex(id)).Select(bson.M{
		"realName": 1,
		"grade":    1,
		"school":   1,
		"learnID":  1,
		"classID":  1,
	}).One(&postData)
	if err != nil {
		return utils.NotFound("can't find the information of this student.")
	}

	postData.Problems = recvData.Problems
	postData.Template = template

	priority, message, err := getDeliverMessage(recvData.ProductID)
	if err != nil {
		log.Printf("failed to get deliver message, err %v\n", err)
		return err
	}

	if err := setLastFileURLToBlank("lastProblemPDF", id); err != nil {
		log.Printf("failed to set last file URL to blank, err %v\n", err)
		return err
	}

	documents.PutToRequestQueue(postData, documents.ProblemRequest, priority, sendGettingFileRequestCallBackCreator("lastProblemPDF", id))

	return c.JSON(http.StatusOK, echo.Map{
		"message": message,
	})
}

// GetAnswersFileHandler 获取答案文件
func GetAnswersFileHandler(c echo.Context) error {
	var id string
	if err := validation.ValidateUser(c, &id); err != nil {
		return err
	}

	if hasMarkTasks(id) {
		return utils.Forbidden("this student hasn't finished his upload tasks.")
	}

	type recvProbType struct {
		ProblemID string `json:"problemID"`
		Location  string `json:"location"`
		Index     int    `json:"index"`
	}

	type recvType struct {
		ProductID string         `json:"productID"`
		Problems  []recvProbType `json:"problems"`
	}

	recvData := recvType{}

	if err := c.Bind(&recvData); err != nil {
		return utils.InvalidParams("invalid input, err:" + err.Error())
	}

	template, err := getTemplate(recvData.ProductID, "answer")
	if err != nil {
		log.Printf("can not get template of this product %s, err: %v\n", recvData.ProductID, err)
		return utils.NotFound("can not get template of this product")
	}

	type postType struct {
		LearnID  int64              `json:"learnID" bson:"learnID"`
		School   string             `json:"school" bson:"school"`
		Grade    string             `json:"grade" bson:"grade"`
		ClassID  int                `json:"classID" bson:"classID"`
		Name     string             `json:"name" bson:"realName"`
		Template templateDetailType `json:"template"`
		Problems []recvProbType     `json:"problems"`
	}

	postData := postType{}
	err = userDB.C("students").FindId(bson.ObjectIdHex(id)).Select(bson.M{
		"realName": 1,
		"grade":    1,
		"school":   1,
		"learnID":  1,
		"classID":  1,
	}).One(&postData)
	if err != nil {
		return utils.NotFound("can't find the information of this student.")
	}

	postData.Problems = recvData.Problems
	postData.Template = template

	priority, message, err := getDeliverMessage(recvData.ProductID)
	if err != nil {
		log.Printf("failed to get deliver message, err %v\n", err)
		return err
	}

	if err := setLastFileURLToBlank("lastAnswerPDF", id); err != nil {
		log.Printf("failed to set last file URL to blank, err %v\n", err)
		return err
	}

	documents.PutToRequestQueue(postData, documents.AnswerRequest, priority, sendGettingFileRequestCallBackCreator("lastAnswerPDF", id))

	return c.JSON(http.StatusOK, echo.Map{
		"message": message,
	})
}

// GetPointsFileHandler 获取知识点文件
func GetPointsFileHandler(c echo.Context) error {
	var id string
	if err := validation.ValidateUser(c, &id); err != nil {
		return err
	}

	contentServer := conf.AppConfig.FilesServer

	type recvType struct {
		ProblemID string `json:"problemID"`
		SubIdx    int    `json:"subIdx"`
		Index     int    `json:"index"`
	}
	recvData := []recvType{}

	if err := c.Bind(&recvData); err != nil {
		return utils.InvalidParams()
	}

	type postType struct {
		LearnID  int64      `json:"learnID" bson:"learnID"`
		School   string     `json:"school" bson:"school"`
		Grade    string     `json:"grade" bson:"grade"`
		ClassID  int        `json:"classID" bson:"classID"`
		Name     string     `json:"name" bson:"realName"`
		Problems []recvType `json:"problems"`
	}

	postData := postType{}
	err := userDB.C("students").FindId(bson.ObjectIdHex(id)).Select(bson.M{
		"realName": 1,
		"grade":    1,
		"school":   1,
		"classID":  1,
		"learnID":  1,
	}).One(&postData)
	if err != nil {
		return utils.NotFound("can't find the information of this student.")
	}

	postData.Problems = recvData

	result := struct {
		DocURL string `json:"docurl"`
		PdfURL string `json:"pdfurl"`
	}{"", ""}

	statusCode, err := utils.PostAndGetData("/getPointsFile/", postData, &result)
	if err != nil {
		log.Println(err)
		return err
	}
	if statusCode != 200 {
		log.Printf("Contacting with content server /getPointsFile/ status code: %d\n", statusCode)
		return echo.NewHTTPError(statusCode)
	}

	if result.DocURL == "" || result.PdfURL == "" {
		return utils.NotFound("Can't find files.")
	}

	result.DocURL = contentServer + result.DocURL
	result.PdfURL = contentServer + result.PdfURL
	return c.JSON(http.StatusOK, result)
}

// GetLastFileURLs 获取上一次纠错本相关文件下载URL
func GetLastFileURLs(c echo.Context) error {
	var id string
	if err := validation.ValidateUser(c, &id); err != nil {
		return err
	}

	data := struct {
		ProblemFileURL string `bson:"lastProblemPDF" json:"problemFileURL"`
		AnswerFileURL  string `bson:"lastAnswerPDF" json:"answerFileURL"`
	}{}

	err := userDB.C("students").FindId(bson.ObjectIdHex(id)).Select(bson.M{
		"lastProblemPDF": 1,
		"lastAnswerPDF":  1,
	}).One(&data)

	if err != nil {
		log.Printf("failed to get last file urls of id %s, err: %v", id, err)
		return err
	}

	return c.JSON(http.StatusOK, data)
}
