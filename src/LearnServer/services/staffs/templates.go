package staffs

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"LearnServer/conf"
	"LearnServer/models/userDB"
	
	"LearnServer/services/students/problempdfs"
	"LearnServer/utils"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

// operationType 模板当中的操作
type operationType struct {
	Type       int      `json:"type" bson:"type"`             // 1 变更格式 2 添加内容
	Font       string   `json:"font" bson:"font"`             // 字体（仅当type为1变更格式时有效）
	FontSize   float64  `json:"fontSize" bson:"fontSize"`     // 字号（仅当type为1变更格式时有效）
	FontEffect []string `json:"fontEffect" bson:"fontEffect"` // 字体效果 bold 加粗， underlined 下划线 italic 倾斜（仅当type为1变更格式时有效）
	Alignment  int      `json:"alignment" bson:"alignment"`   // 对齐，0 左对齐 1 居中对齐 2 右对齐（仅当type为1变更格式时有效）
	RowSpacing float64  `json:"rowSpacing" bson:"rowSpacing"` // 行距（仅当type为1变更格式时有效）
	Content    string   `json:"content" bson:"content"`       // 内容（仅当type为2添加内容时有效）
}

// listItemType 列表变量每一项的定义
type listItemType struct {
	BlankLinesForSelectProFillPro int             `json:"blankLinesForSelectProFillPro" bson:"blankLinesForSelectProFillPro"` // 选择题或者填空题空多少行
	BlankLinesForEachSubPro       int             `json:"blankLinesForEachSubPro" bson:"blankLinesForEachSubPro"`             // 大题每小问空多少行
	MinBlankLinesOfPro            int             `json:"minBlankLinesOfPro" bson:"minBlankLinesOfPro"`                       // 大题空的总行数最少值
	Operations                    []operationType `json:"operations" bson:"operations"`                                       // 操作
}

// TemplateDetailType 模板信息
// 包括所有需要配置的字段，除了 templateID 与 date
type TemplateDetailType struct {
	Name             string          `json:"name" bson:"name"`                         // 模板名称
	Info             string          `json:"info" bson:"info"`                         // 模板说明
	Type             string          `json:"type" bson:"type"`                         // 模板类型
	FileName         string          `json:"fileName" bson:"fileName"`                 // 文件名称
	PageType         string          `json:"pageType" bson:"pageType"`                 // 纸张大小，"A3"或者"A4"
	PageDirection    string          `json:"pageDirection" bson:"pageDirection"`       // 纸张方向
	ColumnCount      int             `json:"columnCount" bson:"columnCount"`           // 分栏数
	MarginTop        float64         `json:"marginTop" bson:"marginTop"`               // 上下页边距
	MarginLeft       float64         `json:"marginLeft" bson:"marginLeft"`             // 左右页边距
	CheckProblemList listItemType    `json:"checkProblemList" bson:"checkProblemList"` // CHECK_PROBLEM_LIST变量每项定义
	ContentList      listItemType    `json:"contentList" bson:"contentList"`           // CONTENT_LIST变量每项定义
	Operations       []operationType `json:"operations" bson:"operations"`             // 文档正文的操作
}

func getDocumentFileName(fileNameInTemplate string, variableMap map[string]string) string {
	// 生成需要下载的文件的文件名， variableMap 为支持的变量与对应值的 map
	fileName := fileNameInTemplate
	for variableName, value := range variableMap {
		fileName = strings.Replace(fileName, variableName, value, -1)
	}
	return fileName
}

func uploadTemplateHandler(c echo.Context) error {
	// 上传新的模板


	type templateType struct {
		TemplateDetailType `bson:",inline"`
		Date               time.Time `bson:"date"`
	}

	uploadData := templateType{}
	if err := c.Bind(&uploadData); err != nil {
		return utils.InvalidParams("invalid input, error: " + err.Error())
	}

	uploadData.Date = time.Now()

	if err := userDB.C("templates").Insert(uploadData); err != nil {
		log.Printf("cannot save new template, err: %v\n", err)
		return err
	}

	return c.JSON(http.StatusOK, "successfully add a template")
}

