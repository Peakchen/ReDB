package students

import (
	"log"
	"net/http"
	"time"

	"LearnServer/models/contentDB"
	"LearnServer/models/userDB"
	"LearnServer/services/students/validation"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

// BookStatusType 资料录入状态类型
type BookStatusType struct {
	Book   string `json:"book" db:"name"`
	Status int    `json:"status"` // 0最近一周有标记，1没有
}

// GetProblemRecordBookStatus 获取资料的录入状态
func GetProblemRecordBookStatus(id string) ([]BookStatusType, error) {
	type problemDBType struct {
		ProblemID  string    `bson:"problemID"`
		Time       time.Time `bson:"assignDate"`
		SourceID   string    `bson:"sourceID"`
		SourceType int       `bson:"sourceType"`
	}
	stu := struct {
		SchoolID bson.ObjectId   `bson:"schoolID"`
		Grade    string          `bson:"grade"`
		Class    int             `bson:"classID"`
		Problems []problemDBType `bson:"problems"`
	}{}

	err := userDB.C("students").FindId(bson.ObjectIdHex(id)).Select(bson.M{
		"schoolID": 1,
		"grade":    1,
		"classID":  1,
		"problems": 1,
	}).One(&stu)
	if err != nil {
		log.Printf("failed to get student info, id %s, error: %v", id, err)
		return []BookStatusType{}, err
	}

	classBookIDs := struct {
		BookIDs []string `bson:"books"`
	}{}
	err = userDB.C("classes").Find(bson.M{
		"schoolID": stu.SchoolID,
		"grade":    stu.Grade,
		"class":    stu.Class,
		"valid":    true,
	}).One(&classBookIDs)
	if err != nil {
		if err.Error() == "not found" {
			return []BookStatusType{}, nil
		}
		log.Printf("failed to get books of class , err %v", err)
		return []BookStatusType{}, err
	}

	bookIDStatusMap := make(map[string]bool)
	for _, p := range stu.Problems {
		if p.Time.Before(time.Now()) && time.Now().AddDate(0, 0, -7).Before(p.Time) && p.SourceType == 1 {
			bookIDStatusMap[p.SourceID] = true
		}
	}

	bookStatuses := []BookStatusType{}

	for _, bookID := range classBookIDs.BookIDs {
		bookStatus := BookStatusType{}
		err := contentDB.GetDB().Get(&bookStatus, "SELECT name from books where bookID = ?;", bookID)
		if err != nil {
			log.Printf("failed to get name of bookID %s, err: %v", bookID, err)
			continue
		}

		status, found := bookIDStatusMap[bookID]
		if !found || !status {
			bookStatus.Status = 1
		} else {
			bookStatus.Status = 0
		}
		bookStatuses = append(bookStatuses, bookStatus)
	}

	return bookStatuses, nil
}

// TestHasMarkTasks 判断是否有未标记的纠错本
func TestHasMarkTasks(id string) bool {
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

// TestHasNotMarkedPapers 判断是否有未标记的试卷
func TestHasNotMarkedPapers(id string) bool {
	paperIDs, err := userDB.GetNotMarkedPaperIDs(id)
	return err == nil && len(paperIDs) > 0
}

func getProblemRecordsHandler(c echo.Context) error {
	// 获取题目录入记录信息
	var id string
	if err := validation.ValidateUser(c, &id); err != nil {
		return err
	}

	result := struct {
		WrongProblemStatus int              `json:"wrongProblemStatus"` // 纠错本状态，1未标记，0已标记
		PaperStatus        int              `json:"paperStatus"`        // 试卷状态，1未标记，0已标记
		BookStatus         []BookStatusType `json:"bookStatus"`
	}{}

	bookStatus, err := GetProblemRecordBookStatus(id)
	if err != nil {
		return err
	}

	result.BookStatus = bookStatus

	if TestHasNotMarkedPapers(id) {
		result.PaperStatus = 1
	} else {
		result.PaperStatus = 0
	}

	if TestHasMarkTasks(id) {
		result.WrongProblemStatus = 1
	} else {
		result.WrongProblemStatus = 0
	}

	return c.JSON(http.StatusOK, result)
}
