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

type AppConfig struct {
	Port        string
	DatabaseURL string
}

func InitConfig() {
	// 从环境变量中读取数据库URL
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		// 如果没有设置环境变量，使用默认的本地数据库连接
		databaseURL = "root:123456@tcp(localhost:3306)/sl2?charset=utf8&parseTime=True&loc=Local"
		fmt.Println("DATABASE_URL environment variable is not set, using default:", databaseURL)
	}

	// 初始化应用配置
	Config = &AppConfig{
		Port:        "8081", // 默认端口
		DatabaseURL: databaseURL,
	}

	// 初始化数据库连接
	var err error
	DB, err = gorm.Open("mysql", databaseURL)

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
