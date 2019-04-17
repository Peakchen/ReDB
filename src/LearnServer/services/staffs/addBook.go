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

func addBookHandler(c echo.Context) error {

	type addBookType struct {
		SchoolID string `json:"schoolID"`
		Grade    string `json:"grade"`
		Class    int    `json:"class"`
		BookID   string `json:"bookID"`
	}

	var book addBookType

	if err := c.Bind(&book); err != nil {
		return utils.InvalidParams("invalid input, err: " + err.Error())
	}

	if book.Class != 0 {
		_, err := userDB.C("classes").Upsert(bson.M{
			"schoolID": bson.ObjectIdHex(book.SchoolID),
			"grade":    book.Grade,
			"class":    book.Class,
			"valid":    true,
		}, bson.M{
			"$addToSet": bson.M{
				"books": book.BookID,
			},
		})

		if err != nil {
			log.Printf("add book failed, book: %v, err: %v\n", book, err)
			return err
		}
	} else {
		// 添加书本给所有班级，此处不能用upsert
		_, err := userDB.C("classes").UpdateAll(bson.M{
			"schoolID": bson.ObjectIdHex(book.SchoolID),
			"grade":    book.Grade,
			"valid":    true,
		}, bson.M{
			"$addToSet": bson.M{
				"books": book.BookID,
			},
		})

		if err != nil {
			log.Printf("add book failed, book: %v, err: %v\n", book, err)
			return err
		}
	}

	return c.JSON(http.StatusOK, "Successfully add a book")
}
