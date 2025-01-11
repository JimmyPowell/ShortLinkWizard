package config

import (
	"fmt"
	"os"

	"shortlink/model"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	DB     *gorm.DB
	Config *AppConfig
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type AppConfig struct {
	Port string
	DB   DatabaseConfig
}

func InitConfig() {
	// 初始化应用配置
	Config = &AppConfig{
		Port: "8081", // 默认端口
		DB: DatabaseConfig{
			Host:     "localhost",
			Port:     "3306",
			User:     "root",
			Password: "123456",
			DBName:   "sl2",
		},
	}

	// 初始化数据库连接
	var err error
	DB, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		Config.DB.User,
		Config.DB.Password,
		Config.DB.Host,
		Config.DB.Port,
		Config.DB.DBName))

	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		os.Exit(1)
	}

	DB.LogMode(true)

	// 自动迁移数据库
	if err := model.AutoMigrate(DB); err != nil {
		fmt.Println("Failed to migrate database:", err)
		os.Exit(1)
	}
}

// Get 返回全局配置
func Get() *AppConfig {
	return Config
}
