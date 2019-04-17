package staffs

import (
	"fmt"
	"net/http"
	"strconv"

	"LearnServer/models/contentDB"
	"LearnServer/models/userDB"

	"LearnServer/utils"
	mapset "github.com/deckarep/golang-set"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

func getClassBooks(schoolID string, grade string, class int) ([]contentDB.BookDetailType, error) {
	type booksType struct {
		BookIDs []interface{} `bson:"books"`
	}
	classBookIDs := []booksType{}
	var selector bson.M
	if class != 0 {
		selector = bson.M{
			"schoolID": bson.ObjectIdHex(schoolID),
			"grade":    grade,
			"class":    class,
			"valid":    true,
		}
	} else {
		// 选择全部班级
		selector = bson.M{
			"schoolID": bson.ObjectIdHex(schoolID),
			"grade":    grade,
			"valid":    true,
		}
	}
	err := userDB.C("classes").Find(selector).All(&classBookIDs)
	if err != nil {
		return nil, err
	}
	if len(classBookIDs) == 0 {
		return nil, fmt.Errorf("not found")
	}

	// bookIDResults 得到 bookID 交集的结果的 slice
	bookIDResults := []string{}
	// 存储 bookID 的交集
	bookIDSets := mapset.NewSetFromSlice(classBookIDs[0].BookIDs)
	for _, b := range classBookIDs {
		bookIDSets = bookIDSets.Intersect(mapset.NewSetFromSlice(b.BookIDs))
	}
	it := bookIDSets.Iterator()
	for bookID := range it.C {
		bookIDResults = append(bookIDResults, bookID.(string))
	}

	return contentDB.GetBooksByBookID(bookIDResults), nil
}

func getBooksOfClassHandler(c echo.Context) error {
	schoolID := c.QueryParam("schoolID")
	grade := c.QueryParam("grade")
	class, err := strconv.Atoi(c.QueryParam("class"))
	if err != nil {
		return utils.InvalidParams("class is not a number!")
	}
	books, err := getClassBooks(schoolID, grade, class)
	if err != nil {
		if err.Error() == "not found" {
			return utils.NotFound("this class has no books")
		}
		return err
	}
	return c.JSON(http.StatusOK, books)
}

func deleteBookFromClassHandler(c echo.Context) error {
	type inputType struct {
		SchoolID string `json:"schoolID"`
		Grade    string `json:"grade"`
		Class    int    `json:"class"`
		BookID   string `json:"bookID"`
	}

	input := inputType{}
	if err := c.Bind(&input); err != nil {
		return utils.InvalidParams("invalid input!" + err.Error())
	}

	var selector bson.M
	if input.Class != 0 {
		selector = bson.M{
			"schoolID": bson.ObjectIdHex(input.SchoolID),
			"grade":    input.Grade,
			"class":    input.Class,
			"valid":    true,
		}
	} else {
		// 选择全部班级
		selector = bson.M{
			"schoolID": bson.ObjectIdHex(input.SchoolID),
			"grade":    input.Grade,
			"valid":    true,
		}
	}

	_, err := userDB.C("classes").UpdateAll(selector, bson.M{
		"$pull": bson.M{
			"books": input.BookID,
		},
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "successfully deleted!")
}
