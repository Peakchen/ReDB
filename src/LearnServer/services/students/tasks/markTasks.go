package tasks

import (
	"net/http"
	"strconv"
	"time"

	"LearnServer/models/userDB"
	"LearnServer/services/students/problempdfs"
	"LearnServer/services/students/validation"
	"LearnServer/utils"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

// ListMarkTasksHandler 列出所有未完成的标记任务
func ListMarkTasksHandler(c echo.Context) error {
	var id string
	if err := validation.ValidateUser(c, &id); err != nil {
		return err
	}

	type taskInfoDBType struct {
		Time time.Time `bson:"time"`
		Type int       `bson:"type"`
	}

	tasksInfoDB := struct {
		Tasks []taskInfoDBType `bson:"tasks"`
	}{}
	err := userDB.C("students").FindId(bson.ObjectIdHex(id)).Select(bson.M{
		"tasks": 1,
	}).One(&tasksInfoDB)
	if err != nil {
		return err
	}

	type taskInfoType struct {
		Time int64 `json:"time"`
		Type int   `json:"type"`
	}

	result := make([]taskInfoType, len(tasksInfoDB.Tasks))
	for i, t := range tasksInfoDB.Tasks {
		result[i] = taskInfoType{
			Time: t.Time.Unix(),
			Type: t.Type,
		}
	}

	// if len(result) == 0 {
	//	return utils.NotFound("No upload tasks.")
	// }

	return c.JSON(http.StatusOK, result)
}

// GetMarkTaskHandler 获取某个标记任务具体内容
func GetMarkTaskHandler(c echo.Context) error {
	var id string
	if err := validation.ValidateUser(c, &id); err != nil {
		return err
	}

	taskTimeUnix, err := strconv.ParseInt(c.Param("time"), 10, 64)
	if err != nil {
		return utils.InvalidParams("path parameter time is invalid")
	}

	problems, err := GetTaskDetail(id, taskTimeUnix)
	if err != nil {
		if err.Error() == "can't find this upload task" {
			return utils.NotFound("Can not find the task.")
		}
		return err
	}

	return c.JSON(http.StatusOK, problems)
}

// DeleteMarkTaskHandler 删除某个标记任务
func DeleteMarkTaskHandler(c echo.Context) error {
	var id string
	if err := validation.ValidateUser(c, &id); err != nil {
		return err
	}

	taskTimeUnix, err := strconv.ParseInt(c.Param("time"), 10, 64)
	if err != nil {
		return utils.InvalidParams("path parameter time is invalid")
	}

	err = DeleteTask(id, taskTimeUnix)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, "Successfully deleted")
}

// CreateMarkTaskHandler 创建一个标记任务
func CreateMarkTaskHandler(c echo.Context) error {
	var id string
	if err := validation.ValidateUser(c, &id); err != nil {
		return err
	}

	type taskType struct {
		Time     int64                                 `json:"time"`
		Type     int                                   `json:"type"`
		Problems []problempdfs.ProblemForCreatingFiles `json:"problems"`
	}

	var markTask taskType
	if err := c.Bind(&markTask); err != nil {
		return utils.InvalidParams()
	}

	if err := SaveTask(id, markTask.Type, markTask.Time, markTask.Problems); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "Successfully create a task.")
}
