package staffs

import (
	"log"
	"net/http"

	"LearnServer/services/students"
	"LearnServer/utils"
	"github.com/labstack/echo"
)

func getProblemRecordsHandler(c echo.Context) error {
	// 获取题目录入记录信息
	learnIDs := []int{}

	if err := c.Bind(&learnIDs); err != nil {
		return utils.InvalidParams("invalid inputs, err: " + err.Error())
	}

	type recordType struct {
		LearnID            int                       `json:"learnID"`
		WrongProblemStatus int                       `json:"wrongProblemStatus"` // 纠错本状态，1未标记，0已标记
		PaperStatus        int                       `json:"paperStatus"`        // 试卷状态，1未标记，0已标记
		BookStatus         []students.BookStatusType `json:"bookStatus"`
	}

	result := []recordType{}

	for _, learnID := range learnIDs {
		studentID, err := getStudentIDByLearnID(learnID)
		if err != nil {
			log.Printf("can not find the information of learnID %d, err %v\n", learnID, err)
			continue
		}

		record := recordType{}
		bookStatus, err := students.GetProblemRecordBookStatus(studentID)
		if err != nil {
			return err
		}

		record.BookStatus = bookStatus

		if students.TestHasNotMarkedPapers(studentID) {
			record.PaperStatus = 1
		} else {
			record.PaperStatus = 0
		}

		if students.TestHasMarkTasks(studentID) {
			record.WrongProblemStatus = 1
		} else {
			record.WrongProblemStatus = 0
		}

		record.LearnID = learnID

		result = append(result, record)
	}

	return c.JSON(http.StatusOK, result)
}
