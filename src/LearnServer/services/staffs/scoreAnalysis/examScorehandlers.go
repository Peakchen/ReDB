package scoreAnalysis

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"LearnServer/models/contentDB"
	"LearnServer/models/userDB"
	"LearnServer/utils"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

func getExams(schoolID string, grade string, class int, startTime time.Time, endTime time.Time, standardFullScore int) (examSlice, error) {
	// 获取班级特定时间之间的考试结果
	type examScoreType struct {
		Time   time.Time         `bson:"time"`
		Scores studentScoreSlice `bson:"scores"`
	}

	classExamRecord := make(map[string]map[string]examScoreType)

	err := userDB.C("classes").Find(bson.M{
		"schoolID": bson.ObjectIdHex(schoolID),
		"grade":    grade,
		"class":    class,
		"valid":    true,
	}).Select(bson.M{
		"examScoreRecords": 1,
	}).One(&classExamRecord)
	if err != nil {
		log.Printf("getting class examRecord failed")
		return nil, err
	}

	examRecords, ok := classExamRecord["examScoreRecords"]
	if !ok {
		return nil, fmt.Errorf("failed to get examRecords")
	}

	db := contentDB.GetDB()
	exams := examSlice{}
	for key, value := range examRecords {
		if value.Time.Before(endTime) && startTime.Before(value.Time) {
			paperID := strings.TrimPrefix(key, "paperID")
			fullScore := 100
			err := db.Get(&fullScore, "SELECT fullScore FROM papers WHERE paperID = ?;", paperID)
			if err != nil {
				log.Printf("failed to get fullScore of paper %s, err %v\n", paperID, err)
			}
			scores := studentScoreSlice{}
			for _, stu := range value.Scores {
				if stu.Score != -1 {
					// 去除缺考学生
					scores = append(scores, stu)
				}
			}
			exam := scoreStandardization(examType{
				Time:    value.Time,
				PaperID: paperID,
				Scores:  scores,
			}, fullScore, standardFullScore)
			exams = append(exams, exam)
		}
	}
	sort.Sort(exams)
	return exams, nil
}

// GetAverageAnalysisHandler 获取某个班级的均分分析
func GetAverageAnalysisHandler(c echo.Context) error {
	schoolID := c.QueryParam("schoolID")
	grade := c.QueryParam("grade")
	class, err := strconv.Atoi(c.QueryParam("class"))
	if err != nil {
		return utils.InvalidParams("class is not a number!")
	}
	startTimeInt64, err := strconv.ParseInt(c.QueryParam("startTime"), 10, 64)
	if err != nil {
		return utils.InvalidParams("startTime is not a number!")
	}
	startTime := time.Unix(startTimeInt64, 0)
	endTimeInt64, err := strconv.ParseInt(c.QueryParam("endTime"), 10, 64)
	if err != nil {
		return utils.InvalidParams("endTime is not a number!")
	}
	endTime := time.Unix(endTimeInt64, 0)
	standardFullScore, err := strconv.Atoi(c.QueryParam("standardFullScore"))
	if err != nil {
		return utils.InvalidParams("standardFullScore is not a number!")
	}

	type resultExamType struct {
		Time         int64   `json:"time"`         // 考试时间，unix时间戳（注意因为界面只选择了日期，用0时0分0秒）
		PaperID      string  `json:"paperID"`      // 考试试卷ID
		AverageScore float32 `json:"averageScore"` // 均分
	}
	results := struct {
		LatestTop10  studentScoreSlice `json:"latestTop10"`  // 最新考试前10名（右侧展示）
		LatestLast10 studentScoreSlice `json:"latestLast10"` // 最新考试后10名（右侧展示）
		Exams        []resultExamType  `json:"exams"`
	}{}

	exams, err := getExams(schoolID, grade, class, startTime, endTime, standardFullScore)
	if err != nil || len(exams) <= 0 {
		return c.JSON(http.StatusOK, results)
	}

	results.Exams = make([]resultExamType, len(exams))
	for i, e := range exams {
		results.Exams[i] = resultExamType{
			Time:         e.Time.Unix(),
			PaperID:      e.PaperID,
			AverageScore: calAverageOfExam(e),
		}
	}
	results.LatestTop10 = getExamTopN(exams[len(exams)-1], 10)
	results.LatestLast10 = getExamLastN(exams[len(exams)-1], 10)

	return c.JSON(http.StatusOK, results)
}

