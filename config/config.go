package config

import (
	"fmt"
	"log"
	"os"

	"shortlink/model"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
)

var (
	DB     *gorm.DB
	Config *AppConfig
)

type AppConfig struct {
	Port     string
	DBUser   string
	DBPass   string
	DBHost   string
	DBPort   string
	DBName   string
	DBParams string
}

// InitConfig 初始化配置
func InitConfig() {
	// 加载.env 文件（如果存在）
	err := godotenv.Load()
	if err != nil {
		log.Println("No.env file found. Using system environment variables.")
	}

	// 从环境变量中读取配置
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbParams := os.Getenv("DB_PARAMS")
	port := os.Getenv("PORT")

	if dbUser == "" {
		log.Fatal("DB_USER environment variable is not set")
	}
	if dbPass == "" {
		log.Fatal("DB_PASS environment variable is not set")
	}
	if dbHost == "" {
		log.Fatal("DB_HOST environment variable is not set")
	}
	if dbPort == "" {
		log.Fatal("DB_PORT environment variable is not set")
	}
	if dbName == "" {
		log.Fatal("DB_NAME environment variable is not set")
	}
	if dbParams == "" {
		log.Fatal("DB_PARAMS environment variable is not set")
	}
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}

	// 初始化应用配置
	Config = &AppConfig{
		Port:     port,
		DBUser:   dbUser,
		DBPass:   dbPass,
		DBHost:   dbHost,
		DBPort:   dbPort,
		DBName:   dbName,
		DBParams: dbParams,
	}

	databaseURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", dbUser, dbPass, dbHost, dbPort, dbName, dbParams)

	// 初始化数据库连接
	DB, err = gorm.Open("mysql", databaseURL)
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		os.Exit(1)
	}

	// 启用数据库日志模式（调试用）
	DB.LogMode(true)

	// 自动迁移数据库
	if err := model.AutoMigrate(DB); err != nil {
		fmt.Println("Failed to migrate database:", err)
		os.Exit(1)
	}

	fmt.Println("Configuration initialized successfully.")
}

// Get 返回全局配置
func Get() *AppConfig {
	return Config
}
