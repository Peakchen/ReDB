package userDB

import (
	"gopkg.in/mgo.v2/bson"
)

// GetNotMarkedPaperIDs 获取一个学生未标记的试卷
func GetNotMarkedPaperIDs(id string) ([]string, error) {
	stu := struct {
		SchoolID bson.ObjectId `bson:"schoolID"`
		Grade    string        `bson:"grade"`
		Class    int           `bson:"classID"`
	}{}

	err := C("students").FindId(bson.ObjectIdHex(id)).One(&stu)
	if err != nil {
		return []string{}, err
	}

	classPaperIDs := struct {
		PaperIDs []string `bson:"papers"`
	}{}
	err = C("classes").Find(bson.M{
		"schoolID": stu.SchoolID,
		"grade":    stu.Grade,
		"class":    stu.Class,
		"valid":    true,
	}).One(&classPaperIDs)
	if err != nil {
		if err.Error() == "not found" {
			return []string{}, nil
		}
		return []string{}, err
	}

	// 去除已经标记的试卷
	paperIDsNotMarked := []string{}
	allProblemsDone := struct {
		Problems []struct {
			ProblemID  string `bson:"problemID"`
			SourceID   string `bson:"sourceID"`
			SourceType int    `bson:"sourceType"`
		} `bson:"problems"`
	}{}
	if err := GetAllProblemsDone(id, &allProblemsDone); err != nil {
		return []string{}, err
	}

	for _, paperID := range classPaperIDs.PaperIDs {
		// 因为录入试卷题目结果必定是整张试卷一起录入的，如果试卷中有一道题目已经录入，说明这张试卷已经标记过了
		found := false
		for _, p := range allProblemsDone.Problems {
			if p.SourceType == 2 && p.SourceID == paperID {
				found = true
				break
			}
		}

		if !found {
			paperIDsNotMarked = append(paperIDsNotMarked, paperID)
		}
	}

	return paperIDsNotMarked, nil
}

// GetMarkedPaperIDs 获取一个学生已经标记的试卷
func GetMarkedPaperIDs(id string) ([]string, error) {
	stu := struct {
		SchoolID bson.ObjectId `bson:"schoolID"`
		Grade    string        `bson:"grade"`
		Class    int           `bson:"classID"`
		Problems []struct {
			ProblemID  string `bson:"problemID"`
			SourceID   string `bson:"sourceID"`
			SourceType int    `bson:"sourceType"`
		} `bson:"problems"`
	}{}

	err := C("students").FindId(bson.ObjectIdHex(id)).One(&stu)
	if err != nil {
		return []string{}, err
	}

	classPaperIDs := struct {
		PaperIDs []string `bson:"papers"`
	}{}
	err = C("classes").Find(bson.M{
		"schoolID": stu.SchoolID,
		"grade":    stu.Grade,
		"class":    stu.Class,
		"valid":    true,
	}).One(&classPaperIDs)
	if err != nil {
		if err.Error() == "not found" {
			return []string{}, nil
		}
		return []string{}, err
	}

	// 已经标记的试卷
	paperIDsMarked := []string{}

	for _, paperID := range classPaperIDs.PaperIDs {
		// 因为录入试卷题目结果必定是整张试卷一起录入的，如果试卷中有一道题目已经录入，说明这张试卷已经标记过了
		found := false
		for _, p := range stu.Problems {
			if p.SourceType == 2 && p.SourceID == paperID {
				found = true
				break
			}
		}

		if found {
			paperIDsMarked = append(paperIDsMarked, paperID)
		}
	}

	return paperIDsMarked, nil
}
