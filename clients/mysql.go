package clients

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type MySqlConf struct {
	Host           string `json:"host" mapstructure:"host"`
	Port           int    `json:"port" mapstructure:"port"`
	ConnectTimeout int    `json:"connectTimeout" mapstructure:"connectTimeout"`
	DBName         string `json:"dbName" mapstructure:"dbName"`
	User           string `json:"user" mapstructure:"user"`
	Password       string `json:"password" mapstructure:"password"`
	SlowThreshold  int    `json:"slowThreshold" mapstructure:"slowThreshold"`
	LogLevel       string `json:"logLevel" mapstructure:"logLevel"`
}

func NewMySqlDB(cfg MySqlConf) *gorm.DB {
	var (
		err      error
		dsn      string
		db       *gorm.DB
		logLevel logger.LogLevel
	)
	dsn = fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%ds",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.ConnectTimeout,
	)
	//connect mysql db
	switch cfg.LogLevel {
	case "info", "INFO":
		logLevel = logger.Info
	case "warn", "WARN":
		logLevel = logger.Warn
	case "error", "ERROR":
		logLevel = logger.Error
	default:
		logLevel = logger.Silent
	}
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Duration(cfg.SlowThreshold) * time.Millisecond,
				LogLevel:                  logLevel,
				IgnoreRecordNotFoundError: false,
				Colorful:                  true,
			},
		),
	})
	if err != nil {
		panic(fmt.Errorf("connect mysql db failed, err: %s", err.Error()))
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(15)
	sqlDB.SetConnMaxIdleTime(time.Hour)
	sqlDB.SetConnMaxIdleTime(20 * time.Minute)
	return db
}
