package contentDB

import (
	"log"

	//"LearnServer/conf"
	"LearnServer/conf"
)

// BookDetailType 书本信息类型
type BookDetailType struct {
	BookID        string `json:"bookID" db:"bookID"`
	Type          int    `json:"type" db:"type"`
	Time          string `json:"time" db:"uploadTime"`
	Name          string `json:"name" db:"name"`
	Term          string `json:"term" db:"term"`
	Version       string `json:"version" db:"version"`
	Year          int    `json:"year" db:"year"`
	Isbn          string `json:"isbn" db:"isbn"`
	EdiYear       string `json:"ediYear" db:"ediYear"`
	EdiMonth      string `json:"ediMonth" db:"ediMonth"`
	EdiVersion    string `json:"ediVersion" db:"ediVersion"`
	ImpYear       string `json:"impYear" db:"impYear"`
	ImpMonth      string `json:"impMonth" db:"impMonth"`
	ImpNum        string `json:"impNum" db:"impNum"`
	CoverFileName string `json:"-" db:"coverFileName"`
	CIPFileName   string `json:"-" db:"cipFileName"`
	PriceFileName string `json:"-" db:"priceFileName"`
	CoverURL      string `json:"coverURL"`
	CIPURL        string `json:"cipURL"`
	PriceURL      string `json:"priceURL"`
}

// GetBooksByBookID 根据书本识别码获取书本信息
func GetBooksByBookID(bookIDs []string) []BookDetailType {
	URL_PREFIX := conf.AppConfig.BookImagesURL

	db := GetDB()
	result := []BookDetailType{}
	for _, id := range bookIDs {
		book := BookDetailType{}
		err := db.Get(&book, `SELECT bookID, type, name, term, version, year, isbn, ediYear, ediMonth, ediVersion, impYear, impMonth, impNum, cipFileName, priceFileName, coverFileName, uploadTime FROM books
							  WHERE bookID = ?;`, id)
		if err != nil {
			log.Printf("getting book %s failed, err: %v \n", id, err)
			continue
		}
		result = append(result, book)
	}

	for i, b := range result {
		result[i].CoverURL = URL_PREFIX + b.CoverFileName
		result[i].CIPURL = URL_PREFIX + b.CIPFileName
		result[i].PriceURL = URL_PREFIX + b.PriceFileName
		// 截取日期
		result[i].Time = b.Time[:10]
	}
	return result
}
