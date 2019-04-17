package staffs

import (
	"log"
	"net/http"

	// "LearnServer/models/userDB"
	// "LearnServer/utils"
	"LearnServer/models/userDB"
	"LearnServer/utils"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

func addPaperHandler(c echo.Context) error {

	type addPaperType struct {
		SchoolID string `json:"schoolID"`
		Grade    string `json:"grade"`
		Class    int    `json:"class"`
		PaperID  string `json:"paperID"`
	}

	var paper addPaperType

	if err := c.Bind(&paper); err != nil {
		return utils.InvalidParams("invalid input, err: " + err.Error())
	}

	if paper.Class != 0 {
		_, err := userDB.C("classes").Upsert(bson.M{
			"schoolID": bson.ObjectIdHex(paper.SchoolID),
			"grade":    paper.Grade,
			"class":    paper.Class,
			"valid":    true,
		}, bson.M{
			"$addToSet": bson.M{
				"papers": paper.PaperID,
			},
		})

		if err != nil {
			log.Printf("add paper failed, paper: %v, err: %v\n", paper, err)
			return err
		}
	} else {
		// 所有班级
		_, err := userDB.C("classes").UpdateAll(bson.M{
			"schoolID": bson.ObjectIdHex(paper.SchoolID),
			"grade":    paper.Grade,
			"valid":    true,
		}, bson.M{
			"$addToSet": bson.M{
				"papers": paper.PaperID,
			},
		})

		if err != nil {
			log.Printf("add paper failed, paper: %v, err: %v\n", paper, err)
			return err
		}
	}

	return c.JSON(http.StatusOK, "Successfully add a paper")
}
