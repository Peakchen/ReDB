package staffs

import (
	"net/http"

	"LearnServer/conf"
	"LearnServer/models/contentDB"
	"github.com/labstack/echo"
)

func searchPaperHandler(c echo.Context) error {

	URL_PREFIX := conf.AppConfig.PaperImagesURL

	choice := c.QueryParam("choice")
	blank := c.QueryParam("blank")
	calculation := c.QueryParam("calculation")

	type paperType struct {
		PaperID       string `json:"paperID" db:"paperID"`
		Name          string `json:"name" db:"name"`
		FullScore     int    `json:"fullScore" db:"fullScore"`
		ImageFileName string `json:"-" db:"imageFileName"`
		ImageURL      string `json:"image"`
	}

	result := []paperType{}
	db := contentDB.GetDB()
	err := db.Select(&result, `SELECT paperID, name, fullScore, imageFileName FROM papers
							  WHERE choice LIKE ? and blank LIKE ? and calculation LIKE ?;`, "%"+choice+"%", "%"+blank+"%", "%"+calculation+"%")
	if err != nil {
		return err
	}

	for i, p := range result {
		result[i].ImageURL = URL_PREFIX + p.ImageFileName
	}
	return c.JSON(http.StatusOK, result)
}
