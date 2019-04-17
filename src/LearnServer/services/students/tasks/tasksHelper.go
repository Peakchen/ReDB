package tasks

import (
	"errors"
	"time"

	"LearnServer/models/userDB"
	"LearnServer/services/students/problempdfs"
	"gopkg.in/mgo.v2/bson"
)

// SaveTask 保存某个任务
func SaveTask(id string, taskType int, taskTimeUnix int64, problems []problempdfs.ProblemForCreatingFiles) error {
	err := userDB.C("students").UpdateId(bson.ObjectIdHex(id), bson.M{
		"$push": bson.M{
			"tasks": bson.M{
				"time":     time.Unix(taskTimeUnix, 0),
				"type":     taskType,
				"problems": problems,
			},
		},
	})
	return err
}

// DeleteTask 删除某个任务
func DeleteTask(id string, timeUnix int64) error {
	err := userDB.C("students").UpdateId(bson.ObjectIdHex(id), bson.M{
		"$pull": bson.M{
			"tasks": bson.M{
				"time": time.Unix(timeUnix, 0),
			},
		},
	})
	if err != nil && err.Error() == "not found" {
		// 删除不存在的任务是跳过但不报错
		return nil
	}
	return err
}

// GetTaskDetail 获取某个任务的内容
func GetTaskDetail(id string, timeUnix int64) ([]problempdfs.ProblemForCreatingFiles, error) {
	detail := struct {
		Task []struct {
			Problems []problempdfs.ProblemForCreatingFiles `bson:"problems"`
		} `bson:"tasks"`
	}{}

	err := userDB.C("students").Find(bson.M{
		"_id":        bson.ObjectIdHex(id),
		"tasks.time": time.Unix(timeUnix, 0),
	}).Select(bson.M{
		"tasks.$": 1,
	}).One(&detail)

	if err != nil {
		if err.Error() == "not found" {
			return nil, errors.New("can't find this upload task")
		}
		return nil, err
	}

	if len(detail.Task) == 0 {
		return nil, errors.New("can't find this upload task")
	}
	return detail.Task[0].Problems, nil
}
