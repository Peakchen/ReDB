package errorRateAnalysis

import (
	"log"
	"net/http"

	"LearnServer/conf"
	"LearnServer/utils"
	"github.com/labstack/echo"
)

// GetProblemsFileHandler 获取题目文件
func GetProblemsFileHandler(c echo.Context) error {

	contentServer := conf.AppConfig.FilesServer

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
		SchoolID string         `json:"schoolID"`
		Grade    string         `json:"grade"`
		Class    int            `json:"class"`
		Problems []recvProbType `json:"problems"`
	}

	recvData := recvType{}

	if err := c.Bind(&recvData); err != nil {
		return utils.InvalidParams("data input is invalid!")
	}

	type postType struct {
		SchoolID    string         `json:"schoolID"`
		Grade       string         `json:"grade"`
		ClassID     int            `json:"classID"`
		PageType    string         `json:"pageType"`
		Problems    []recvProbType `json:"problems"`
		ProblemType string         `json:"problemType"`
	}

	postData := postType{
		SchoolID:    recvData.SchoolID,
		Grade:       recvData.Grade,
		ClassID:     recvData.Class,
		PageType:    "A3",
		Problems:    recvData.Problems,
		ProblemType: "practice",
	}

	result := struct {
		DocURL string `json:"docurl"`
		PdfURL string `json:"pdfurl"`
	}{"", ""}

	statusCode, err := utils.PostAndGetData("/getProblemsFile/", postData, &result)
	if err != nil {
		log.Println(err)
		return err
	}
	if statusCode != 200 {
		log.Printf("Contacting with content server /getProblemsFile/ status code: %d\n", statusCode)
		return echo.NewHTTPError(statusCode)
	}

	if result.DocURL == "" || result.PdfURL == "" {
		return utils.NotFound("Can't find files.")
	}

	result.DocURL = contentServer + result.DocURL
	result.PdfURL = contentServer + result.PdfURL
	return c.JSON(http.StatusOK, result)
}

// GetAnswersFileHandler 获取答案文件
func GetAnswersFileHandler(c echo.Context) error {

	contentServer := conf.AppConfig.FilesServer

	type recvProbType struct {
		ProblemID string `json:"problemID"`
		Location  string `json:"location"`
		Index     int    `json:"index"`
	}

	type recvType struct {
		SchoolID string         `json:"schoolID"`
		Grade    string         `json:"grade"`
		Class    int            `json:"class"`
		Problems []recvProbType `json:"problems"`
	}

	recvData := recvType{}

	if err := c.Bind(&recvData); err != nil {
		return utils.InvalidParams("invalid inputs!")
	}

	type postType struct {
		SchoolID    string         `json:"schoolID"`
		Grade       string         `json:"grade"`
		ClassID     int            `json:"classID"`
		PageType    string         `json:"pageType"`
		Problems    []recvProbType `json:"problems"`
		ProblemType string         `json:"problemType"`
	}

	postData := postType{
		SchoolID:    recvData.SchoolID,
		Grade:       recvData.Grade,
		ClassID:     recvData.Class,
		PageType:    "A3",
		Problems:    recvData.Problems,
		ProblemType: "practice",
	}

	result := struct {
		DocURL string `json:"docurl"`
		PdfURL string `json:"pdfurl"`
	}{"", ""}

	statusCode, err := utils.PostAndGetData("/getAnswersFile/", postData, &result)
	if err != nil {
		log.Println(err)
		return err
	}
	if statusCode != 200 {
		log.Printf("Contacting with content server /getAnswersFile/ status code: %d\n", statusCode)
		return echo.NewHTTPError(statusCode)
	}

	if result.DocURL == "" || result.PdfURL == "" {
		return utils.NotFound("Can't find files.")
	}

	result.DocURL = contentServer + result.DocURL
	result.PdfURL = contentServer + result.PdfURL
	return c.JSON(http.StatusOK, result)
}

// GetPackedFileHandler 获取题目答案压缩包文件
func GetPackedFileHandler(c echo.Context) error {

	type uploadType struct {
		SchoolID string `json:"schoolID"`
		Grade    string `json:"grade"`
		Class    int    `json:"class"`
	}

	input := uploadType{}

	if err := c.Bind(&input); err != nil {
		return utils.InvalidParams("invalid inputs!" + err.Error())
	}

	contentServer := conf.AppConfig.FilesServer

	result := struct {
		URL string `json:"URL"`
	}{}

	statusCode, err := utils.PostAndGetData("/getPackedPracticeFiles/", input, &result)
	if err != nil {
		log.Println(err)
		return err
	}
	if statusCode != 200 {
		log.Printf("Contacting with content server /getPackedPracticeFiles/ status code: %d\n", statusCode)
		return echo.NewHTTPError(statusCode)
	}

	if result.URL == "" {
		return utils.NotFound("can't get this file.")
	}

	result.URL = contentServer + result.URL
	return c.JSON(http.StatusOK, result)
}
