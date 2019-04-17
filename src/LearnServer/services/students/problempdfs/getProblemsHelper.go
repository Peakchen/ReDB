package problempdfs

import (
	"fmt"
	"log"
	"strconv"
	"time"

	// "LearnServer/models/contentDB"
	// "LearnServer/models/userDB"

	"LearnServer/models/contentDB"
	"LearnServer/models/userDB"
)

func getNewestWrongProblems(id string, timeBefore time.Time) ([]detailedProblem, error) {
	// 获取 timeBefore 之前的最新做错的题目
	type problemWithTimeAndCorrectDB struct {
		ProblemID string    `bson:"problemID" json:"problemID"`
		SubIdx    int       `bson:"subIdx" json:"subIdx"`
		Correct   bool      `bson:"correct" json:"isCorrect"`
		Time      time.Time `bson:"assignDate" json:"assignDate"`
	}

	// 所有做过的题目
	allProblemsDoneRaw := struct {
		Problems []problemWithTimeAndCorrectDB `bson:"problems"`
	}{}
	if err := userDB.GetAllProblemsDone(id, &allProblemsDoneRaw); err != nil {
		return []detailedProblem{}, err
	}

	wrongProblemsAll := []detailedProblem{}

	// 得到最新的错题（最新做错的题目），并且加到wrongProblemsAll中
	for _, p := range allProblemsDoneRaw.Problems {
		if p.Correct || timeBefore.Before(p.Time) {
			continue
		}

		newest := true

		for _, pBefore := range allProblemsDoneRaw.Problems {
			if (!timeBefore.Before(pBefore.Time)) && p.ProblemID == pBefore.ProblemID && p.SubIdx == pBefore.SubIdx {
				if p.Time.Before(pBefore.Time) {
					// pBefore比p更加新
					newest = false
					break
				}
			}
		}

		if newest {
			wrongProblemsAll = append(wrongProblemsAll, detailedProblem{
				ProblemID: p.ProblemID,
				SubIdx:    p.SubIdx,
				Full:      true,
			})
		}
	}

	return wrongProblemsAll, nil
}

func getNewestWrongProblemsOfChapSect(id string, chapter int, section int, timeBefore time.Time) ([]detailedProblem, error) {
	// 获取某一章节的最新做错的题目

	// 获取最新错题
	wrongProblemsAll, err := getNewestWrongProblems(id, timeBefore)
	if err != nil {
		return []detailedProblem{}, err
	}

	contentdb := contentDB.GetDB()
	tx, err := contentdb.Beginx()
	if err != nil {
		return []detailedProblem{}, err
	}
	defer tx.Rollback()

	// problems符合章节要求的最新错题
	problems := []detailedProblem{}
	stmtGetWrongProblemsOfChaptSect, err := tx.Preparex("SELECT chapNum as c, sectNum as s FROM probmetas where problemID = ?;")
	if err != nil {
		return []detailedProblem{}, err
	}

	for _, p := range wrongProblemsAll {
		var c, s int
		stmtGetWrongProblemsOfChaptSect.QueryRowx(p.ProblemID).Scan(&c, &s)
		if c == chapter && s == section {
			problems = append(problems, p)
		}
	}

	return problems, nil
}

func getOneProblemNotDone(id string, problemsToCheck []problem) (problem, error) {
	// 从problemsToCheck中找一道没做过的题目
	type problemWithTimeAndCorrectDB struct {
		ProblemID string    `bson:"problemID" json:"problemID"`
		SubIdx    int       `bson:"subIdx" json:"subIdx"`
		Correct   bool      `bson:"correct" json:"isCorrect"`
		Time      time.Time `bson:"assignDate" json:"assignDate"`
	}

	// 所有做过的题目
	allProblemsDoneRaw := struct {
		Problems []problemWithTimeAndCorrectDB `bson:"problems"`
	}{}
	if err := userDB.GetAllProblemsDone(id, &allProblemsDoneRaw); err != nil {
		return problem{}, err
	}

	for _, p := range problemsToCheck {
		done := false
		for _, pDone := range allProblemsDoneRaw.Problems {
			if p.ProblemID == pDone.ProblemID && p.SubIdx == pDone.SubIdx {
				done = true
				break
			}
		}

		if !done {
			return p, nil
		}
	}

	return problem{}, fmt.Errorf("can find a problem that hasn't been done")
}

