package dbhelper

import (
	"fmt"
	"log"
	"testing"
	"time"

	"database/sql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestGormConnection(t *testing.T) {
	// 测试数据库连接
	db := GormConnectionTest()
	result := map[string]interface{}{}
	db.Table("users").Take(&result)
	fmt.Println(result)
}

func GormConnectionTest() *gorm.DB {
	if gormEngine == nil {
		host := "127.0.0.1:3306"
		user := "root"
		password := "12345"
		name := "test"
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
