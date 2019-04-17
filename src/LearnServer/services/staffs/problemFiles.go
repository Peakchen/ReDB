package staffs

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"LearnServer/conf"
	"LearnServer/models/userDB"
	"LearnServer/services/students/problempdfs"
	"LearnServer/tools/documents"
	"LearnServer/utils"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

func hasMarkTasks(learnID int) bool {
	// 判断是否有未标记的纠错本
	type taskInfoType struct {
		Time time.Time `bson:"time"`
	}

	tasksInfo := struct {
		Tasks []taskInfoType `bson:"tasks"`
	}{}
	err := userDB.C("students").Find(bson.M{
		"learnID": learnID,
		"valid":   true,
	}).Select(bson.M{
		"tasks": 1,
	}).One(&tasksInfo)
	if err != nil {
		log.Printf("fail to get tasks: learnID: %d, err: %v", learnID, err)
		return false
	}

	return len(tasksInfo.Tasks) > 0
}

func getTemplate(productID string, docType int) (TemplateDetailType, error) {
	// 获取用户产品中设定的模板 docType 文档类型 1 生成题目文件 2 生成答案文件
	template := TemplateDetailType{}

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
	if docType == 2 {
		templateDBID = bson.ObjectIdHex(templateID.AnswerTemplateID)
	}
	if err := userDB.C("templates").FindId(templateDBID).One(&template); err != nil {
		return template, err
	}

	return template, nil
}

func sendGettingFileRequestCallbackCreator(docType int, batchID string, studentIndex int) func(documents.DocumentURLType, int) {
	return func(urls documents.DocumentURLType, statusCode int) {
		batchInfo := batchDownloadMap[batchID]
		batchInfo.lastFileTime = time.Now()
		// 更新 batchDownloadMap 中的学生文件状态信息
		if docType == 1 {
			// 1 代表生成题目文件
			batchInfo.Students[studentIndex].ProblemFileStatus = statusCode == 200
			batchInfo.Students[studentIndex].ProblemStatusCode = statusCode
		} else if docType == 2 {
			// 2 代表生成答案文件
			batchInfo.Students[studentIndex].AnswerFileStatus = statusCode == 200
			batchInfo.Students[studentIndex].AnswerStatusCode = statusCode
		}
		batchDownloadMap[batchID] = batchInfo
	}
}

