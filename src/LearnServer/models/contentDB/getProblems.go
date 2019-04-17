package contentDB

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

// DetailedProblem 数据库中的problem具体信息
type DetailedProblem struct {
	Book       string `json:"book,omitempty" db:"book"`
	Page       int    `json:"page,omitempty" db:"page"`
	LessonName string `json:"lessonName" db:"lessonName"`
	Column     string `json:"column" db:"column"`
	ProblemID  string `json:"problemID" db:"problemID"`
	Idx        int    `json:"idx" db:"idx"`
	SubIdx     int    `json:"subIdx" db:"subIdx"`
	Type       string `json:"-" db:"type"`
	SourceID   string `json:"-" db:"-"`
}

// GetNonExampleProblemsByBookAndPage 根据某一本书某一页获取不是例题的题目
func GetNonExampleProblemsByBookAndPage(bookID string, page int) ([]DetailedProblem, error) {
	db := GetDB()
	tx, err := db.Beginx()
	if err != nil {
		return []DetailedProblem{}, err
	}
	defer tx.Rollback()
	stmt, err := tx.Preparex("SELECT m.problemID, m.`column`, m.num as idx, t.subIdx FROM probmetas as m, probtypes as t WHERE m.bookID = ? and m.page = ? and t.problemID = m.problemID and not EXISTS(SELECT * FROM examples as e WHERE e.problemID = m.problemID and e.subIdx = t.subIdx LIMIT 1) ORDER BY idx, subIdx;")
	if err != nil {
		return []DetailedProblem{}, err
	}
	rows, err := stmt.Queryx(bookID, page)
	if err != nil {
		return []DetailedProblem{}, err
	}

	problems := []DetailedProblem{}
	for rows.Next() {
		// 未知课时名称
		p := DetailedProblem{
			LessonName: "",
		}
		err := rows.StructScan(&p)
		if err != nil {
			log.Println(err)
		} else {
			if p.Column != "练习题" {
				problems = append(problems, p)
			}
		}
	}
	return problems, nil
}

// GetProblemsByChapterSection 得到某一章节的题目
func GetProblemsByChapterSection(chapter int, section int) ([]DetailedProblem, error) {
	db := GetDB()
	tx, err := db.Beginx()
	if err != nil {
		return []DetailedProblem{}, err
	}
	defer tx.Rollback()

	stmt, err := tx.Preparex("SELECT m.problemID, m.subIdx FROM probtypes as m, probmetas as n where m.problemID = n.problemID and n.chapNum = ? and n.sectNum = ?;")
	if err != nil {
		return []DetailedProblem{}, err
	}

	var result []DetailedProblem
	rows, err := stmt.Queryx(chapter, section)
	if err != nil {
		return []DetailedProblem{}, err
	}
	for rows.Next() {
		tmpP := DetailedProblem{}
		if err := rows.StructScan(&tmpP); err != nil {
			log.Println(err)
		} else {
			result = append(result, tmpP)
		}
	}

	return result, nil
}

// GetProblemsByBookPage 得到某一本书某些页的题目
func GetProblemsByBookPage(bookID string, startPage int, endPage int) ([]DetailedProblem, error) {
	db := GetDB()
	tx, err := db.Beginx()
	if err != nil {
		return []DetailedProblem{}, err
	}
	defer tx.Rollback()

	stmt, err := tx.Preparex(`SELECT m.problemID, m.subIdx, n.num as idx
							  FROM probtypes as m, probmetas as n
							  WHERE m.problemID = n.problemID and n.bookID = ? and n.page >= ? and n.page <= ?
							  ORDER BY idx, subIdx;`)
	if err != nil {
		return []DetailedProblem{}, err
	}

	var result []DetailedProblem
	rows, err := stmt.Queryx(bookID, startPage, endPage)
	if err != nil {
		return []DetailedProblem{}, err
	}
	for rows.Next() {
		tmpP := DetailedProblem{}
		if err := rows.StructScan(&tmpP); err != nil {
			log.Println(err)
		} else {
			result = append(result, tmpP)
		}
	}

	return result, nil
}

