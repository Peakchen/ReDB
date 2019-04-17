package documents

import (
	"log"

	// "LearnServer/utils"
	"LearnServer/utils"
)

const (
	// ProblemRequest 在 requests 中对应了 sendRequestToGetProblemFile 的下标
	ProblemRequest = iota
	// AnswerRequest 在 requests 中对应了 sendRequestToGetAnswerFile 的下标
	AnswerRequest
)

// DocumentURLType 生成文档返回的文档下载URL类型
type DocumentURLType struct {
	DocURL string `json:"docurl"`
	PdfURL string `json:"pdfurl"`
}

// TODO: 无需区分题目与答案文档

// sendRequestToGetProblemFile 发送请求生成题目文档
func sendRequestToGetProblemFile(postData interface{}, callback func(fileURL DocumentURLType, statusCode int)) {
	result := DocumentURLType{"", ""}
	statusCode, err := utils.PostAndGetData("/createDocuments/", postData, &result)
	if err != nil {
		log.Println(err)
		return
	}
	if statusCode != 200 {
		log.Printf("Contacting with content server /createDocuments/ status code: %d\n", statusCode)
	}

	// 成功之后调用callback
	callback(result, statusCode)
}

// sendRequestToGetAnswerFile 发送请求生成题目文档
func sendRequestToGetAnswerFile(postData interface{}, callback func(fileURL DocumentURLType, statusCode int)) {
	result := DocumentURLType{"", ""}
	statusCode, err := utils.PostAndGetData("/createDocuments/", postData, &result)
	if err != nil {
		log.Println(err)
		return
	}
	if statusCode != 200 {
		log.Printf("Contacting with content server /createDocuments/ status code: %d\n", statusCode)
	}

	callback(result, statusCode)
}

var requests = [2]func(interface{}, func(DocumentURLType, int)){sendRequestToGetProblemFile, sendRequestToGetAnswerFile}
