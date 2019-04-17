package students

import (
	"net/http"
	"time"

	"LearnServer/conf"
	"LearnServer/models/userDB"
	"LearnServer/utils"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

type TStudentUser struct {
	ID       int64  `json:"learnID"`
	Password string `json:"password"`
	Remember bool   `json:"remember"`
}

func createJwtCookie(id bson.ObjectId, remember bool) (http.Cookie, error) {
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	now := time.Now()
	exp := now.AddDate(0, 1, 0) // 1 month
	claims["iat"] = now.Unix()
	claims["exp"] = exp.Unix()
	claims["id"] = id
	claims["userType"] = "students"

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(conf.AppConfig.Secret))
	if err != nil {
		return http.Cookie{}, err
	}
	tokCookie := http.Cookie{
		Name:     "TOKEN",
		Path:     "/api/v3/students/",
		Value:    t,
		HttpOnly: true,
	}
	if remember {
		tokCookie.MaxAge = int(exp.Sub(now) / time.Second)
	}
	return tokCookie, nil
}

func loginHandler(c echo.Context) error {
	u := TStudentUser{}
	if err := c.Bind(&u); err != nil {
		return utils.InvalidParams()
	}
	id := u.ID
	password := u.Password
	remember := u.Remember

	var objectID struct {
		ID bson.ObjectId `bson:"_id,omitempty"`
	}

	err := userDB.C("students").Find(bson.M{
		"learnID":  id,
		"password": password,
		"valid":    true,
	}).Select(bson.M{
		"_id": 1,
	}).One(&objectID)
	if err != nil {
		return utils.Unauthorized("用户名或者密码错误！")
	}

	var tokCookie http.Cookie
	tokCookie, err = createJwtCookie(objectID.ID, remember)
	if err != nil {
		return err
	}

	c.SetCookie(&tokCookie)
	return c.NoContent(http.StatusOK)
}