// GetProblemsByPaper 得到某一Paper的题目
func GetProblemsByPaper(paperID string) ([]DetailedProblem, error) {
	db := GetDB()
	tx, err := db.Beginx()
	if err != nil {
		return []DetailedProblem{}, err
	}
	defer tx.Rollback()

	stmt, err := tx.Preparex(`SELECT m.problemID, n.subIdx, m.problemIndex as idx
							  FROM examproblem as m, probtypes as n
							  where m.problemID = n.problemID and m.examPaperID = ?
							  ORDER BY idx, subIdx;`)
	if err != nil {
		return []DetailedProblem{}, err
	}

	var result []DetailedProblem
	rows, err := stmt.Queryx(paperID)
	if err != nil {
		return []DetailedProblem{}, err
	}
	for rows.Next() {
		tmpP := DetailedProblem{}
		if err := rows.StructScan(&tmpP); err != nil {
			log.Println(err)
		} else {
			result = append(result, tmpP)
		}
	}

	return result, nil
}

// GetAllProblems 获取数据库中所有题目
func GetAllProblems() ([]DetailedProblem, error) {
	db := GetDB()
	problems := []DetailedProblem{}
	err := db.Select(&problems, "SELECT problemID, subIdx FROM probtypes;")
	return problems, err
}

// GetHow 得到这道题的出题方式
func GetHow(problemID string) (string, error) {
	db := GetDB()
	how := ""
	err := db.Get(&how, "SELECT how FROM hows WHERE problemID = ?;", problemID)
	return how, err
}

// ScanDetailedProblem 得到题目具体信息, data: 保存得到的数据的指针, bookIDsFilter 题目来源优先匹配这些bookID
func ScanDetailedProblem(problemID string, subIdx int, data interface{}, bookIDsFilter []string) error {
	db := GetDB()
	query, args, err := sqlx.In(
		"SELECT m.problemID, t.subIdx, m.source as book, m.`column`, m.num as idx, m.page, t.typename as type "+
			"FROM probmetas as m, probtypes as t "+
			"WHERE m.problemID = t.problemID and m.problemID = ? and t.subIdx = ? and m.bookID IN (?);", problemID, subIdx, bookIDsFilter)
	if err != nil {
		return err
	}
	err = db.QueryRowx(db.Rebind(query), args...).StructScan(data)
	if err != nil {
		err = db.QueryRowx(
			"SELECT m.problemID, t.subIdx, m.source as book, m.`column`, m.num as idx, m.page, t.typename as type "+
				"FROM probmetas as m, probtypes as t "+
				"WHERE m.problemID = t.problemID and m.problemID = ? and t.subIdx = ?;", problemID, subIdx).StructScan(data)
		if err != nil {
			err = db.QueryRowx(
				"SELECT m.problemID, t.subIdx, m.examPaperName as book, m.problemIndex as idx, t.typename as type "+
					"FROM examproblem as m, probtypes as t "+
					"WHERE m.problemID = t.problemID and m.problemID = ? and t.subIdx = ?;", problemID, subIdx).StructScan(data)
		}
	}
	return err
}

// GetSourceType 获取某个题目的来源类型
func GetSourceType(problemID string) (int, error) {
	db := GetDB()
	bookType := -1
	err := db.Get(&bookType, "select b.type from books as b, probmetas as p where p.problemID = ? and p.bookID = b.bookID;", problemID)
	if err != nil {
		return -1, err
	}
	if bookType == 1 {
		// 该题目在课本中
		return 3, nil
	}
	if bookType == 2 {
		// 该题目在普通辅导书中
		return 2, nil
	}
	examPaperID := ""
	err = db.Get(&examPaperID, "select examPaperID from examproblem where problemID = ?;", problemID)
	if err != nil {
		return -1, err
	}
	if examPaperID == "" {
		return -1, fmt.Errorf("can not find paperID of problemID %s", problemID)
	}
	// 题目在试卷中
	return 1, nil
}
