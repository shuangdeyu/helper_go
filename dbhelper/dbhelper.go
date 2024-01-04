package dbhelper

import (
	"log"
	"time"

	"github.com/Unknwon/goconfig"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/redis.v5"
	"xorm.io/core"
	"xorm.io/xorm"
)

var (
	NewEngine          *xorm.Engine
	NewSession         *xorm.Session
	RedisClient        *redis.Client
	RedisClusterClient *redis.ClusterClient
)

type ConfInit struct {
	FilePath string
}

var (
	confFilePath *ConfInit
	iniCfg       *goconfig.ConfigFile
)

// 配置文件初始化
func (p *ConfInit) FileInit() {
	if confFilePath == nil {
		confFilePath = p
	}

	var tmpErr error
	iniCfg, tmpErr = goconfig.LoadConfigFile(confFilePath.FilePath)
	if tmpErr != nil {
		panic("读取配置文件失败:" + tmpErr.Error())
	}
}

// 读取配置文件
func InitConfig(path string) {
	var tmpErr error
	iniCfg, tmpErr = goconfig.LoadConfigFile(path)
	if tmpErr != nil {
		panic("读取配置文件失败: " + tmpErr.Error())
	}
}

// 加载配置文件
func LoadIni(param1 string, param2 string) string {
	result, err := iniCfg.GetValue(param1, param2)
	if err != nil {
		log.Fatal("无法获取键值", err)
		//log.Println("无法获取键值", err)
		return ""
	}
	return result
}

// 数据库连接初始化
func NewEngineInit() *xorm.Engine {
	if NewEngine == nil {
		host := LoadIni("DataBase", "db_host")
		user := LoadIni("DataBase", "db_user")
		password := LoadIni("DataBase", "db_password")
		name := LoadIni("DataBase", "db_name")
		prefix := LoadIni("DataBase", "table_prefix")
		// 连接数据库
		source := user + ":" + password + "@tcp(" + host + ")/" + name + "?charset=utf8"
		engine, err := xorm.NewEngine("mysql", source)
		if err != nil {
			log.Println("数据库连接出错: ", err.Error())
		}
		// 设置连接最长使用时间
		engine.SetConnMaxLifetime(time.Duration(3600) * time.Second)
		//engine.DB().SetMaxOpenConns(3)
		// 设置表前缀映射
		if prefix != "" {
			tbMapper := core.NewPrefixMapper(core.SnakeMapper{}, prefix)
			engine.SetTableMapper(tbMapper)
		}
		NewEngine = engine
	}
	return NewEngine
}

// session
func NewSessionInit() *xorm.Session {
	if NewSession == nil {
		NewSession = NewEngineInit().NewSession()
	}
	return NewSession
}