// GetRankingLevelAverageAnalysisHandler 获取某个班级的分层均分分析
func GetRankingLevelAverageAnalysisHandler(c echo.Context) error {
	schoolID := c.QueryParam("schoolID")
	grade := c.QueryParam("grade")
	class, err := strconv.Atoi(c.QueryParam("class"))
	if err != nil {
		return utils.InvalidParams("class is not a number!")
	}
	startTimeInt64, err := strconv.ParseInt(c.QueryParam("startTime"), 10, 64)
	if err != nil {
		return utils.InvalidParams("startTime is not a number!")
	}
	startTime := time.Unix(startTimeInt64, 0)
	endTimeInt64, err := strconv.ParseInt(c.QueryParam("endTime"), 10, 64)
	if err != nil {
		return utils.InvalidParams("endTime is not a number!")
	}
	endTime := time.Unix(endTimeInt64, 0)
	standardFullScore, err := strconv.Atoi(c.QueryParam("standardFullScore"))
	if err != nil {
		return utils.InvalidParams("standardFullScore is not a number!")
	}

	type examScoreType struct {
		Time         int64   `json:"time"`         // 考试时间，unix时间戳（注意因为界面只选择了日期，用0时0分0秒）
		PaperID      string  `json:"paperID"`      // 考试试卷ID
		AverageScore float32 `json:"averageScore"` // 某一层均分
	}
	type levelScoreType struct {
		LatestRankinglevel string          `json:"latestRankinglevel"` // 最新考试排名段
		Data               []examScoreType `json:"data"`
	}
	results := []levelScoreType{}
	extendResults := func(newSize int) {
		for len(results) < newSize {
			results = append(results, levelScoreType{})
		}
	}

	exams, err := getExams(schoolID, grade, class, startTime, endTime, standardFullScore)
	if err != nil || len(exams) <= 0 {
		return c.JSON(http.StatusOK, results)
	}

	const rankingLevels int = 5
	for _, e := range exams {
		levelInterval := len(e.Scores) / rankingLevels
		if len(e.Scores)%rankingLevels != 0 {
			levelInterval++
		}
		levelAverages := calLevelAverageOfExam(e, levelInterval)
		extendResults(len(levelAverages))
		startRanking := 1
		for i, average := range levelAverages {
			if len(results[i].Data) <= 0 {
				// 还没有该排名段的数据，初始化
				results[i] = levelScoreType{
					LatestRankinglevel: fmt.Sprintf("%d~%d名", startRanking, startRanking+levelInterval-1),
					Data:               []examScoreType{},
				}
			}
			// 更新 LatestRankinglevel 使得排名段信息为最新的
			results[i].LatestRankinglevel = fmt.Sprintf("%d~%d名", startRanking, startRanking+levelInterval-1)
			results[i].Data = append(results[i].Data, examScoreType{
				Time:         e.Time.Unix(),
				PaperID:      e.PaperID,
				AverageScore: average,
			})
			startRanking += levelInterval
		}
	}

	return c.JSON(http.StatusOK, results)
}