func getProblemsOfHowAndType(howCode int, typeName string) ([]problem, error) {
	// 获取特定类型和题型的题目
	db := contentDB.GetDB()
	tx, err := db.Beginx()
	if err != nil {
		return []problem{}, err
	}
	defer tx.Rollback()

	stmtStr := ""
	switch howCode {
	case 0:
		stmtStr = "SELECT t.problemID as problemID, t.subIdx FROM hows as h, probtypes as t WHERE t.problemID = h.problemID and t.typeName = ? and h.how = '选择题';"
	case 1:
		stmtStr = "SELECT t.problemID as problemID, t.subIdx FROM hows as h, probtypes as t WHERE t.problemID = h.problemID and t.typeName = ? and h.how = '填空题';"
	default:
		stmtStr = "SELECT t.problemID as problemID, t.subIdx FROM hows as h, probtypes as t WHERE t.problemID = h.problemID and t.typeName = ? and h.how != '选择题' and h.how != '填空题';"
	}

	stmt, err := tx.Preparex(stmtStr)
	if err != nil {
		return []problem{}, err
	}

	result := []problem{}
	rows, err := stmt.Queryx(typeName)
	if err != nil {
		return []problem{}, err
	}
	for rows.Next() {
		tmpP := problem{}
		if err := rows.StructScan(&tmpP); err != nil {
			log.Println(err)
		} else {
			result = append(result, tmpP)
		}
	}

	return result, nil
}

// GetNewestWrongProblemsOfBookPage 获取某一本书某些页的最新做错的题目
func GetNewestWrongProblemsOfBookPage(id string, bookID string, startPage int, endPage int, timeBefore time.Time) ([]detailedProblem, error) {

	// 获取最新错题
	wrongProblemsAll, err := getNewestWrongProblems(id, timeBefore)
	if err != nil {
		return []detailedProblem{}, err
	}

	contentdb := contentDB.GetDB()
	tx, err := contentdb.Beginx()
	if err != nil {
		return []detailedProblem{}, err
	}
	defer tx.Rollback()

	// problems符合章节要求的最新错题
	problems := []detailedProblem{}
	stmtGetWrongProblemsOfChaptSect, err := tx.Preparex("SELECT bookID as b, page as p FROM probmetas where problemID = ?;")
	if err != nil {
		return []detailedProblem{}, err
	}

	for _, pro := range wrongProblemsAll {
		var b string
		var p int
		stmtGetWrongProblemsOfChaptSect.QueryRowx(pro.ProblemID).Scan(&b, &p)
		if b == bookID && p >= startPage && p <= endPage {
			problems = append(problems, pro)
		}
	}

	return problems, nil
}

func getOnceWrongProblems(id string, timeBefore time.Time) ([]detailedProblem, error) {
	// 获取 timeBefore 之前曾经错过的题目
	type problemWithTimeAndCorrectDB struct {
		ProblemID string    `bson:"problemID" json:"problemID"`
		SubIdx    int       `bson:"subIdx" json:"subIdx"`
		Correct   bool      `bson:"correct" json:"isCorrect"`
		Time      time.Time `bson:"assignDate" json:"assignDate"`
	}

	// 所有做过的题目
	allProblemsDoneRaw := struct {
		Problems []problemWithTimeAndCorrectDB `bson:"problems"`
	}{}
	if err := userDB.GetAllProblemsDone(id, &allProblemsDoneRaw); err != nil {
		return []detailedProblem{}, err
	}

	onceWrongProblems := []detailedProblem{}

	// 得到曾经错过的，并且加到onceWrongProblems中
	for _, p := range allProblemsDoneRaw.Problems {
		if p.Correct || timeBefore.Before(p.Time) {
			continue
		}

		found := false
		for _, pPut := range onceWrongProblems {
			if p.ProblemID == pPut.ProblemID && p.SubIdx == pPut.SubIdx {
				found = true
				break
			}
		}
		if !found {
			// 避免重复添加
			onceWrongProblems = append(onceWrongProblems, detailedProblem{
				ProblemID: p.ProblemID,
				SubIdx:    p.SubIdx,
				Full:      true,
			})
		}
	}
	return onceWrongProblems, nil
}

