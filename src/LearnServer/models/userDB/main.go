package userDB

import (
	"fmt"
	"log"
	"net/url"

	"gopkg.in/mgo.v2"

	// "LearnServer/conf"
	"LearnServer/conf"
)

var (
	globalSession *mgo.Session
	gloablDB      *mgo.Database
)

// C get collection of the database
func C(name string) *mgo.Collection {
	return gloablDB.C(name)
}

func init() {
	host := conf.AppConfig.UserDB.HostAndPort
	if conf.AppConfig.UserDB.User != "" || conf.AppConfig.UserDB.Passwd != "" {
		host = fmt.Sprintf("%s:%s@%s",
			url.QueryEscape(conf.AppConfig.UserDB.User),
			url.QueryEscape(conf.AppConfig.UserDB.Passwd),
			host)
	}
	var err error
	globalSession, err = mgo.Dial(fmt.Sprintf("mongodb://%s/%s", host, conf.AppConfig.UserDB.Name))
	if err != nil {
		panic(fmt.Sprintf("Initialize mongodb error: %v", err))
	}
	if err := globalSession.Ping(); err != nil {
		panic(fmt.Sprintf("MongoDB ping error: %v", err))
	}
	gloablDB = globalSession.DB(conf.AppConfig.UserDB.Name)

	// Optional. Switch the session to a monotonic behavior.
	globalSession.SetMode(mgo.Monotonic, true)

	log.Println("MongoDB initialize success.")
	mgo.SetDebug(conf.AppConfig.Debug)
}
