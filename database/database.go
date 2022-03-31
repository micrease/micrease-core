package database

import (
	"github.com/jinzhu/gorm"
	log "github.com/micro/go-micro/v2/logger"
	"micrease-core/config"
	"time"
)

func Connect(config config.DatabaseSection) *gorm.DB {
	if len(config.Driver) == 0 {
		config.Driver = "mysql"
	}

	if db, err := gorm.Open(config.Driver, config.DataSourceName); err != nil {
		log.Fatal("连接数据库失败:"+config.DataSourceName, err)
		panic(err)
	} else {
		// 表名使用单数
		db.SingularTable(true)
		db.LogMode(true)

		if config.MaxOpenConns > 0 {
			db.DB().SetMaxIdleConns(int(config.MaxIdleConns))
		}

		if config.MaxOpenConns > 0 {
			db.DB().SetMaxOpenConns(int(config.MaxOpenConns))
		}

		db.DB().SetConnMaxLifetime(10 * time.Second)
		log.Info("连接数据库成功:" + config.DataSourceName)
		return db
	}
}
