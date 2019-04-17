package students

import (
	"log"
	"net/http"
	"strconv"

	"LearnServer/models/contentDB"
	"LearnServer/utils"
	"github.com/labstack/echo"
)

func getInfoHandler(c echo.Context) error {
	block, err := strconv.Atoi(c.QueryParam("block"))
	if err != nil {
		block = -1
	}

	chapter, err := strconv.Atoi(c.QueryParam("chapter"))
	if err != nil {
		chapter = -1
	}

	section, err := strconv.Atoi(c.QueryParam("section"))
	if err != nil {
		section = -1
	}

	point, err := strconv.Atoi(c.QueryParam("point"))
	if err != nil {
		point = -1
	}

	type infoType struct {
		Block       string `db:"block" json:"block"`
		Chapter     int    `db:"chapter" json:"chapter"`
		ChapterName string `db:"chapterName" json:"chapterName"`
		Section     int    `db:"section" json:"section"`
		SectionName string `db:"sectionName" json:"sectionName"`
		Point       string `db:"point" json:"point"`
	}

	db := contentDB.GetDB()
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 此处拼接SQL表达式不会引起SQL注入
	queryColumns := ""
	if block == 1 {
		queryColumns += "c.block,"
	}
	if chapter == 1 {
		queryColumns += "c.num as chapter, c.name as chapterName,"
	}
	if section == 1 {
		queryColumns += "b.sectNum as section, s.name as sectionName,"
	}
	if point == 1 {
		queryColumns += "b.name as point,"
	}

	if queryColumns == "" {
		return utils.InvalidParams("Please specify block or chapter or section or point.")
	}

	queryColumns = queryColumns[:len(queryColumns)-1]

	stmt, err := tx.Preparex("SELECT DISTINCT " + queryColumns + " FROM chapters as c, blocks as b, sections as s WHERE c.num = b.chapNum and s.chapNum = c.num and s.num = b.sectNum;")
	if err != nil {
		return err
	}
	rows, err := stmt.Queryx()
	if err != nil {
		return err
	}

	infos := []infoType{}
	for rows.Next() {
		info := infoType{}
		err := rows.StructScan(&info)
		if err != nil {
			log.Println(err)
		} else {
			infos = append(infos, info)
		}
	}
	return c.JSON(http.StatusOK, infos)
}