// GetScoreProportionAnalysisHandler 获取某个班级的分数段占比分析
func GetScoreProportionAnalysisHandler(c echo.Context) error {
	schoolID := c.QueryParam("schoolID")
	grade := c.QueryParam("grade")
	class, err := strconv.Atoi(c.QueryParam("class"))
	if err != nil {
		return utils.InvalidParams("class is not a number!")
	}
	startTimeInt64, err := strconv.ParseInt(c.QueryParam("startTime"), 10, 64)
	if err != nil {
		return utils.InvalidParams("startTime is not a number!")
	}
	startTime := time.Unix(startTimeInt64, 0)
	endTimeInt64, err := strconv.ParseInt(c.QueryParam("endTime"), 10, 64)
	if err != nil {
		return utils.InvalidParams("endTime is not a number!")
	}
	endTime := time.Unix(endTimeInt64, 0)
	standardFullScore, err := strconv.Atoi(c.QueryParam("standardFullScore"))
	if err != nil {
		return utils.InvalidParams("standardFullScore is not a number!")
	}

	type examScoreType struct {
		Time    int64   `json:"time"`    // 考试时间，unix时间戳（注意因为界面只选择了日期，用0时0分0秒）
		PaperID string  `json:"paperID"` // 考试试卷ID
		Rate    float32 `json:"rate"`    // 该分数层这次分数段占比
	}
	type levelScoreType struct {
		ScoreSegment int             `json:"scoreSegment"` // 分数段，1 为 60-， 2 为 60-70 ，以下类推
		Data         []examScoreType `json:"data"`
	}
	results := []levelScoreType{}
	extendResults := func(newSize int) {
		for len(results) < newSize {
			results = append(results, levelScoreType{})
		}
	}

	exams, err := getExams(schoolID, grade, class, startTime, endTime, standardFullScore)
	if err != nil || len(exams) <= 0 {
		return c.JSON(http.StatusOK, results)
	}

	for _, e := range exams {
		proportions := calScoreProportionOfExam(e, standardFullScore)
		extendResults(len(proportions))
		for i, rate := range proportions {
			if len(results[i].Data) <= 0 {
				// 还没有该分数段的数据，初始化
				results[i] = levelScoreType{
					ScoreSegment: i + 1,
					Data:         []examScoreType{},
				}
			}
			results[i].Data = append(results[i].Data, examScoreType{
				Time:    e.Time.Unix(),
				PaperID: e.PaperID,
				Rate:    rate,
			})
		}
	}

	return c.JSON(http.StatusOK, results)
}

// GetStudentScoreAnalysisHandler 获取某个班级的个人分数分析
func GetStudentScoreAnalysisHandler(c echo.Context) error {
	schoolID := c.QueryParam("schoolID")
	grade := c.QueryParam("grade")
	class, err := strconv.Atoi(c.QueryParam("class"))
	if err != nil {
		return utils.InvalidParams("class is not a number!")
	}
	startTimeInt64, err := strconv.ParseInt(c.QueryParam("startTime"), 10, 64)
	if err != nil {
		return utils.InvalidParams("startTime is not a number!")
	}
	startTime := time.Unix(startTimeInt64, 0)
	endTimeInt64, err := strconv.ParseInt(c.QueryParam("endTime"), 10, 64)
	if err != nil {
		return utils.InvalidParams("endTime is not a number!")
	}
	endTime := time.Unix(endTimeInt64, 0)
	standardFullScore, err := strconv.Atoi(c.QueryParam("standardFullScore"))
	if err != nil {
		return utils.InvalidParams("standardFullScore is not a number!")
	}

	type examScoreType struct {
		Time    int64   `json:"time"`    // 考试时间，unix时间戳（注意因为界面只选择了日期，用0时0分0秒）
		PaperID string  `json:"paperID"` // 考试试卷ID
		Score   float32 `json:"score"`   // 这次考试分数
	}
	type studentScoreType struct {
		LearnID       int             `json:"learnID"`       // 学习号
		Name          string          `json:"name"`          // 学生名字
		LatestRanking int             `json:"latestRanking"` // 最新排名（右侧表格展示）
		Data          []examScoreType `json:"data"`
	}
	type levelScoreType struct {
		Level    int                `json:"level"`
		Students []studentScoreType `json:"students"`
	}
	results := []levelScoreType{}

	exams, err := getExams(schoolID, grade, class, startTime, endTime, standardFullScore)
	if err != nil || len(exams) <= 0 {
		return c.JSON(http.StatusOK, results)
	}

	studentsMap := make(map[int]studentScoreType)
	for _, exam := range exams {
		for _, stu := range exam.Scores {
			examScore := examScoreType{
				Time:    exam.Time.Unix(),
				PaperID: exam.PaperID,
				Score:   stu.Score,
			}
			if studentInMap, ok := studentsMap[stu.LearnID]; !ok {
				studentsMap[stu.LearnID] = studentScoreType{
					LearnID: stu.LearnID,
					Name:    stu.Name,
					Data: []examScoreType{
						examScore,
					},
				}
			} else {
				studentInMap.Data = append(studentInMap.Data, examScore)
				studentsMap[stu.LearnID] = studentInMap
			}
		}
	}

	latestScores := exams[len(exams)-1].Scores
	sort.Sort(latestScores)
	lastRanking := 1
	var lastScore float32 = -100
	for index, stu := range latestScores {
		studentInMap := studentsMap[stu.LearnID]
		if lastScore != stu.Score {
			studentInMap.LatestRanking = index + 1
			lastRanking = index + 1
			lastScore = stu.Score
		} else {
			// 成绩相同则排名相同
			studentInMap.LatestRanking = lastRanking
		}
		studentsMap[stu.LearnID] = studentInMap
	}

	// 对学生分层整理
	levelMap := make(map[int]levelScoreType)
	for _, stu := range studentsMap {
		level, err := getStudentLevel(stu.LearnID)
		if err != nil || level == 0 || level == -1 {
			return utils.Forbidden("some students do not have levels")
		}
		if levelInMap, ok := levelMap[level]; !ok {
			levelMap[level] = levelScoreType{
				Level: level,
				Students: []studentScoreType{
					stu,
				},
			}
		} else {
			levelInMap.Students = append(levelInMap.Students, stu)
			levelMap[level] = levelInMap
		}
	}

	// 转为list
	for _, value := range levelMap {
		results = append(results, value)
	}
	sort.Slice(results, func(i, j int) bool { return results[i].Level < results[j].Level })

	for index, value := range results {
		students := value.Students
		sort.Slice(students, func(i, j int) bool { return students[i].LatestRanking < students[j].LatestRanking })
		results[index].Students = students
	}

	return c.JSON(http.StatusOK, results)
}

