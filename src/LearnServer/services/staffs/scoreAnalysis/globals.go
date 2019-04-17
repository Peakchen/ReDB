package scoreAnalysis

import "time"

type studentScoreType struct {
	LearnID int     `json:"learnID" bson:"learnID"` // 学生学习号
	Name    string  `json:"name" bson:"name"`       // 学生姓名
	Score   float32 `json:"score" bson:"score"`     // 成绩
}

type studentScoreSlice []studentScoreType

// 注意sort.Sort是将分数从高到低排序
func (c studentScoreSlice) Len() int {
	return len(c)
}
func (c studentScoreSlice) Swap(i int, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c studentScoreSlice) Less(i int, j int) bool {
	return c[i].Score > c[j].Score
}

type examType struct {
	Time    time.Time         `bson:"time"`    // 考试时间
	PaperID string            `bson:"paperID"` // 考试试卷ID
	Scores  studentScoreSlice `bson:"scores"`  //  考试成绩
}

type examSlice []examType

func (c examSlice) Len() int {
	return len(c)
}
func (c examSlice) Swap(i int, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c examSlice) Less(i int, j int) bool {
	return c[i].Time.Before(c[j].Time)
}
