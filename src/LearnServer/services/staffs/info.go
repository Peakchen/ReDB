package staffs

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"LearnServer/conf"
	"LearnServer/models/contentDB"
	"LearnServer/models/userDB"
	
	"LearnServer/utils"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

func getChaptersSectionsHandler(c echo.Context) error {
	// 获取某个学期的章节信息
	semester := c.QueryParam("semester")

	chapterMinMax := struct {
		ChapterMin int `db:"chapMin"`
		ChapterMax int `db:"chapMax"`
	}{}

	db := contentDB.GetDB()
	if semester == "全部" {
		err := db.Get(&chapterMinMax, `SELECT MIN(chapMin) as chapMin, MAX(chapMax) as chapMax FROM extremumsChapMinMax;`)
		if err != nil {
			log.Printf("failed to get chapterMinMax, semester: %s, err: %v\n", semester, err)
			return err
		}
	} else {
		err := db.Get(&chapterMinMax, `SELECT chapMin, chapMax FROM extremumsChapMinMax WHERE grade = ? AND semester = ?;`, string([]rune(semester)[0:1]), string([]rune(semester)[1:2]))
		if err != nil {
			log.Printf("failed to get chapterMinMax, semester: %s, err: %v\n", semester, err)
			return err
		}
	}

	chaptersSections := []struct {
		Chapter     int    `json:"chapter" db:"chapter"`
		ChapterName string `json:"chapterName" db:"chapterName"`
		Section     int    `json:"section" db:"section"`
		SectionName string `json:"sectionName" db:"sectionName"`
	}{}

	err := db.Select(&chaptersSections, `
		SELECT DISTINCT c.num as chapter, c.name as chapterName, s.num as section, s.name as sectionName
		FROM chapters as c, sections as s
		WHERE s.chapNum = c.num AND c.num >= ? AND c.num <= ?;`, chapterMinMax.ChapterMin, chapterMinMax.ChapterMax)
	if err != nil {
		log.Printf("failed to get chapter section info, semester: %s, err: %v\n", semester, err)
	}

	return c.JSON(http.StatusOK, chaptersSections)
}

type examProblemType struct {
	ProblemID string `db:"problemID" bson:"problemID"`
	SubIdx    int    `db:"subIdx" bson:"subIdx"`
	Chapter   int    `db:"chapter"`
	Section   int    `db:"section"`
	TypeName  string `db:"typeName"`
}

func getClassExamProblems(schoolID string, grade string, class int) ([]examProblemType, error) {
	// 获取一个班级所有试卷的所有题目
	classExams := struct {
		ExamScoreRecords []struct {
			PaperID string `bson:"paperID"`
		} `bson:"examScoreRecords"`
	}{}
	err := userDB.C("classes").Find(bson.M{
		"schoolID": bson.ObjectIdHex(schoolID),
		"grade":    grade,
		"class":    class,
		"valid":    true,
	}).Select(bson.M{
		"examScoreRecords": 1,
	}).One(&classExams)
	if err != nil {
		log.Printf("failed to get exam score records of class, err:%v\n", err)
		return nil, err
	}

	examProblems := []examProblemType{}

	for _, exam := range classExams.ExamScoreRecords {
		problems := []examProblemType{}
		err := contentDB.GetDB().Select(&problems, `
			SELECT DISTINCT p.typeName, p.problemID, p.subIdx, e.chapterNum as chapter, e.sectionNum as section
			FROM examproblem as e, probtypes as p
			WHERE p.problemID = e.problemID AND e.examPaperID = ?;`, exam.PaperID)
		if err != nil {
			log.Printf("failed to get exam problems, err: %v\n", err)
			continue
		}
		examProblems = append(examProblems, problems...)
	}
	return examProblems, nil
}

