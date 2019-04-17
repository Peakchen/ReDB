package problemSorting

import "time"

type typeInfo struct {
	Type     string `db:"typename"`
	Category string `db:"category"`
	Priority int    `db:"priority"`
	Chapter  int    `db:"typeChapter"`
	Section  int    `db:"typeSection"`
}

type problemWithTypeTimeAndCorrectDB struct {
	TypeInfo  typeInfo
	ProblemID string    `bson:"problemID" json:"problemID"`
	SubIdx    int       `bson:"subIdx" json:"subIdx"`
	Correct   bool      `bson:"correct" json:"isCorrect"`
	Time      time.Time `bson:"assignDate" json:"assignDate"`
}
