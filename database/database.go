package database

import (
	"encoding/json"
	"fmt"
	"github.com/micrease/micrease-core/config"
	log "github.com/micro/go-micro/v2/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

func Connect(config config.DatabaseSection) *gorm.DB {
	if len(config.Driver) == 0 {
		config.Driver = "mysql"
	}
	conf := &gorm.Config{}
	if config.Debug {
		conf.Logger = logger.Default.LogMode(logger.Info)
	}

	if db, err := gorm.Open(mysql.Open(config.DataSourceName), conf); err != nil {
		log.Fatal("连接数据库失败:"+config.DataSourceName, err)
		panic(err)
	} else {

		sqlDb, err := db.DB()
		if err != nil {
			log.Fatal("连接数据库失败:"+config.DataSourceName, err)
		}

		if config.MaxIdleConns > 0 {
			sqlDb.SetMaxIdleConns(int(config.MaxIdleConns))
		}

		if config.MaxOpenConns > 0 {
			sqlDb.SetMaxOpenConns(int(config.MaxOpenConns))
		}

		sqlDb.SetConnMaxLifetime(10 * time.Second)
		if config.Debug {
			data, _ := json.Marshal(sqlDb.Stats())
			fmt.Println(string(data))
		}
		log.Info("连接数据库成功:" + config.DataSourceName)

		return db
	}
}
