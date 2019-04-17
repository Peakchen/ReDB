package contentDB

import (
	"log"

	//"LearnServer/conf"
	"LearnServer/conf"
)

// PaperDetailType 试卷信息类型
type PaperDetailType struct {
	PaperID       string `json:"paperID" db:"paperID"`
	Time          string `json:"time" db:"uploadTime"`
	Name          string `json:"name" db:"name"`
	Type          string `json:"type" db:"type"`
	Version       string `json:"version" db:"version"`
	Year          int    `json:"year" db:"year"`
	FullScore     int    `json:"fullScore" db:"fullScore"`
	Choice        string `json:"choice" db:"choice"`
	Blank         string `json:"blank" db:"blank"`
	ImageFileName string `json:"-" db:"imageFileName"`
	ImageURL      string `json:"imageURL"`
}

// GetPapersByPaperID 根据试卷识别码获取试卷信息
func GetPapersByPaperID(paperIDs []string) []PaperDetailType {
	URL_PREFIX := conf.AppConfig.PaperImagesURL

	db := GetDB()
	result := []PaperDetailType{}
	for _, id := range paperIDs {
		paper := PaperDetailType{}
		err := db.Get(&paper, `SELECT paperID, uploadTime, name, type, version, year, fullScore, choice, blank, imageFileName FROM papers
						       WHERE paperID = ?;`, id)
		if err != nil {
			log.Printf("getting paper %s failed, err: %v \n", id, err)
			continue
		}
		result = append(result, paper)
	}

	for i, p := range result {
		result[i].ImageURL = URL_PREFIX + p.ImageFileName
		// 截取日期
		result[i].Time = p.Time[:10]
	}
	return result
}
