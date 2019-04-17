package staffs

import (
	"net/http"
	"strconv"
	"time"

	"LearnServer/conf"
	"LearnServer/models/userDB"
	"LearnServer/services/staffs/validation"
	"LearnServer/utils"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

func testLoginHandler(c echo.Context) error {
	var id string
	if err := validation.ValidateStaff(c, &id); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, "ok")
}

type staffUser struct {
	ID       string `json:"staffID"`
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
	claims["userType"] = "staffs"

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(conf.AppConfig.Secret))
	if err != nil {
		return http.Cookie{}, err
	}
	tokCookie := http.Cookie{
		Name:     "TOKEN",
		Path:     "/api/v3/staffs/",
		Value:    t,
		HttpOnly: true,
	}
	if remember {
		tokCookie.MaxAge = int(exp.Sub(now) / time.Second)
	}
	return tokCookie, nil
}

func loginHandler(c echo.Context) error {
	u := staffUser{}
	if err := c.Bind(&u); err != nil {
		return utils.InvalidParams()
	}
	id, err := strconv.Atoi(u.ID)
	if err != nil {
		utils.Unauthorized("invalid staffID or password")
	}
	password := u.Password
	remember := u.Remember

	var objectID struct {
		ID bson.ObjectId `bson:"_id,omitempty"`
	}

	err = userDB.C("staffs").Find(bson.M{
		"staffID":  id,
		"password": password,
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
	return c.JSON(http.StatusOK, "Successfully logged in!")
}
