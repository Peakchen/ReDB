package staffs

import (
	"net/http"

	"LearnServer/conf"
	"LearnServer/models/contentDB"
	"github.com/labstack/echo"
)

func searchBookHandler(c echo.Context) error {
	URL_PREFIX := conf.AppConfig.BookImagesURL

	isbn := c.QueryParam("isbn")
	ediYear := c.QueryParam("ediYear")
	ediMonth := c.QueryParam("ediMonth")
	ediVersion := c.QueryParam("ediVersion")
	impYear := c.QueryParam("impYear")
	impMonth := c.QueryParam("impMonth")
	impNum := c.QueryParam("impNum")

	type bookType struct {
		BookID        string `json:"bookID" db:"bookID"`
		CoverFileName string `json:"-" db:"coverFileName"`
		CIPFileName   string `json:"-" db:"cipFileName"`
		PriceFileName string `json:"-" db:"priceFileName"`
		CoverURL      string `json:"coverURL"`
		CIPURL        string `json:"cipURL"`
		PriceURL      string `json:"priceURL"`
	}

	result := []bookType{}
	db := contentDB.GetDB()
	err := db.Select(&result, `SELECT bookID, cipFileName, coverFileName, priceFileName FROM books
							  WHERE isbn LIKE ? and ediYear LIKE ? and ediMonth LIKE ? and ediVersion LIKE ? and impYear Like ? and impMonth LIKE ? and impNum like ?;`, "%"+isbn+"%", "%"+ediYear+"%", "%"+ediMonth+"%", "%"+ediVersion+"%", "%"+impYear+"%", "%"+impMonth+"%", "%"+impNum+"%")
	if err != nil {
		return err
	}

	for i, b := range result {
		result[i].CoverURL = URL_PREFIX + b.CoverFileName
		result[i].CIPURL = URL_PREFIX + b.CIPFileName
		result[i].PriceURL = URL_PREFIX + b.PriceFileName
	}
	return c.JSON(http.StatusOK, result)
}
