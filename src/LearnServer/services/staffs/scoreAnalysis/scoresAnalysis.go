package scoreAnalysis

import (
	"math"
	"sort"

	"LearnServer/models/userDB"
	"gopkg.in/mgo.v2/bson"
)

// calAverageOfExam 计算考试平均分
func calAverageOfExam(exam examType) float32 {
	if len(exam.Scores) <= 0 {
		return 0
	}
	var count float32
	for _, stu := range exam.Scores {
		count += stu.Score
	}
	return count / float32(len(exam.Scores))
}

// 获取考试前N名
func getExamTopN(exam examType, n int) studentScoreSlice {
	if len(exam.Scores) <= n {
		result := make(studentScoreSlice, len(exam.Scores))
		copy(result, exam.Scores)
		return result
	}
	scores := exam.Scores
	sort.Sort(scores)
	result := make(studentScoreSlice, n)
	copy(result, scores[:n])
	return result
}

// 获取考试后N名
func getExamLastN(exam examType, n int) studentScoreSlice {
	if len(exam.Scores) <= n {
		result := make(studentScoreSlice, len(exam.Scores))
		copy(result, exam.Scores)
		return result
	}
	scores := exam.Scores
	sort.Sort(sort.Reverse(scores))
	result := make(studentScoreSlice, n)
	copy(result, scores[:n])
	return result
}

// calLevelAverageOfExam 计算排名层均分, levelInterval 每个level包含的名次数目
func calLevelAverageOfExam(exam examType, levelInterval int) []float32 {
	scores := exam.Scores
	sort.Sort(scores)
	if levelInterval <= 0 {
		levelInterval = len(scores)
	}
	stuLeft := len(scores)
	averages := []float32{}
	for stuLeft > 0 {
		var scoreCount float32
		stuNum := 0
		for stuNum < levelInterval && stuLeft > 0 {
			scoreCount += scores[len(scores)-stuLeft].Score
			stuNum++
			stuLeft--
		}
		averages = append(averages, scoreCount/float32(stuNum))
	}
	return averages
}

// calScoreProportionOfExam 计算分数段占比
func calScoreProportionOfExam(exam examType, standardFullScore int) []float32 {
	scores := exam.Scores
	sort.Sort(sort.Reverse(scores))
	// 分数段节点， MaxFloat32 两个节点设置方便判断
	nodes := []float32{-math.MaxFloat32, 60, 70, 80, 90, 95, math.MaxFloat32}
	for i := range nodes {
		if i != 0 && i != len(nodes)-1 {
			nodes[i] = nodes[i] / float32(100) * float32(standardFullScore)
		}
	}
	proportions := make([]float32, len(nodes)-1)
	stuCounts := 0
	currNodeIndex := 0
	for _, stu := range scores {
		nodeIndexChanged := false
		for !(stu.Score >= nodes[currNodeIndex] && stu.Score < nodes[currNodeIndex+1]) {
			if !nodeIndexChanged {
				// 即将修改 nodeIndex，先将之前统计的数据写入
				proportions[currNodeIndex] = float32(stuCounts) / float32(len(scores))
				stuCounts = 0
			}
			currNodeIndex++
			nodeIndexChanged = true
		}
		stuCounts++
	}
	// 写入最后一个分数段的统计
	proportions[currNodeIndex] = float32(stuCounts) / float32(len(scores))
	return proportions
}

// scoreStandardization 把 exam 由 currFullScore 统一转换到 standardFullScore
func scoreStandardization(exam examType, currFullScore int, standardFullScore int) examType {
	scores := exam.Scores
	for i, stu := range scores {
		stu.Score = stu.Score / float32(currFullScore) * float32(standardFullScore)
		scores[i] = stu
	}
	exam.Scores = scores
	return exam
}

func getStudentLevel(learnID int) (int, error) {
	stuLevel := struct {
		Level int `db:"level"` // 层级
	}{}
	err := userDB.C("students").Find(bson.M{
		"learnID": learnID,
		"valid":   true,
	}).Select(bson.M{
		"level": 1,
	}).One(&stuLevel)
	return stuLevel.Level, err
}

// getStudentsLevelRanking 获取学生在该层级的排名
func getStudentsLevelRanking(exam examType, studentLevelMap map[int]int) map[int]int {
	scores := exam.Scores
	sort.Sort(scores)
	stuRankingMap := make(map[int]int)
	// 存放 level 对应的 最大名次
	levelRankingMaxMap := make(map[int]int)
	for _, stu := range scores {
		level := studentLevelMap[stu.LearnID]
		if rankingMax, ok := levelRankingMaxMap[level]; ok {
			stuRankingMap[stu.LearnID] = rankingMax + 1
			levelRankingMaxMap[level] = rankingMax + 1
		} else {
			stuRankingMap[stu.LearnID] = 1
			levelRankingMaxMap[level] = 1
		}
	}
	return stuRankingMap
}
