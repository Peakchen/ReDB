package students

type problem struct {
	ProblemID string `json:"problemID" bson:"problemID"`
	SubIdx    int    `json:"subIdx" bson:"subIdx"`
}