// GetStudentRankingAnalysisHandler 获取某个班级的个人排名分析
func GetStudentRankingAnalysisHandler(c echo.Context) error {
	schoolID := c.QueryParam("schoolID")
	grade := c.QueryParam("grade")
	class, err := strconv.Atoi(c.QueryParam("class"))
	if err != nil {
		return utils.InvalidParams("class is not a number!")
	}
	startTimeInt64, err := strconv.ParseInt(c.QueryParam("startTime"), 10, 64)
	if err != nil {
		return utils.InvalidParams("startTime is not a number!")
	}
	startTime := time.Unix(startTimeInt64, 0)
	endTimeInt64, err := strconv.ParseInt(c.QueryParam("endTime"), 10, 64)
	if err != nil {
		return utils.InvalidParams("endTime is not a number!")
	}
	endTime := time.Unix(endTimeInt64, 0)
	standardFullScore, err := strconv.Atoi(c.QueryParam("standardFullScore"))
	if err != nil {
		return utils.InvalidParams("standardFullScore is not a number!")
	}

	type examRankingType struct {
		Time    int64  `json:"time"`    // 考试时间，unix时间戳（注意因为界面只选择了日期，用0时0分0秒）
		PaperID string `json:"paperID"` // 考试试卷ID
		Ranking int    `json:"ranking"` // 这次考试排名
	}
	type studentRankingType struct {
		LearnID       int               `json:"learnID"`       // 学习号
		Name          string            `json:"name"`          // 学生名字
		LatestRanking int               `json:"latestRanking"` // 最新排名（右侧表格展示）
		Data          []examRankingType `json:"data"`
	}
	type levelRankingType struct {
		Level    int                  `json:"level"`
		Students []studentRankingType `json:"students"`
	}
	results := []levelRankingType{}

	exams, err := getExams(schoolID, grade, class, startTime, endTime, standardFullScore)
	if err != nil || len(exams) <= 0 {
		return c.JSON(http.StatusOK, results)
	}

	studentsMap := make(map[int]studentRankingType)
	for _, exam := range exams {
		scores := exam.Scores
		sort.Sort(scores)
		lastRanking := 1
		var lastScore float32 = -100
		for index, stu := range scores {
			examRanking := examRankingType{
				Time:    exam.Time.Unix(),
				PaperID: exam.PaperID,
			}
			if lastScore != stu.Score {
				examRanking.Ranking = index + 1
				lastRanking = index + 1
				lastScore = stu.Score
			} else {
				// 成绩相同则排名相同
				examRanking.Ranking = lastRanking
			}
			if studentInMap, ok := studentsMap[stu.LearnID]; !ok {
				studentsMap[stu.LearnID] = studentRankingType{
					LearnID: stu.LearnID,
					Name:    stu.Name,
					Data: []examRankingType{
						examRanking,
					},
				}
			} else {
				studentInMap.Data = append(studentInMap.Data, examRanking)
				studentsMap[stu.LearnID] = studentInMap
			}
		}
	}

	latestScores := exams[len(exams)-1].Scores
	sort.Sort(latestScores)
	lastRanking := 1
	var lastScore float32 = -100
	for index, stu := range latestScores {
		studentInMap := studentsMap[stu.LearnID]
		if lastScore != stu.Score {
			studentInMap.LatestRanking = index + 1
			lastRanking = index + 1
			lastScore = stu.Score
		} else {
			// 成绩相同则排名相同
			studentInMap.LatestRanking = lastRanking
		}
		studentsMap[stu.LearnID] = studentInMap
	}

	// 对学生分层整理
	levelMap := make(map[int]levelRankingType)
	for _, stu := range studentsMap {
		level, err := getStudentLevel(stu.LearnID)
		if err != nil || level == 0 || level == -1 {
			return utils.Forbidden("some students do not have levels")
		}
		if levelInMap, ok := levelMap[level]; !ok {
			levelMap[level] = levelRankingType{
				Level: level,
				Students: []studentRankingType{
					stu,
				},
			}
		} else {
			levelInMap.Students = append(levelInMap.Students, stu)
			levelMap[level] = levelInMap
		}
	}

	// 转为list
	for _, value := range levelMap {
		results = append(results, value)
	}
	sort.Slice(results, func(i, j int) bool { return results[i].Level < results[j].Level })

	for index, value := range results {
		students := value.Students
		sort.Slice(students, func(i, j int) bool { return students[i].LatestRanking < students[j].LatestRanking })
		results[index].Students = students
	}

	return c.JSON(http.StatusOK, results)
}

