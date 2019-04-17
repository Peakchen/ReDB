package staffs

import (
	"log"
	"net/http"
	"sort"
	"strconv"

	"LearnServer/models/contentDB"
	"LearnServer/models/userDB"

	"LearnServer/utils"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

type targetType struct {
	Status     bool   `json:"status" bson:"-"`
	Chapter    int    `json:"chapter" bson:"chapter" db:"chapNum"` // 章
	Section    int    `json:"section" bson:"section" db:"sectNum"` // 节
	Typename   string `json:"typename" bson:"typename" db:"name"`  // 题型名称
	OriginalKP string `json:"originalKP" bson:"-" db:"originalKP"` // 最新原始知识点
	Lesson     int    `json:"-" bson:"-" db:"lesson"`              // 课时
	Priority   int    `json:"-" bson:"-" db:"priority"`            // 学习顺序
}

type targetSlice []targetType

func (c targetSlice) Len() int { return len(c) }

func (c targetSlice) Swap(i int, j int) { c[i], c[j] = c[j], c[i] }

func (c targetSlice) Less(i int, j int) bool {
	if c[i].Status && !c[j].Status {
		return true
	}
	if !c[i].Status && c[j].Status {
		return false
	}
	if c[i].Chapter != c[j].Chapter {
		return c[i].Chapter < c[j].Chapter
	}
	if c[i].Section != c[j].Section {
		return c[i].Section < c[j].Section
	}
	if c[i].Lesson != c[j].Lesson {
		return c[i].Lesson < c[j].Lesson
	}
	return c[i].Priority < c[j].Priority
}

func getTargetsWithStatus(allTargets targetSlice, statusTrueTargets targetSlice) targetSlice {
	// 在　allTargets　中找出在　statusTrueTargets　中的 targets，将 status　赋值 true　其余赋值 false
	result := make(targetSlice, len(allTargets))
	for i, target := range allTargets {
		status := false
		for _, trueTarget := range statusTrueTargets {
			if target.Chapter == trueTarget.Chapter && target.Section == trueTarget.Section && target.Typename == trueTarget.Typename {
				status = true
				break
			}
		}
		result[i] = targetType{
			Chapter:    target.Chapter,
			Section:    target.Section,
			Typename:   target.Typename,
			OriginalKP: target.OriginalKP,
			Status:     status,
			Lesson:     target.Lesson,
			Priority:   target.Priority,
		}
	}
	return result
}

func getClassTargetsHandler(c echo.Context) error {
	// 获取某个班级的目标规划信息
	schoolID := c.QueryParam("schoolID")
	grade := c.QueryParam("grade")
	classID, err := strconv.Atoi(c.QueryParam("class"))
	if err != nil {
		return utils.InvalidParams("class is invalid!")
	}
	levelStr := c.QueryParam("level")
	chapterStr := c.QueryParam("chapter")
	sectionStr := c.QueryParam("section")
	typename := c.QueryParam("typename")
	semester := c.QueryParam("semester")
	exam := c.QueryParam("exam")

	chapter, err := strconv.Atoi(chapterStr)
	if err != nil {
		return utils.InvalidParams("chapter is not a number")
	}
	section, err := strconv.Atoi(sectionStr)
	if err != nil {
		return utils.InvalidParams("section is not a number")
	}

	db := contentDB.GetDB()
	var allTargets targetSlice
	var whereClause string
	if chapter <= 0 {
		// 全部章
		whereClause += "chapNum > ?"
		// 先要有一个 > ? 与 chapter 匹配
		if semester != "全部" {
			// 只获取那个学期的全部章
			chapterInfo := struct {
				Min int `db:"chapMin"`
				Max int `db:"chapMax"`
			}{}
			err := db.Get(&chapterInfo, `SELECT chapMin, chapMax FROM extremumsChapMinMax WHERE grade = ? AND semester = ?;`, string([]rune(semester)[0:1]), string([]rune(semester)[1:2]))
			if err != nil {
				log.Printf("failed to get chapterMinMax, semester: %s, err: %v\n", semester, err)
				return err
			}
			whereClause += " AND chapNum >= " + strconv.Itoa(chapterInfo.Min) + " AND chapNum <= " + strconv.Itoa(chapterInfo.Max)
		}
	} else {
		whereClause += "chapNum = ?"
	}
	if section <= 0 {
		// 全部节
		whereClause += " AND sectNum > ?"
	} else {
		whereClause += " AND sectNum = ?"
	}
	whereClause += " AND name LIKE ?;"
	err = db.Select(&allTargets, "SELECT chapNum, sectNum, name, originalKP, lesson, priority FROM typenames WHERE "+whereClause, chapter, section, "%"+typename+"%")
	if err != nil {
		log.Printf("failed to get allTargets, err %v\n", err)
		return err
	}

	statusTrueTargetsMap := make(map[string]targetSlice)
	queryKey := "level" + levelStr + "targets" + exam
	statusTrueTargetsMap[queryKey] = targetSlice{}
	err = userDB.C("classes").Find(bson.M{
		"schoolID": bson.ObjectIdHex(schoolID),
		"grade":    grade,
		"class":    classID,
		"valid":    true,
	}).Select(bson.M{
		queryKey: 1,
	}).One(&statusTrueTargetsMap)
	if err != nil {
		log.Printf("failed to get statusTrueTargets, err %v\n", err)
		return err
	}

	targets := getTargetsWithStatus(allTargets, statusTrueTargetsMap[queryKey])
	sort.Sort(targets)

	return c.JSON(http.StatusOK, targets)
}

func addClassTargetsHandler(c echo.Context) error {
	// 为某个班级层级添加目标

	type uploadType struct {
		SchoolID string      `json:"schoolID"`
		Grade    string      `json:"grade"`
		ClassID  int         `json:"class"`
		Level    int         `json:"level"`
		Exam     string      `json:"exam"`
		Targets  targetSlice `json:"targets"`
	}

	uploadData := uploadType{}
	if err := c.Bind(&uploadData); err != nil {
		return utils.InvalidParams("invalid input!" + err.Error())
	}

	fieldName := "level" + strconv.Itoa(uploadData.Level) + "targets" + uploadData.Exam
	_, err := userDB.C("classes").Upsert(bson.M{
		"schoolID": bson.ObjectIdHex(uploadData.SchoolID),
		"grade":    uploadData.Grade,
		"class":    uploadData.ClassID,
		"valid":    true,
	}, bson.M{
		"$addToSet": bson.M{
			fieldName: bson.M{
				"$each": uploadData.Targets,
			},
		},
	})
	if err != nil {
		log.Printf("failed to add targets, err %v\n", err)
		return err
	}

	return c.JSON(http.StatusOK, "successfully added")
}

func deleteClassTargetsHandler(c echo.Context) error {
	// 为某个班级层级移除目标

	type uploadType struct {
		SchoolID string      `json:"schoolID"`
		Grade    string      `json:"grade"`
		ClassID  int         `json:"class"`
		Level    int         `json:"level"`
		Exam     string      `json:"exam"`
		Targets  targetSlice `json:"targets"`
	}

	uploadData := uploadType{}
	if err := c.Bind(&uploadData); err != nil {
		return utils.InvalidParams("invalid input!" + err.Error())
	}

	fieldName := "level" + strconv.Itoa(uploadData.Level) + "targets" + uploadData.Exam
	_, err := userDB.C("classes").Upsert(bson.M{
		"schoolID": bson.ObjectIdHex(uploadData.SchoolID),
		"grade":    uploadData.Grade,
		"class":    uploadData.ClassID,
		"valid":    true,
	}, bson.M{
		"$pull": bson.M{
			fieldName: bson.M{
				"$in": uploadData.Targets,
			},
		},
	})
	if err != nil {
		log.Printf("failed to delete targets, err %v\n", err)
		return err
	}

	return c.JSON(http.StatusOK, "successfully deleted")
}