// GetOnceWrongProblemsOfBookPage 获取某一本书某些页的曾经做错的题目
func GetOnceWrongProblemsOfBookPage(id string, bookID string, startPage int, endPage int, timeBefore time.Time) ([]detailedProblem, error) {

	// 获取全部曾经做错的题
	wrongProblemsAll, err := getOnceWrongProblems(id, timeBefore)
	if err != nil {
		return []detailedProblem{}, err
	}

	contentdb := contentDB.GetDB()
	tx, err := contentdb.Beginx()
	if err != nil {
		return []detailedProblem{}, err
	}
	defer tx.Rollback()

	// problems符合章节要求的最新错题
	problems := []detailedProblem{}
	stmtGetWrongProblemsOfChaptSect, err := tx.Preparex("SELECT bookID as b, page as p FROM probmetas where problemID = ?;")
	if err != nil {
		return []detailedProblem{}, err
	}

	for _, pro := range wrongProblemsAll {
		var b string
		var p int
		stmtGetWrongProblemsOfChaptSect.QueryRowx(pro.ProblemID).Scan(&b, &p)
		if b == bookID && p >= startPage && p <= endPage {
			problems = append(problems, pro)
		}
	}

	return problems, nil
}

func getKnownWrongProblems(id string) ([]detailedProblem, error) {
	// 获取曾经错过现在做对的题目
	type problemWithTimeAndCorrectDB struct {
		ProblemID string    `bson:"problemID" json:"problemID"`
		SubIdx    int       `bson:"subIdx" json:"subIdx"`
		Correct   bool      `bson:"correct" json:"isCorrect"`
		Time      time.Time `bson:"assignDate" json:"assignDate"`
	}

	// 所有做过的题目
	allProblemsDoneRaw := struct {
		Problems []problemWithTimeAndCorrectDB `bson:"problems"`
	}{}
	if err := userDB.GetAllProblemsDone(id, &allProblemsDoneRaw); err != nil {
		return []detailedProblem{}, err
	}

	// problemStatusMap 记录problemID+subIdx与题目状态的对应关系
	type statusType struct {
		ProblemID string
		SubIdx    int
		Status    bool // false错了，true代表有错过但最新是对的
		Time      time.Time
	}
	problemStatusMap := make(map[string]statusType)

	for _, p := range allProblemsDoneRaw.Problems {
		// 先找到所有做错过的题目
		if !p.Correct {
			pkey := p.ProblemID + strconv.Itoa(p.SubIdx)
			status, ok := problemStatusMap[pkey]
			if !ok {
				problemStatusMap[pkey] = statusType{
					ProblemID: p.ProblemID,
					SubIdx:    p.SubIdx,
					Time:      p.Time,
					Status:    false,
				}
			} else {
				if status.Time.Before(p.Time) {
					// 更新时间，确保status中记录的是错题做错的最新时间
					problemStatusMap[pkey] = statusType{
						ProblemID: status.ProblemID,
						SubIdx:    status.SubIdx,
						Time:      p.Time,
						Status:    false,
					}
				}
			}
		}
	}

	for _, p := range allProblemsDoneRaw.Problems {
		if p.Correct {
			pkey := p.ProblemID + strconv.Itoa(p.SubIdx)
			status, ok := problemStatusMap[pkey]
			if ok && status.Time.Before(p.Time) {
				// 之后做对了
				problemStatusMap[pkey] = statusType{
					ProblemID: status.ProblemID,
					SubIdx:    status.SubIdx,
					Time:      status.Time,
					Status:    true,
				}
			}
		}
	}

	knownWrongProblems := []detailedProblem{}
	for _, status := range problemStatusMap {
		if status.Status {
			knownWrongProblems = append(knownWrongProblems, detailedProblem{
				ProblemID: status.ProblemID,
				SubIdx:    status.SubIdx,
				Full:      true,
			})
		}
	}

	return knownWrongProblems, nil
}