func getProblemTypesHandler(c echo.Context) error {
	// 获取特定课时的题型信息


	chapterStart, err := strconv.Atoi(c.QueryParam("chapterStart"))
	if err != nil {
		return utils.InvalidParams("chapterStart is not a number!")
	}
	sectionStart, err := strconv.Atoi(c.QueryParam("sectionStart"))
	if err != nil {
		return utils.InvalidParams("sectionStart is not a number!")
	}
	lessonStart, err := strconv.Atoi(c.QueryParam("lessonStart"))
	if err != nil {
		return utils.InvalidParams("lessonStart is not a number!")
	}
	chapterEnd, err := strconv.Atoi(c.QueryParam("chapterEnd"))
	if err != nil {
		return utils.InvalidParams("chapterEnd is not a number!")
	}
	sectionEnd, err := strconv.Atoi(c.QueryParam("sectionEnd"))
	if err != nil {
		return utils.InvalidParams("sectionEnd is not a number!")
	}
	lessonEnd, err := strconv.Atoi(c.QueryParam("lessonEnd"))
	if err != nil {
		return utils.InvalidParams("lessonEnd is not a number!")
	}
	schoolID := c.QueryParam("schoolID")
	grade := c.QueryParam("grade")
	class, err := strconv.Atoi(c.QueryParam("class"))
	if err != nil {
		return utils.InvalidParams("class is not a number!")
	}

	type problemType struct {
		Chapter            int     `json:"chapter" db:"chapter"`
		Section            int     `json:"section" db:"section"`
		Lesson             int     `json:"lesson" db:"lesson"`                         // 课时序号
		TypeName           string  `json:"typeName" db:"typeName"`                     // 题型名称
		Priority           int     `json:"priority" db:"priority"`                     // 学习顺序
		PriorityTotal      int     `json:"priorityTotal" db:"priorityTotal"`           // 学习顺序总数
		Category           string  `json:"category" db:"category"`                     // 题型大类
		UnitExamProb       float32 `json:"unitExamProb" db:"unitExamProb"`             // 单元考试概率
		MidtermProb        float32 `json:"midtermProb" db:"midtermProb"`               // 期中考试概率
		FinalProb          float32 `json:"finalProb" db:"finalProb"`                   // 期末考试概率
		SeniorEntranceProb float32 `json:"seniorEntranceProb" db:"seniorEntranceProb"` // 中考概率
		ExamCount          int     `json:"examCount"`                                  // 已考次数
	}

	results := []problemType{}
	// 不能直接t.chapNum >= ? AND t.chapNum <= ? AND t.sectNum >= ? AND t.sectNum <= ? AND lesson >= ? AND lesson <= ?
	// 如范围为11.1.1 - 12.1.1，这样的判断无法获取11.2.2的数据
	err = contentDB.GetDB().Select(&results, `
		SELECT DISTINCT t.chapNum as chapter, t.sectNum as section, t.lesson, t.name as typeName, t.priority, t.category, COALESCE(e.unitExamProb, 0.01) as unitExamProb, COALESCE(e.midtermProb, 0.01) as midtermProb, COALESCE(e.finalProb, 0.01) as finalProb, COALESCE(e.seniorEntranceProb, 0.01) as seniorEntranceProb, (select MAX(priority) FROM typenames WHERE chapNum = t.chapNum AND sectNum = t.sectNum) as priorityTotal
		FROM typenames as t LEFT JOIN typeExamProbability as e
		ON t.name = e.typeName
		WHERE 
			((t.chapNum = ? AND t.sectNum = ? AND lesson >= ?) OR (t.chapNum = ? AND t.sectNum > ?) OR (t.chapNum > ?)) AND
			((t.chapNum < ?) OR (t.chapNum = ? AND t.sectNum < ?) OR (t.chapNum = ? AND t.sectNum = ? AND lesson <= ?));`,
		chapterStart, sectionStart, lessonStart, chapterStart, sectionStart, chapterStart,
		chapterEnd, chapterEnd, sectionEnd, chapterEnd, sectionEnd, lessonEnd)
	if err != nil {
		log.Printf("failed to get problem types info, err: %v\n", err)
		return err
	}

	examProblems, err := getClassExamProblems(schoolID, grade, class)
	if err != nil {
		return err
	}

	typeExamCountMap := make(map[string]int)
	for _, p := range examProblems {
		key := strconv.Itoa(p.Chapter) + strconv.Itoa(p.Section) + p.TypeName
		if count, ok := typeExamCountMap[key]; !ok {
			typeExamCountMap[key] = 1
		} else {
			typeExamCountMap[key] = count + 1
		}
	}
	for i, t := range results {
		key := strconv.Itoa(t.Chapter) + strconv.Itoa(t.Section) + t.TypeName
		if count, ok := typeExamCountMap[key]; ok {
			results[i].ExamCount = count
		}
	}

	return c.JSON(http.StatusOK, results)
}

