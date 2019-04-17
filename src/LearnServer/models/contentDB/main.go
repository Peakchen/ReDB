package contentDB

import (
	"log"
	"time"

	//"LearnServer/conf"
	"LearnServer/conf"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var _DB *sqlx.DB

// GetDB 获取内容数据库
func GetDB() *sqlx.DB {
	if _DB != nil {
		return _DB
	}
	cfg := mysql.Config{
		User:                 conf.AppConfig.ContentDB.User,
		Passwd:               conf.AppConfig.ContentDB.Password,
		DBName:               conf.AppConfig.ContentDB.Name,
		Net:                  conf.AppConfig.ContentDB.Protocol,
		Addr:                 conf.AppConfig.ContentDB.Addr,
		AllowNativePasswords: true,
		ParseTime:            true,
	}
	var err error
	// 确保 _DB 修改的是全局变量而不是 GetDB 中的局部变量
	_DB, err = sqlx.Connect("mysql", cfg.FormatDSN())
	if err != nil {
		panic(err)
	}
	if err := _DB.Ping(); err != nil {
		panic(err)
	}

	return _DB
}

func init() {
	db := GetDB()

	// 每10分钟心跳检测，确保数据库不会因为超时断开
	ticker := time.NewTicker(time.Minute * 10)
	go func() {
		for range ticker.C {
			if err := db.Ping(); err != nil {
				log.Printf("Ping failed. Error: %v", err)
			}
		}
	}()
}