func getKnownWrongProblemsOfBookPage(id string, book string, startPage int, endPage int) ([]detailedProblem, error) {
	// 获取某一本书某些页的曾经做错现在会做的题目

	// 获取全部曾经做错现在会做的题
	wrongProblemsAll, err := getKnownWrongProblems(id)
	if err != nil {
		return []detailedProblem{}, err
	}

	contentdb := contentDB.GetDB()
	tx, err := contentdb.Beginx()
	if err != nil {
		return []detailedProblem{}, err
	}
	defer tx.Rollback()

	// problems符合章节要求的最新错题
	problems := []detailedProblem{}
	stmtGetWrongProblemsOfChaptSect, err := tx.Preparex("SELECT source as b, page as p FROM probmetas where problemID = ?;")
	if err != nil {
		return []detailedProblem{}, err
	}

	for _, pro := range wrongProblemsAll {
		var b string
		var p int
		stmtGetWrongProblemsOfChaptSect.QueryRowx(pro.ProblemID).Scan(&b, &p)
		if b == book && p >= startPage && p <= endPage {
			problems = append(problems, pro)
		}
	}

	return problems, nil
}

// GetNewestWrongPaperProblems 获取某个试卷的最新做错的题目
func GetNewestWrongPaperProblems(id string, paperID string, timeBefore time.Time) ([]detailedProblem, error) {

	// 获取最新错题
	wrongProblemsAll, err := getNewestWrongProblems(id, timeBefore)
	if err != nil {
		return []detailedProblem{}, err
	}

	contentdb := contentDB.GetDB()
	tx, err := contentdb.Beginx()
	if err != nil {
		return []detailedProblem{}, err
	}
	defer tx.Rollback()

	// problems符合要求的最新错题
	problems := []detailedProblem{}
	stmtGetWrongProblemsOfPaper, err := tx.Preparex("SELECT examPaperID FROM examproblem where problemID = ?;")
	if err != nil {
		return []detailedProblem{}, err
	}

	for _, pro := range wrongProblemsAll {
		var examPaperID string
		stmtGetWrongProblemsOfPaper.QueryRowx(pro.ProblemID).Scan(&examPaperID)
		if examPaperID == paperID {
			problems = append(problems, pro)
		}
	}

	return problems, nil
}

// GetOncePaperWrongProblems 获取某个试卷曾经做错的题目
func GetOncePaperWrongProblems(id string, paperID string, timeBefore time.Time) ([]detailedProblem, error) {

	// 获取全部曾经做错的题
	wrongProblemsAll, err := getOnceWrongProblems(id, timeBefore)
	if err != nil {
		return []detailedProblem{}, err
	}

	contentdb := contentDB.GetDB()
	tx, err := contentdb.Beginx()
	if err != nil {
		return []detailedProblem{}, err
	}
	defer tx.Rollback()

	// problems符合要求的最新错题
	problems := []detailedProblem{}
	stmtGetWrongProblemsOfPaper, err := tx.Preparex("SELECT examPaperID FROM examproblem where problemID = ?;")
	if err != nil {
		return []detailedProblem{}, err
	}

	for _, pro := range wrongProblemsAll {
		var examPaperID string
		stmtGetWrongProblemsOfPaper.QueryRowx(pro.ProblemID).Scan(&examPaperID)
		if examPaperID == paperID {
			problems = append(problems, pro)
		}
	}

	return problems, nil
}
