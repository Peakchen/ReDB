package staffs

import (
	"log"
	"net/http"

	"LearnServer/models/userDB"
	"LearnServer/utils"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

func getSchoolsHandler(c echo.Context) error {
	province := c.QueryParam("province")
	city := c.QueryParam("city")
	district := c.QueryParam("district")
	county := c.QueryParam("county")

	query := make(bson.M)
	if province != "" {
		query["province"] = province
	}
	if city != "" {
		query["city"] = city
	}
	if district != "" {
		query["district"] = district
	}
	if county != "" {
		query["county"] = county
	}

	type schoolTypeDB struct {
		Name     string        `bson:"name"`
		SchoolID bson.ObjectId `bson:"_id"`
	}
	schoolsDB := []schoolTypeDB{}

	err := userDB.C("schools").Find(query).All(&schoolsDB)
	if err != nil {
		log.Printf("unable to get schools! %v", err)
		return err
	}

	type schoolType struct {
		Name     string `json:"name"`
		SchoolID string `json:"schoolID"`
	}
	schools := []schoolType{}

	for _, s := range schoolsDB {
		schools = append(schools, schoolType{
			Name:     s.Name,
			SchoolID: s.SchoolID.Hex(),
		})
	}

	return c.JSON(http.StatusOK, schools)
}

func uploadSchoolHandler(c echo.Context) error {
	type uploadSchoolType struct {
		Province string `json:"province" bson:"province"`
		City     string `json:"city" bson:"city"`
		District string `json:"district" bson:"district"`
		County   string `json:"county" bson:"county"`
		Name     string `json:"name" bson:"name"`
	}

	var uploadSchool uploadSchoolType
	if err := c.Bind(&uploadSchool); err != nil {
		return utils.InvalidParams("invalid input! " + err.Error())
	}

	count, err := userDB.C("schools").Find(bson.M{
		"province": uploadSchool.Province,
		"city":     uploadSchool.City,
		"district": uploadSchool.District,
		"county":   uploadSchool.County,
		"name":     uploadSchool.Name,
	}).Count()
	if err != nil {
		return err
	}

	if count != 0 {
		return utils.Forbidden("this school already exist!")
	}

	err = userDB.C("schools").Insert(uploadSchool)
	if err != nil {
		log.Printf("unable to save school! err: %v", err)
		return err
	}

	objID := struct {
		ID bson.ObjectId `bson:"_id"`
	}{}

	err = userDB.C("schools").Find(bson.M{
		"province": uploadSchool.Province,
		"city":     uploadSchool.City,
		"district": uploadSchool.District,
		"county":   uploadSchool.County,
		"name":     uploadSchool.Name,
	}).One(&objID)
	if err != nil {
		log.Printf("unable to get school! err: %v", err)
		return err
	}

	return c.JSON(http.StatusOK, bson.M{
		"schoolID": objID.ID.Hex(),
	})
}
