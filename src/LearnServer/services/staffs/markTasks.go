package staffs

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"LearnServer/models/userDB"
	
	"LearnServer/services/students/problempdfs"
	"LearnServer/services/students/tasks"
	"LearnServer/utils"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

func createMarkTasksHandler(c echo.Context) error {
	type taskType struct {
		Time     int64                                 `json:"time"`
		Type     int                                   `json:"type"`
		LearnID  int                                   `json:"learnID"`
		Problems []problempdfs.ProblemForCreatingFiles `json:"problems"`
	}

	var markTasks []taskType
	if err := c.Bind(&markTasks); err != nil {
		return utils.InvalidParams("invalid input!")
	}

	failedIDs := []int{}
	for _, task := range markTasks {
		id, err := getStudentIDByLearnID(task.LearnID)
		if err != nil {
			log.Printf("can't find the student with learnID %d. \n", task.LearnID)
			failedIDs = append(failedIDs, task.LearnID)
			continue
		}

		if err := tasks.SaveTask(id, task.Type, task.Time, task.Problems); err != nil {
			failedIDs = append(failedIDs, task.LearnID)
			log.Printf("saving task for learnID %d failed, task: %v \n", task.LearnID, task)
			continue
		}
	}

	if len(failedIDs) != 0 {
		return c.JSON(http.StatusNotFound, failedIDs)
	}
	return c.JSON(http.StatusOK, "Successfully create a task.")
}

type taskInfoType struct {
	Time int64 `json:"time"`
	Type int   `json:"type"`
}

func getMarkTaskOfStudent(id bson.ObjectId) ([]taskInfoType, error) {
	// 获取某个学生的标记任务
	type taskInfoDBType struct {
		Time time.Time `bson:"time"`
		Type int       `bson:"type"`
	}

	tasksInfoDB := struct {
		Tasks []taskInfoDBType `bson:"tasks"`
	}{}
	err := userDB.C("students").FindId(id).Select(bson.M{
		"tasks": 1,
	}).One(&tasksInfoDB)
	if err != nil {
		return nil, err
	}

	result := make([]taskInfoType, len(tasksInfoDB.Tasks))
	for i, t := range tasksInfoDB.Tasks {
		result[i] = taskInfoType{
			Time: t.Time.Unix(),
			Type: t.Type,
		}
	}

	return result, nil
}

// listMarkTasksHandler 列出某个学生所有未完成的标记任务
func listMarkTasksHandler(c echo.Context) error {

	learnID, err := strconv.Atoi(c.Param("learnID"))
	if err != nil {
		return utils.InvalidParams("learnID is invalid")
	}
	studentID, err := getStudentIDByLearnID(learnID)
	if err != nil {
		return utils.NotFound("can not find the information of this learnID")
	}

	result, err := getMarkTaskOfStudent(bson.ObjectIdHex(studentID))
	if err != nil {
		log.Printf("fail to get upload task of student %d, err %v\n", learnID, err)
		return err
	}

	// if len(result) == 0 {
	//	return utils.NotFound("No upload tasks.")
	// }

	return c.JSON(http.StatusOK, result)
}

// getMarkTaskHandler 获取某个标记任务具体内容
func getMarkTaskHandler(c echo.Context) error {
	learnID, err := strconv.Atoi(c.Param("learnID"))
	if err != nil {
		return utils.InvalidParams("learnID is invalid")
	}
	studentID, err := getStudentIDByLearnID(learnID)
	if err != nil {
		return utils.NotFound("can not find the information of this learnID")
	}

	taskTimeUnix, err := strconv.ParseInt(c.Param("time"), 10, 64)
	if err != nil {
		return utils.InvalidParams("path parameter time is invalid")
	}

	problems, err := tasks.GetTaskDetail(studentID, taskTimeUnix)
	if err != nil {
		if err.Error() == "can't find this upload task" {
			return utils.NotFound("Can not find the task.")
		}
		return err
	}

	return c.JSON(http.StatusOK, problems)
}

func deleteMarkTaskOfStudent(learnID int, timeUnix int64) error {
	// 删除某个学生的特定标记任务
	studentID, err := getStudentIDByLearnID(learnID)
	if err != nil {
		return utils.NotFound("can not find the information of this learnID")
	}

	err = tasks.DeleteTask(studentID, timeUnix)
	if err != nil {
		return err
	}

	return nil
}

// deleteMarkTaskHandler 删除学生标记任务接口
func deleteMarkTaskHandler(c echo.Context) error {
	learnID, err := strconv.Atoi(c.Param("learnID"))
	if err != nil {
		return utils.InvalidParams("learnID is invalid")
	}

	taskTimeUnix, err := strconv.ParseInt(c.Param("time"), 10, 64)
	if err != nil {
		return utils.InvalidParams("path parameter time is invalid")
	}

	if err := deleteMarkTaskOfStudent(learnID, taskTimeUnix); err != nil {
		log.Printf("failed to delete upload task of %d, err %v\n", learnID, err)
		return err
	}

	return c.JSON(http.StatusOK, "Successfully deleted")
}

// listClassMarkTasksHandler 获取班级所有学生的标记任务
func listClassMarkTasksHandler(c echo.Context) error {


	schoolID := c.QueryParam("schoolID")
	grade := c.QueryParam("grade")
	classID, err := strconv.Atoi(c.QueryParam("class"))
	if err != nil {
		return utils.InvalidParams("class is not a number!")
	}

	students, err := userDB.GetStudents(schoolID, grade, classID, 0, "", "", "", "")
	if err != nil {
		if err.Error() == "can't find students of this school and class" {
			return utils.NotFound("can't find students of this school and class")
		}
		log.Printf("failed to get students of schoolID %s, grade %s, classID %d, err %v \n", schoolID, grade, classID, err)
		return err
	}

	type taskType struct {
		LearnID int64  `json:"learnID"` // 学生学习号
		Name    string `json:"name"`    // 学生名字
		Time    int64  `json:"time"`    // 任务发生的时间对应的unix时间戳
		Type    int    `json:"type"`    // 任务类型（没标记的是错题为1，没标记的是检验题为2）
	}
	tasks := []taskType{}
	for _, stu := range students {
		tasksOfStu, err := getMarkTaskOfStudent(stu.ID)
		if err != nil {
			log.Printf("fail to get upload task of student %d, err %v\n", stu.LearnID, err)
			continue
		}
		for _, t := range tasksOfStu {
			tasks = append(tasks, taskType{
				LearnID: stu.LearnID,
				Name:    stu.Name,
				Time:    t.Time,
				Type:    t.Type,
			})
		}
	}

	if len(tasks) == 0 {
		return utils.NotFound("No upload tasks.")
	}

	return c.JSON(http.StatusOK, tasks)
}

// deleteBundleMarkTasksHandler 批量删除标记任务
func deleteBundleMarkTasksHandler(c echo.Context) error {


	uploadData := []struct {
		LearnID int   `json:"learnID"` // 学生学习号
		Time    int64 `json:"time"`    // 任务发生的时间对应的unix时间戳
	}{}

	if err := c.Bind(&uploadData); err != nil {
		return utils.InvalidParams("invalid input!" + err.Error())
	}

	for _, task := range uploadData {
		if err := deleteMarkTaskOfStudent(task.LearnID, task.Time); err != nil {
			log.Printf("failed to delete upload task of %d, err %v\n", task.LearnID, err)
			return err
		}
	}

	return c.JSON(http.StatusOK, "Successfully deleted")
}