func getProblemsOfTypeHandler(c echo.Context) error {
	// 获取特定范围的题目信息


	typeName := c.QueryParam("typeName")
	schoolID := c.QueryParam("schoolID")
	grade := c.QueryParam("grade")
	class, err := strconv.Atoi(c.QueryParam("class"))
	if err != nil {
		return utils.InvalidParams("class is not a number!")
	}

	results := []struct {
		How          string `json:"how" db:"how"`   // 出题方式
		HowCode      int    `json:"-" db:"howCode"` // 出题方式代码，用于排序
		ProblemID    string `json:"problemID" db:"problemID"`
		SubIdx       int    `json:"subIdx" db:"subIdx"`
		Used         int    `json:"used"` // 使用情况， 1 已布置， 2 已考  3 已布置已考 4 未使用
		HTMLFileName string `json:"-" db:"htmlFileName"`
		WordFileName string `json:"-" db:"wordFileName"`
		HTMLURL      string `json:"htmlURL"` // html 文件URL
		WordURL      string `json:"wordURL"` // word 文件URL
	}{}

	err = contentDB.GetDB().Select(&results, `
		SELECT DISTINCT p.problemID, p.subIdx, h.how, COALESCE(z.htmlFileName, '') as htmlFileName, COALESCE(z.filename, '') as wordFileName, IF(h.how='选择题', 1, IF(h.how='填空题', 2, 3)) as howCode
		FROM probtypes as p, hows as h LEFT JOIN problemzip as z
		ON h.problemID = z.problemID
		WHERE p.problemID = h.problemID AND p.typeName = ?
		ORDER BY howCode, p.problemID, p.subIdx;`, typeName)
	if err != nil {
		log.Printf("failed to get problems of type %s, err: %v\n", typeName, err)
		return err
	}

	examProblems, err := getClassExamProblems(schoolID, grade, class)
	if err != nil {
		return err
	}
	homeworkProblemsDB := struct {
		Problems []examProblemType `bson:"assignments"`
	}{}
	err = userDB.C("classes").Find(bson.M{
		"schoolID": bson.ObjectIdHex(schoolID),
		"grade":    grade,
		"class":    class,
		"valid":    true,
	}).Select(bson.M{
		"assignments": 1,
	}).One(&homeworkProblemsDB)
	if err != nil {
		log.Printf("failed to get assignments of class, err:%v\n", err)
		return err
	}

	homeworkProblems := homeworkProblemsDB.Problems
	// 统计题目是否出现
	examProblemsDoneMap := make(map[string]bool)
	homeworkProblemsDoneMap := make(map[string]bool)
	for _, p := range examProblems {
		examProblemsDoneMap[p.ProblemID+strconv.Itoa(p.SubIdx)] = true
	}
	for _, p := range homeworkProblems {
		homeworkProblemsDoneMap[p.ProblemID+strconv.Itoa(p.SubIdx)] = true
	}

	for i, p := range results {
		_, existInPaper := examProblemsDoneMap[p.ProblemID+strconv.Itoa(p.SubIdx)]
		_, existInHomework := homeworkProblemsDoneMap[p.ProblemID+strconv.Itoa(p.SubIdx)]
		switch {
		case !existInPaper && existInHomework:
			results[i].Used = 1
		case existInPaper && !existInHomework:
			results[i].Used = 2
		case existInPaper && existInHomework:
			results[i].Used = 3
		default:
			results[i].Used = 4
		}
		if p.HTMLFileName != "" {
			results[i].HTMLURL = conf.AppConfig.HTMLURL + results[i].HTMLFileName
		}
		if p.WordFileName != "" {
			results[i].WordURL = conf.AppConfig.ProblemDocURL + results[i].WordFileName
		}
	}

	return c.JSON(http.StatusOK, results)
}

