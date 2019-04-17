package staffs

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"LearnServer/models/contentDB"
	"LearnServer/models/userDB"

	"LearnServer/utils"
	set "github.com/deckarep/golang-set"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

func getClassPapers(schoolID string, grade string, class int) ([]contentDB.PaperDetailType, error) {
	type papersType struct {
		PaperIDs []interface{} `bson:"papers"`
	}
	classPaperIDs := []papersType{}
	var selector bson.M
	if class != 0 {
		selector = bson.M{
			"schoolID": bson.ObjectIdHex(schoolID),
			"grade":    grade,
			"class":    class,
			"valid":    true,
		}
	} else {
		// 选择全部班级
		selector = bson.M{
			"schoolID": bson.ObjectIdHex(schoolID),
			"grade":    grade,
			"valid":    true,
		}
	}
	err := userDB.C("classes").Find(selector).All(&classPaperIDs)
	if err != nil {
		return nil, err
	}
	if len(classPaperIDs) == 0 {
		return nil, fmt.Errorf("not found")
	}

	// paperIDResults 得到 paperID 交集的结果的 slice
	paperIDResults := []string{}
	// 存储 paperID 的交集
	paperIDSets := set.NewSetFromSlice(classPaperIDs[0].PaperIDs)
	for _, b := range classPaperIDs {
		paperIDSets = paperIDSets.Intersect(set.NewSetFromSlice(b.PaperIDs))
	}
	it := paperIDSets.Iterator()
	for paperID := range it.C {
		paperIDResults = append(paperIDResults, paperID.(string))
	}

	return contentDB.GetPapersByPaperID(paperIDResults), nil
}

func getPapersOfClassHandler(c echo.Context) error {
	schoolID := c.QueryParam("schoolID")
	grade := c.QueryParam("grade")
	class, err := strconv.Atoi(c.QueryParam("class"))
	if err != nil {
		return utils.InvalidParams("class is not a number!")
	}
	papers, err := getClassPapers(schoolID, grade, class)
	if err != nil {
		if err.Error() == "not found" {
			return utils.NotFound("this class has no papers")
		}
		return err
	}
	return c.JSON(http.StatusOK, papers)
}

func getPapersOfClassForMarkingScoreHandler(c echo.Context) error {
	// 获取某个班级的用于标记考试成绩的试卷信息
	schoolID := c.QueryParam("schoolID")
	grade := c.QueryParam("grade")
	class, err := strconv.Atoi(c.QueryParam("class"))
	if err != nil {
		return utils.InvalidParams("class is not a number!")
	}
	papers, err := getClassPapers(schoolID, grade, class)
	if err != nil {
		if err.Error() == "not found" {
			return utils.NotFound("this class has no papers")
		}
		return err
	}

	type examType struct {
		Time time.Time `bson:"time"`
	}
	classExamRecord := make(map[string]map[string]examType)

	err = userDB.C("classes").Find(bson.M{
		"schoolID": bson.ObjectIdHex(schoolID),
		"grade":    grade,
		"class":    class,
		"valid":    true,
	}).Select(bson.M{
		"examScoreRecords": 1,
	}).One(&classExamRecord)
	if err != nil {
		log.Printf("getting class examRecord failed")
		return err
	}
	examRecord, ok := classExamRecord["examScoreRecords"]
	if !ok {
		// 置空
		examRecord = make(map[string]examType)
	}

	type resultPaperType struct {
		PaperID   string `json:"paperID"`
		Name      string `json:"name"`
		FullScore int    `json:"fullScore"`
		Marked    bool   `json:"marked"`
	}
	results := []resultPaperType{}
	// 对 papers 顺序进行调整，将已经标记的试卷放后面
	markedPapers := []resultPaperType{}
	for _, p := range papers {
		if _, ok := examRecord["paperID"+p.PaperID]; !ok {
			// 还没标记过
			results = append(results, resultPaperType{
				PaperID:   p.PaperID,
				Name:      p.Name,
				FullScore: p.FullScore,
				Marked:    false,
			})
		} else {
			// 已经标记过
			markedPapers = append(markedPapers, resultPaperType{
				PaperID:   p.PaperID,
				Name:      p.Name,
				FullScore: p.FullScore,
				Marked:    true,
			})
		}
	}
	results = append(results, markedPapers...)

	return c.JSON(http.StatusOK, results)
}

func deletePaperFromClassHandler(c echo.Context) error {
	type inputType struct {
		SchoolID string `json:"schoolID"`
		Grade    string `json:"grade"`
		Class    int    `json:"class"`
		PaperID  string `json:"paperID"`
	}

	input := inputType{}
	if err := c.Bind(&input); err != nil {
		return utils.InvalidParams("invalid input!" + err.Error())
	}

	var selector bson.M
	if input.Class != 0 {
		selector = bson.M{
			"schoolID": bson.ObjectIdHex(input.SchoolID),
			"grade":    input.Grade,
			"class":    input.Class,
			"valid":    true,
		}
	} else {
		// 选择全部班级
		selector = bson.M{
			"schoolID": bson.ObjectIdHex(input.SchoolID),
			"grade":    input.Grade,
			"valid":    true,
		}
	}

	_, err := userDB.C("classes").UpdateAll(selector, bson.M{
		"$pull": bson.M{
			"papers": input.PaperID,
		},
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "successfully deleted!")
}

// getNotMarkedPapersHandler 获取有哪些未标记试卷
func getNotMarkedPapersHandler(c echo.Context) error {
	learnID, err := strconv.Atoi(c.Param("learnID"))
	if err != nil {
		return utils.InvalidParams("learnID is invalid")
	}
	studentID, err := getStudentIDByLearnID(learnID)
	if err != nil {
		return utils.NotFound("can not find the information of this learnID")
	}

	paperIDs, err := userDB.GetNotMarkedPaperIDs(studentID)
	if err != nil {
		log.Printf("getting not marked papers of student id %s failed, err %v", studentID, err)
		return err
	}
	if len(paperIDs) <= 0 {
		return utils.NotFound("no papers.")
	}

	return c.JSON(http.StatusOK, contentDB.GetPapersByPaperID(paperIDs))
}

// 获取有哪些已经标记的试卷
func getMarkedPapersHandler(c echo.Context) error {
	learnID, err := strconv.Atoi(c.Param("learnID"))
	if err != nil {
		return utils.InvalidParams("learnID is invalid")
	}
	studentID, err := getStudentIDByLearnID(learnID)
	if err != nil {
		return utils.NotFound("can not find the information of this learnID")
	}

	paperIDs, err := userDB.GetMarkedPaperIDs(studentID)
	if err != nil {
		log.Printf("getting marked papers of student id %s failed, err %v", studentID, err)
		return err
	}
	if len(paperIDs) <= 0 {
		return utils.NotFound("no papers.")
	}

	return c.JSON(http.StatusOK, contentDB.GetPapersByPaperID(paperIDs))
}