func listTemplatesHandler(c echo.Context) error {
	// 获取符合条件的模板信息


	templateDataType := c.QueryParam("type")

	type templateType struct {
		TemplateDetailType `bson:",inline"`
		DatabaseID         bson.ObjectId `json:"-" bson:"_id"`
		TemplateID         string        `json:"templateID" bson:"-"`
		DateUnix           int64         `json:"date" bson:"-"` // 设计日期
		Date               time.Time     `json:"-" bson:"date"`
	}

	templates := []templateType{}

	query := bson.M{}
	if templateDataType != "all" {
		query["type"] = templateDataType
	}

	if err := userDB.C("templates").Find(query).All(&templates); err != nil {
		log.Printf("finding templates failed, err: %v\n", err)
		return utils.NotFound("finding templates failed")
	}

	for i := range templates {
		templates[i].DateUnix = templates[i].Date.Unix()
		templates[i].TemplateID = templates[i].DatabaseID.Hex()
	}

	return c.JSON(http.StatusOK, templates)
}

func retriveTemplateHandler(c echo.Context) error {
	// 根据ID获取一个模板信息


	type templateType struct {
		TemplateDetailType `bson:",inline"`
		DatabaseID         bson.ObjectId `json:"-" bson:"_id"`
		TemplateID         string        `json:"templateID" bson:"-"`
		DateUnix           int64         `json:"date" bson:"-"` // 设计日期
		Date               time.Time     `json:"-" bson:"date"`
	}

	result := templateType{}
	templateID := c.Param("templateID")

	err := userDB.C("templates").FindId(bson.ObjectIdHex(templateID)).One(&result)
	if err != nil {
		log.Printf("can not find this template, err: %v\n", err)
		return utils.NotFound("can not find this template")
	}

	result.DateUnix = result.Date.Unix()
	result.TemplateID = result.DatabaseID.Hex()

	return c.JSON(http.StatusOK, result)
}

func updateTemplateHandler(c echo.Context) error {
	// 修改一个模板


	type templateType struct {
		TemplateDetailType `bson:",inline"`
		Date               time.Time `bson:"date"`
		PdfFileURL         string    `json:"-" bson:"pdfURL"`
		DocFileURL         string    `json:"-" bson:"docURL"` // 重置预览URL
	}

	templateID := c.Param("templateID")

	uploadData := templateType{}
	if err := c.Bind(&uploadData); err != nil {
		return utils.InvalidParams("invalid input, error: " + err.Error())
	}

	uploadData.Date = time.Now()

	err := userDB.C("templates").UpdateId(bson.ObjectIdHex(templateID), bson.M{
		"$set": uploadData,
	})
	if err != nil {
		log.Printf("can not update this template, err: %v\n", err)
		return err
	}

	return c.JSON(http.StatusOK, "Successfully updated this template")
}

func deleteTemplateHandler(c echo.Context) error {
	// 删除一个模板


	templateID := c.Param("templateID")
	if err := userDB.C("templates").RemoveId(bson.ObjectIdHex(templateID)); err != nil {
		log.Printf("can not delete this template, err: %v\n", err)
		return err
	}

	return c.JSON(http.StatusOK, "Successfully deleted this template")
}