// UploadExamThoughtsHandler 上传教学思考
func UploadExamThoughtsHandler(c echo.Context) error {
	type examThoughtsType struct {
		SchoolID     string `json:"schoolID"`     // 学校识别码
		Grade        string `json:"grade"`        // 年级
		Class        int    `json:"class"`        // 班级, 0 代表全部
		PaperID      string `json:"paperID"`      // 考试试卷ID
		ImageType    int    `json:"imageType"`    // 1 班级均分分析图，2 班级排名段分析图 3 班级分数段分析图  4 个人分数分析图  5 个人排名分析图
		ScoreSegment int    `json:"scoreSegment"` // 分数段(当为班级分数段分析图才有效)
		Level        int    `json:"level"`        // 层级(个人分数分析或者个人排名分析才有效)
		Thoughts     string `json:"thoughts"`     // 教学思考
	}

	var uploadedData examThoughtsType
	if err := c.Bind(&uploadedData); err != nil {
		return utils.InvalidParams("invalid input!" + err.Error())
	}

	key := fmt.Sprintf("%s*%d*%d*%d", uploadedData.PaperID, uploadedData.ImageType, uploadedData.ScoreSegment, uploadedData.Level)

	_, err := userDB.C("classes").Upsert(bson.M{
		"schoolID": bson.ObjectIdHex(uploadedData.SchoolID),
		"grade":    uploadedData.Grade,
		"class":    uploadedData.Class,
		"valid":    true,
	}, bson.M{
		"$set": bson.M{
			"examThoughts." + key: uploadedData.Thoughts,
		},
	})
	if err != nil {
		log.Printf("failed to save thoughts in classes, error %v\n", err)
		return err
	}
	return c.JSON(http.StatusOK, "successfully added thoughts")
}
