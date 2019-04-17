package staffs

import (
	"fmt"
	"log"
	"net/http"
	"time"

	// "LearnServer/services/students/problempdfs"
	// "LearnServer/utils"

	"LearnServer/services/students/problempdfs"
	"LearnServer/utils"
	"github.com/labstack/echo"
)

type wrongProblemsResultType []problempdfs.ProblemForCreatingFiles

type batchDownloadInfo struct {
	BatchID                string    `json:"batchID"`
	CreateTime             int64     `json:"createTime"` // 该任务创建时间
	FinishTime             int64     `json:"finishTime"` // 完成时间或预计完成时间
	lastFileTime           time.Time // 上一份文档完成时间
	ProblemFilesFinished   int       `json:"problemFilesFinished"`   // 已处理的题目文件数目（无论是否成功）
	AnswerFilesFinished    int       `json:"answerFilesFinished"`    // 已处理的答案文件数目（无论是否成功）
	ProblemFilesSuccessful int       `json:"problemFilesSuccessful"` // 已完成并成功生成的题目文件数目
	AnswerFilesSuccessful  int       `json:"answerFilesSuccessful"`  // 已完成并成功生成的答案文件数目
	School                 string    `json:"school"`                 // 学校名称
	Grade                  string    `json:"grade"`                  // 年级
	Class                  int       `json:"class"`                  // 班级
	Students               []struct {
		Name              string                  `json:"name"`
		LearnID           int                     `json:"learnID"`
		ProblemFileStatus bool                    `json:"problemFileStatus"` // 题目文件生成状态(是否成功)
		AnswerFileStatus  bool                    `json:"answerFileStatus"`  // 答案文件生成状态(是否成功)
		ProblemStatusCode int                     `json:"problemStatusCode"` // 题目文件生成请求得到的状态码
		AnswerStatusCode  int                     `json:"answerStatusCode"`  // 答案文件生成请求得到的状态码
		Problems          wrongProblemsResultType `json:"problems"`          // 错题信息
	} `json:"students"` // 需要生成纠错本的学生的信息
}

var batchDownloadMap = make(map[string]batchDownloadInfo)

func createBatchDownloadTask(c echo.Context) error {
	// 创建新的批量下载文档任务
	batchDownloadTask := batchDownloadInfo{}

	if err := c.Bind(&batchDownloadTask); err != nil {
		return utils.InvalidParams("invalid inputs, error: " + err.Error())
	}

	batchDownloadTask.CreateTime = time.Now().Unix()
	for i := range batchDownloadTask.Students {
		batchDownloadTask.Students[i].ProblemFileStatus = false
		batchDownloadTask.Students[i].AnswerFileStatus = false
		batchDownloadTask.Students[i].ProblemStatusCode = 0
		batchDownloadTask.Students[i].AnswerStatusCode = 0
	}

	uuid, err := utils.UUID()
	if err != nil {
		log.Printf("create uuid failed, err: %v\n", err)
		return err
	}
	batchDownloadTask.BatchID = uuid
	batchDownloadMap[uuid] = batchDownloadTask

	return c.JSON(http.StatusOK, echo.Map{
		"batchID": uuid,
	})
}

func deleteBatchDownloadTask(c echo.Context) error {
	// 删除批量下载文档任务记录，但是不影响后台继续生成文档
	batchID := c.Param("batchID")
	delete(batchDownloadMap, batchID)
	return c.JSON(http.StatusOK, "successfully deleted batchID")
}

func getStudentIndex(batchID string, learnID int) (int, error) {
	// 获取该 learnID 在 batchID 对应的记录中的 students 数组中的下标
	batchInfo, ok := batchDownloadMap[batchID]
	if !ok {
		return -1, fmt.Errorf("this batchID doesn't exist")
	}

	for index, stu := range batchInfo.Students {
		if stu.LearnID == learnID {
			return index, nil
		}
	}

	return -1, fmt.Errorf("can not find this student in this batch download")
}

func listBatchDownloadTasks(c echo.Context) error {

	results := []batchDownloadInfo{}

	for _, batchInfo := range batchDownloadMap {
		problemFilesSuccessful := 0
		answerFilesSuccessful := 0
		problemFilesFinished := 0
		answerFilesFinished := 0
		totalFiles := len(batchInfo.Students) * 2
		for _, stu := range batchInfo.Students {
			if stu.ProblemFileStatus {
				problemFilesSuccessful++
			}
			if stu.AnswerFileStatus {
				answerFilesSuccessful++
			}
			if stu.ProblemStatusCode != 0 {
				problemFilesFinished++
			}
			if stu.AnswerStatusCode != 0 {
				answerFilesFinished++
			}
		}

		batchInfo.ProblemFilesSuccessful = problemFilesSuccessful
		batchInfo.AnswerFilesSuccessful = answerFilesSuccessful
		batchInfo.ProblemFilesFinished = problemFilesFinished
		batchInfo.AnswerFilesFinished = answerFilesFinished

		const timePerFile = 30
		filesLeft := totalFiles - problemFilesFinished - answerFilesFinished
		if filesLeft == 0 {
			batchInfo.FinishTime = batchInfo.lastFileTime.Unix()
		} else {
			batchInfo.FinishTime = time.Now().Add(time.Duration(filesLeft*timePerFile) * time.Second).Unix()
		}

		results = append(results, batchInfo)
	}

	return c.JSON(http.StatusOK, results)
}
