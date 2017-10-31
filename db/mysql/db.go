package mysql

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	DB *gorm.DB
)

type DBConfig struct {
	Driver string
	DN     string
}

func Init(config *DBConfig) error {

	database, err := gorm.Open(config.Driver, config.DN)
	if err != nil {
		return err
	}
	database = database.AutoMigrate(
		&ID{},
		&MysqlBucket{},
		&MysqlObject{},
		&HistoryObject{},
		&MysqlUploadInfo{},
		&AbandonUploadPart{},
	)
	database = database.Set("gorm:table_options", "CHARSET=utf8 COLLATE=utf8_bin ROW_FORMAT=DYNAMIC")
	err = database.AutoMigrate(&ObjectList{}).Error
	if err != nil {
		return err
	}
	DB = database
	return nil
}