// createDocumentHandler 生成文档
func createDocumentHandler(c echo.Context) error {
	learnID, err := strconv.Atoi(c.Param("learnID"))
	if err != nil {
		return utils.InvalidParams("learnID is invalid.")
	}

	type recvType struct {
		ProductID string                                `json:"productID"`
		BatchID   string                                `json:"batchID"`
		DocType   int                                   `json:"docType"` // 文档类型 1 生成题目文件 2 生成答案文件
		Contents  []problempdfs.ProblemForCreatingFiles `json:"contents"`
	}

	recvData := recvType{}

	if err := c.Bind(&recvData); err != nil {
		return utils.InvalidParams("data input is invalid! error: " + err.Error())
	}

	template, err := getTemplate(recvData.ProductID, recvData.DocType)
	if err != nil {
		log.Printf("can not get template of this product %s and this student, id: %d, err: %v\n", recvData.ProductID, learnID, err)
		return utils.NotFound("can not get template of this product")
	}

	studentIndex, err := getStudentIndex(recvData.BatchID, learnID)
	if err != nil {
		return utils.InvalidParams(err.Error())
	}

	batchInfo := batchDownloadMap[recvData.BatchID]
	batchInfo.lastFileTime = time.Now()
	// 重置状态（使得可以重发请求）
	if recvData.DocType == 1 {
		// 1 代表生成题目文件
		batchInfo.Students[studentIndex].ProblemFileStatus = false
		batchInfo.Students[studentIndex].ProblemStatusCode = 0
	} else {
		// 2 代表生成答案文件
		batchInfo.Students[studentIndex].AnswerFileStatus = false
		batchInfo.Students[studentIndex].AnswerStatusCode = 0
	}

	if hasMarkTasks(learnID) {
		if recvData.DocType == 1 {
			// 1 代表生成题目文件
			batchInfo.Students[studentIndex].ProblemFileStatus = false
			batchInfo.Students[studentIndex].ProblemStatusCode = 403
		} else {
			// 2 代表生成答案文件
			batchInfo.Students[studentIndex].AnswerFileStatus = false
			batchInfo.Students[studentIndex].AnswerStatusCode = 403
		}
		batchDownloadMap[recvData.BatchID] = batchInfo
		return c.JSON(http.StatusOK, "this student hasn't finished his upload tasks.")
	}

	type postType struct {
		BatchID    string                                `json:"batchID"`
		DocType    int                                   `json:"docType"`
		LearnID    int64                                 `json:"learnID" bson:"learnID"`
		School     string                                `json:"school" bson:"school"`
		SchoolID   string                                `json:"schoolID" bson:"-"`
		SchoolIDDB bson.ObjectId                         `json:"-" bson:"schoolID"`
		Grade      string                                `json:"grade" bson:"grade"`
		ClassID    int                                   `json:"classID" bson:"classID"`
		Name       string                                `json:"name" bson:"realName"`
		Template   TemplateDetailType                    `json:"template"`
		Contents   []problempdfs.ProblemForCreatingFiles `json:"contents"`
	}

	postData := postType{}
	err = userDB.C("students").Find(bson.M{
		"learnID": learnID,
		"valid":   true,
	}).Select(bson.M{
		"realName": 1,
		"grade":    1,
		"school":   1,
		"learnID":  1,
		"classID":  1,
		"schoolID": 1,
	}).One(&postData)
	if err != nil {
		return utils.NotFound("can't find the information of this student.")
	}

	postData.BatchID = recvData.BatchID
	postData.DocType = recvData.DocType
	postData.SchoolID = postData.SchoolIDDB.Hex()
	postData.Contents = recvData.Contents
	variableMap := map[string]string{
		"{SCHOOL}":   postData.School,
		"{SCHOOLID}": postData.SchoolID,
		"{GRADE}":    postData.Grade,
		"{CLASS}":    strconv.Itoa(postData.ClassID),
		"{LEARNID}":  strconv.FormatInt(postData.LearnID, 10),
		"{NAME}":     postData.Name,
		"{DATE}":     time.Now().Format("20060102"),
	}
	template.FileName = getDocumentFileName(template.FileName, variableMap)
	postData.Template = template

	// 工作人员生成文档优先级为0
	documents.PutToRequestQueue(postData, documents.ProblemRequest, 0, sendGettingFileRequestCallbackCreator(recvData.DocType, recvData.BatchID, studentIndex))
	return c.JSON(http.StatusOK, "file is generating")
}

func getPackedFileHandler(c echo.Context) error {
	type uploadType struct {
		Grade   string `json:"grade"`
		Class   int    `json:"class"`
		BatchID string `json:"batchID"`
	}

	uploadedData := uploadType{}

	if err := c.Bind(&uploadedData); err != nil {
		return utils.InvalidParams("invalid learnID inputs!")
	}

	contentServer := conf.AppConfig.FilesServer

	result := struct {
		URL string `json:"URL"`
	}{}

	statusCode, err := utils.PostAndGetData("/getPackedFiles/", uploadedData, &result)
	if err != nil {
		log.Println(err)
		return err
	}
	if statusCode != 200 {
		log.Printf("Contacting with content server /getPackedFile/ status code: %d\n", statusCode)
		return echo.NewHTTPError(statusCode)
	}

	if result.URL == "" {
		return utils.NotFound("can't get this file.")
	}

	result.URL = contentServer + result.URL
	return c.JSON(http.StatusOK, result)
}
