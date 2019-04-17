package students

import (
	"log"
	"net/http"

	"LearnServer/models/contentDB"
	"LearnServer/models/userDB"
	"LearnServer/services/students/validation"
	"LearnServer/utils"
	"github.com/labstack/echo"
)

// getNotMarkedPapersHandler 获取有哪些未标记试卷
func getNotMarkedPapersHandler(c echo.Context) error {
	var id string
	err := validation.ValidateUser(c, &id)
	if err != nil {
		return err
	}

	paperIDs, err := userDB.GetNotMarkedPaperIDs(id)
	if err != nil {
		log.Printf("getting not marked papers of student id %s failed, err %v", id, err)
		return err
	}
	if len(paperIDs) <= 0 {
		return utils.NotFound("no papers.")
	}

	return c.JSON(http.StatusOK, contentDB.GetPapersByPaperID(paperIDs))
}

// 获取有哪些已经标记的试卷
func getMarkedPapersHandler(c echo.Context) error {
	var id string
	err := validation.ValidateUser(c, &id)
	if err != nil {
		return err
	}

	paperIDs, err := userDB.GetMarkedPaperIDs(id)
	if err != nil {
		log.Printf("getting marked papers of student id %s failed, err %v", id, err)
		return err
	}
	if len(paperIDs) <= 0 {
		return utils.NotFound("no papers.")
	}

	return c.JSON(http.StatusOK, contentDB.GetPapersByPaperID(paperIDs))
}