func getProblemDocsZipHandler(c echo.Context) error {
	// 获取题目word文件压缩包


	type uploadType struct {
		ZipFileName string `json:"zipFileName"`
		TypeName    string `json:"typeName"`
		How         string `json:"how"`
		Problems    []struct {
			ProblemID string `json:"problemID"`
		} `json:"problems"`
	}

	uploadedData := uploadType{}
	if err := c.Bind(&uploadedData); err != nil {
		return utils.InvalidParams("invalid inputs!, error " + err.Error())
	}

	uploadedData.ZipFileName = fmt.Sprintf("%s（%s）", uploadedData.TypeName, uploadedData.How)

	contentServer := conf.AppConfig.FilesServer

	result := struct {
		URL string `json:"URL"`
	}{}

	statusCode, err := utils.PostAndGetData("/compressProblemDocs/", uploadedData, &result)
	if err != nil {
		log.Println(err)
		return err
	}
	if statusCode != 200 {
		log.Printf("Contacting with content server /compressProblemDocs/ status code: %d\n", statusCode)
		return echo.NewHTTPError(statusCode)
	}

	if result.URL == "" {
		return utils.NotFound("can't get packed file.")
	}

	result.URL = contentServer + result.URL
	return c.JSON(http.StatusOK, result)
}

func getBookProblemsInfoHandler(c echo.Context) error {
	// 获取某一本书某一页的题目


	bookID := c.QueryParam("bookID")
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		return utils.InvalidParams("page is invalid")
	}

	problems, err := contentDB.GetProblemsByBookPage(bookID, page, page)
	if err != nil {
		return err
	}
	if len(problems) == 0 {
		return utils.NotFound("no problems in this page in the book")
	}

	return c.JSON(http.StatusOK, problems)
}

func getPaperProblemsInfoHandler(c echo.Context) error {
	// 获取某个试卷的题目


	paperID := c.QueryParam("paperID")

	problems, err := contentDB.GetProblemsByPaper(paperID)
	if err != nil {
		return err
	}
	if len(problems) == 0 {
		return utils.NotFound("no problems in this page of the paper")
	}

	return c.JSON(http.StatusOK, problems)
}

func getKnowledgePointOfChapSectHandler(c echo.Context) error {
	// 获取某个某章节的知识点


	chapter, err := strconv.Atoi(c.QueryParam("chapter"))
	if err != nil {
		return utils.InvalidParams("chapter is invalid")
	}
	section, err := strconv.Atoi(c.QueryParam("section"))
	if err != nil {
		return utils.InvalidParams("section is invalid")
	}

	results := []struct {
		KnowledgeNum  int    `json:"knowledgeNum" db:"num" `   // 知识点序号
		KnowledgeName string `json:"knowledgeName" db:"name" ` // 知识点名称
	}{}

	err = contentDB.GetDB().Select(&results, "SELECT num, name FROM blocks where chapNum = ? and sectNum = ?;", chapter, section)
	if err != nil {
		log.Printf("failed to get knowledge from table blocks")
		return err
	}
	if len(results) == 0 {
		return utils.NotFound("no knowledge points in this chapter and section")
	}

	return c.JSON(http.StatusOK, results)
}
