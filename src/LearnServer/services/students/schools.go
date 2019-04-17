package students

import (
	"log"
	"net/http"

	"LearnServer/models/userDB"
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
