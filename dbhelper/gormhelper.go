package dbhelper

import (
	"database/sql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

var gormEngine *gorm.DB

func GormConnection() *gorm.DB {
	if gormEngine == nil {
		host := LoadIni("DataBase", "db_host")
		user := LoadIni("DataBase", "db_user")
		password := LoadIni("DataBase", "db_password")
		name := LoadIni("DataBase", "db_name")
		//prefix := LoadIni("DataBase", "table_prefix")
		// 连接数据库
		source := user + ":" + password + "@tcp(" + host + ")/" + name + "?charset=utf8mb4&parseTime=True&loc=Local"
		sqlDB, err := sql.Open("mysql", source)
		if err != nil {
			log.Println("数据库连接出错 sqlDb: ", err.Error())
		}
		// 设置连接最长使用时间
		sqlDB.SetConnMaxLifetime(time.Duration(3600) * time.Second)
		//engine.DB().SetMaxOpenConns(3)

		engine, err := gorm.Open(mysql.New(mysql.Config{
			Conn: sqlDB,
		}), &gorm.Config{})
		if err != nil {
			log.Println("数据库连接出错 gorm: ", err.Error())
		}
		gormEngine = engine
	}
	return gormEngine
}