func previewTemplateHandler(c echo.Context) error {
	// 预览模板


	type templateType struct {
		TemplateDetailType `bson:",inline"`
		PdfFileURL         string `bson:"pdfURL"`
		DocFileURL         string `bson:"docURL"`
	}

	template := templateType{}
	templateID := c.Param("templateID")

	err := userDB.C("templates").FindId(bson.ObjectIdHex(templateID)).One(&template)
	if err != nil {
		log.Printf("can not find this template, err: %v\n", err)
		return utils.NotFound("can not find this template")
	}

	if template.PdfFileURL != "" {
		return c.JSON(http.StatusOK, echo.Map{
			"pdfurl": template.PdfFileURL,
			"docurl": template.DocFileURL,
		})
	}

	fakeLearnID, err := userDB.GetNewID("fakeLearnIDs")
	variableMap := map[string]string{
		"{SCHOOL}":   "测试学校",
		"{SCHOOLID}": "abcdefg",
		"{GRADE}":    "八",
		"{CLASS}":    "5",
		"{LEARNID}":  strconv.FormatInt(fakeLearnID, 10),
		"{NAME}":     "张三",
		"{DATE}":     time.Now().Format("20060102"),
	}
	// 将模板预览文件存储在 templates 下
	template.FileName = "templates/" + getDocumentFileName(template.FileName, variableMap)
	if err != nil {
		log.Printf("failed to create fake learnID, err %v\n", err)
		return err
	}
	fakeLearnID = -fakeLearnID
	contentServer := conf.AppConfig.FilesServer
	var postData interface{}
	type basicProblemForCreatingFiles = problempdfs.BasicProblemForCreatingFiles

	type problemForCreatingFiles = problempdfs.ProblemForCreatingFiles

	postData = struct {
		BatchID  string                    `json:"batchID"`
		DocType  int                       `json:"docType"`
		LearnID  int64                     `json:"learnID"`
		School   string                    `json:"school"`
		SchoolID string                    `json:"schoolID"`
		Grade    string                    `json:"grade"`
		ClassID  int                       `json:"classID"`
		Name     string                    `json:"name"`
		Contents []problemForCreatingFiles `json:"contents"`
		Template templateType              `json:"template"`
	}{
		BatchID:  "000",
		DocType:  1,
		LearnID:  fakeLearnID,
		School:   "测试学校",
		SchoolID: "abcdefg",
		Grade:    "八",
		ClassID:  5,
		Name:     "张三",
		Contents: []problemForCreatingFiles{
			problemForCreatingFiles{
				basicProblemForCreatingFiles{"类型1", "八下全品作业本", 1, "栏目1", 1, "90001", []int{-1}, true, "选择题", "依据1"},
				[]basicProblemForCreatingFiles{
					basicProblemForCreatingFiles{"类型1", "八下课本", 4, "栏目1", 3, "26008", []int{1, 3}, false, "计算题", "依据2"},
					basicProblemForCreatingFiles{"类型2", "八下全品作业本", 5, "栏目1", 5, "26010", []int{1, 2}, true, "计算题", "依据3"},
				},
			},
			problemForCreatingFiles{
				basicProblemForCreatingFiles{"类型1", "八下课本", 4, "栏目1", 3, "26008", []int{1, 3}, false, "计算题", "依据2"},
				[]basicProblemForCreatingFiles{
					basicProblemForCreatingFiles{"类型1", "八下全品作业本", 1, "栏目1", 1, "90001", []int{-1}, true, "选择题", "依据1"},
					basicProblemForCreatingFiles{"类型2", "八下全品作业本", 5, "栏目1", 5, "26010", []int{1, 2}, true, "计算题", "依据3"},
				},
			},
			problemForCreatingFiles{
				basicProblemForCreatingFiles{"类型2", "八下全品作业本", 5, "栏目1", 5, "26010", []int{1, 2}, true, "计算题", "依据3"},
				[]basicProblemForCreatingFiles{
					basicProblemForCreatingFiles{"类型1", "八下全品作业本", 1, "栏目1", 1, "90001", []int{-1}, true, "选择题", "依据1"},
					basicProblemForCreatingFiles{"类型1", "八下课本", 4, "栏目1", 3, "26008", []int{1, 3}, false, "计算题", "依据2"},
				},
			},
		},
		Template: template,
	}
	result := struct {
		DocURL string `json:"docurl"`
		PdfURL string `json:"pdfurl"`
	}{"", ""}

	statusCode, err := utils.PostAndGetData("/createDocuments/", postData, &result)
	if err != nil {
		log.Println(err)
		return err
	}
	if statusCode != 200 {
		log.Printf("Contacting with content server %s status code: %d\n", "/createDocuments/", statusCode)
		return echo.NewHTTPError(statusCode)
	}

	if result.DocURL == "" || result.PdfURL == "" {
		return utils.NotFound("Can't find files.")
	}

	result.DocURL = contentServer + result.DocURL
	result.PdfURL = contentServer + result.PdfURL

	err = userDB.C("templates").UpdateId(bson.ObjectIdHex(templateID), bson.M{
		"$set": bson.M{
			"pdfURL": result.PdfURL,
			"docURL": result.DocURL,
		},
	})
	if err != nil {
		log.Printf("can not update template preview url, err: %v\n", err)
	}
	return c.JSON(http.StatusOK, result)
}
