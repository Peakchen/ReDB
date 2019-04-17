package problempdfs

import (
	"log"
	"sort"
	"strings"
	"time"

	// "LearnServer/models/contentDB"
	// "LearnServer/models/userDB"
	
	"LearnServer/models/contentDB"
	"LearnServer/models/userDB"
)

type problemDBType struct {
	ProblemID string    `bson:"problemID"`
	SubIdx    int       `bson:"subIdx"`
	Correct   bool      `bson:"correct"`
	Time      time.Time `bson:"assignDate"`
	Type      int       `bson:"type"`
}

// 实现sort，用来合并错题状态
type problemsToMergeSliceType []problemDBType

func (c problemsToMergeSliceType) Len() int {
	return len(c)
}

func (c problemsToMergeSliceType) Swap(i int, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c problemsToMergeSliceType) Less(i int, j int) bool {
	if strings.Compare(c[i].ProblemID, c[j].ProblemID) != 0 {
		return strings.Compare(c[i].ProblemID, c[j].ProblemID) < 0
	}

	if c[i].SubIdx != c[j].SubIdx {
		return c[i].SubIdx < c[j].SubIdx
	}

	return c[i].Time.Before(c[j].Time)
}

type problemsForPickSliceType []detailedProblem

func (c problemsForPickSliceType) Len() int {
	return len(c)
}

func (c problemsForPickSliceType) Swap(i int, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c problemsForPickSliceType) Less(i int, j int) bool {
	if !c[i].newestState && c[j].newestState {
		return true
	}
	if c[i].newestState && !c[j].newestState {
		return false
	}
	if c[i].markTimes != c[j].markTimes {
		return c[i].markTimes < c[j].markTimes
	}
	if c[i].wrongTimes != c[j].wrongTimes {
		return c[i].wrongTimes < c[j].wrongTimes
	}
	if (!c[i].newestState) && (c[i].markTimes == 1) {
		// c[i], c[j] 最新状态、标记次数此时必定一样
		// 间隔按从小到大排序
		if c[i].newestMarkDuration != c[j].newestMarkDuration {
			return c[i].newestMarkDuration < c[j].newestMarkDuration
		}
	} else {
		// 间隔按从大到小排序
		if c[i].newestMarkDuration != c[j].newestMarkDuration {
			return c[i].newestMarkDuration > c[j].newestMarkDuration
		}
	}
	if c[i].sourceType != c[j].sourceType {
		return c[i].sourceType < c[j].sourceType
	}
	if c[i].howCode != c[j].howCode {
		return c[i].howCode < c[j].howCode
	}

	if strings.Compare(c[i].ProblemID, c[j].ProblemID) != 0 {
		return strings.Compare(c[i].ProblemID, c[j].ProblemID) < 0
	}

	return c[i].SubIdx < c[j].SubIdx
}

func existProbAndAnswerDoc(problemID string) bool {
	// 判断是否题目答案文档都存在
	db := contentDB.GetDB()
	problemPath := ""
	err := db.Get(&problemPath, "select path from problemzip where problemID = ?;", problemID)
	if err != nil || problemPath == "" {
		return false
	}
	answerPath := ""
	err = db.Get(&answerPath, "select path from answerzip where problemID = ?;", problemID)
	if err != nil || answerPath == "" {
		return false
	}
	return true
}

func getProblemReason(p detailedProblem) string {
	// 获取选题依据
	if (!p.newestState) && p.markTimes == 1 {

		switch {
		case p.newestMarkDuration.Hours() <= 168:
			// 24 * 7 = 168 一周168小时
			return "最近错题及时复习"
		case p.newestMarkDuration.Hours() <= 840:
			{
				// 24 * 7 * 5 = 840 5周
				switch p.sourceType {
				case 1:
					return "试卷错题优先对待"
				case 2:
					return "教辅错题考试常见"
				case 3:
					return "课本错题回归基础"
				}
			}
		default:
			return "该错题已超过1个月没复习"

		}
	}
	if (!p.newestState) && p.wrongTimes == 1 {
		return "该题开始做对后来做错"
	}
	if (!p.newestState) && p.wrongTimes == 2 {
		return "该题曾经做错2次"
	}
	if (!p.newestState) && p.wrongTimes >= 3 {
		return "该题曾经做错3次以上"
	}

	if p.newestState {
		switch p.how {
		case "计算题":
			return "多练计算题提高正确率"
		case "应用题":
			return "多练应用题提高审题能力"
		default:
			return "错题连对3次就是高手"
		}
	}

	return ""
}

