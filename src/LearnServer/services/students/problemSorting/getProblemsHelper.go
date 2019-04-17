package problemSorting

import (
	// "LearnServer/models/contentDB"
	// "LearnServer/models/userDB"
	"LearnServer/models/contentDB"
	"LearnServer/models/userDB"
)

func getAllProblemsDone(id string) ([]problemWithTypeTimeAndCorrectDB, error) {
	// 获取所有做过的题目
	allProblemsDoneRaw := struct {
		Problems []problemWithTypeTimeAndCorrectDB `bson:"problems"`
	}{}
	err := userDB.GetAllProblemsDone(id, &allProblemsDoneRaw)
	return allProblemsDoneRaw.Problems, err
}

func getProblemsDoneOfChapSect(id string, chapter int, section int) ([]problemWithTypeTimeAndCorrectDB, error) {
	// 获取符合章节要求的做过的题目
	allProblemsDone, err := getAllProblemsDone(id)
	if err != nil {
		return nil, err
	}

	problems := []problemWithTypeTimeAndCorrectDB{}

	db := contentDB.GetDB()
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	stmtGetWrongProblemsOfChaptSect, err := tx.Preparex("SELECT chapNum as c, sectNum as s FROM probmetas where problemID = ?;")
	if err != nil {
		return nil, err
	}

	for _, p := range allProblemsDone {
		var c, s int
		stmtGetWrongProblemsOfChaptSect.QueryRowx(p.ProblemID).Scan(&c, &s)
		if c == chapter && s == section {
			problems = append(problems, p)
		}
	}

	return problems, nil
}