func formatRawProblems(rawProblems problemsToMergeSliceType) problemsForPickSliceType {
	// 将数据库中经过排序待合并的符合要求的题目数据转换成 detailedProblem，并补充完整题目状态等数据
	problems := problemsForPickSliceType{}
	lastProblemID := ""
	lastSubIdx := -2
	// 上一道题标记次数
	lastMarkTimes := 0
	// 错误次数
	lastWrongTimes := 0
	// 上一道题最新状态
	lastNewestState := false
	// 上一道题最近标记与现在的时间间隔
	var lastNewestMarkDuration time.Duration

	getHowCode := func(how string) int {
		if how == "选择题" {
			return 1
		}
		if how == "填空题" {
			return 2
		}
		if how == "计算题" {
			return 3
		}
		return 4
	}

	addToProblems := func(p problemDBType) {
		// 21 * 24 = 504  3周的小时数
		if (lastNewestState && (lastMarkTimes == 1 || lastNewestMarkDuration.Hours() <= 504)) || !existProbAndAnswerDoc(lastProblemID) {
			// 不选出这些题目
			return
		}

		sourceType, err := contentDB.GetSourceType(lastProblemID)
		if err != nil {
			log.Printf("getting book type failed, problemID: %s, error: %v", lastProblemID, err)
			return
		}

		how, err := contentDB.GetHow(lastProblemID)
		if err != nil {
			log.Printf("getting how failed, problemID: %s, error: %v", lastProblemID, err)
			return
		}

		newProblem := detailedProblem{
			ProblemID:          lastProblemID,
			SubIdx:             lastSubIdx,
			newestState:        lastNewestState,
			markTimes:          lastMarkTimes,
			wrongTimes:         lastWrongTimes,
			sourceType:         sourceType,
			newestMarkDuration: lastNewestMarkDuration,
			how:                how,
			Full:               true,
			howCode:            getHowCode(how),
		}
		newProblem.Reason = getProblemReason(newProblem)

		problems = append(problems, newProblem)
	}

	for i, p := range rawProblems {
		if (lastProblemID != p.ProblemID || lastSubIdx != p.SubIdx) && i != 0 {
			// 到了下一题
			// 因为判断到下一道题才把上一道题添加进problems中，所以跳过第一道
			addToProblems(p)

			lastMarkTimes = 0
			lastWrongTimes = 0
		}
		lastMarkTimes++
		if !p.Correct {
			lastWrongTimes++
		}
		lastNewestState = p.Correct
		lastNewestMarkDuration = time.Now().Sub(p.Time)
		lastProblemID = p.ProblemID
		lastSubIdx = p.SubIdx

		if i == len(rawProblems)-1 {
			// 添加最后一道
			addToProblems(p)
		}
	}

	return problems
}

func getSliceOfProblems(max int, problems []detailedProblem) []detailedProblem {
	// 获取前 max 道题目
	if max < 0 || len(problems) <= max {
		// 若len(problems) == 0，必进入这个if，故之后可认为problems[0]必定存在
		return problems
	}

	count := 0
	index := 1
	for i, p := range problems {
		found := false
		for j := 0; j < i; j++ {
			if problems[j].ProblemID == p.ProblemID {
				found = true
				break
			}
		}
		if !found {
			count++
		}
		index = i + 1
		if count >= max {
			break
		}
	}
	return problems[:index]
}

// PickWrongProblems 挑选EPU2用来生成纠错本的错题
func PickWrongProblems(id string, max int) ([]detailedProblem, error) {

	// 所有做过的题目
	allProblemsDoneRaw := struct {
		Problems problemsToMergeSliceType `bson:"problems"`
	}{}
	if err := userDB.GetAllProblemsDone(id, &allProblemsDoneRaw); err != nil {
		return []detailedProblem{}, err
	}

	// sort 之后同一个题目位置必定连续，而且时间先的记录在前面
	sort.Sort(allProblemsDoneRaw.Problems)

	problems := formatRawProblems(allProblemsDoneRaw.Problems)
	sort.Sort(problems)
	return getSliceOfProblems(max, problems), nil
}
