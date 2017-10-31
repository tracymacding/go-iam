package mysql

import (
	"github.com/go-iam/db"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"strings"
)

type MysqlDriver struct{}

func (msd *MysqlDriver) Open(args ...interface{}) (db.Service, error) {

	dn := ""
	for _, s := range args {
		dn = dn + s.(string) + ","
	}
	dn = strings.Trim(dn, ",")

	database, err := gorm.Open("mysql", dn)
	if err != nil {
		return nil, err
	}
	// database = database.AutoMigrate(
	// 	&ID{},
	// 	&MysqlBucket{},
	// 	&MysqlObject{},
	// 	&HistoryObject{},
	// 	&MysqlUploadInfo{},
	// 	&db.UploadPart{},
	// 	&AbandonUploadPart{},
	// )
	// database = database.Set("gorm:table_options", "CHARSET=utf8 COLLATE=utf8_bin ROW_FORMAT=DYNAMIC")
	// err = database.AutoMigrate(&ObjectList{}).Error
	// if err != nil {
	// 	return nil, err
	// }
	return &mysqlService{dn, database}, nil
}

func init() {
	db.RegisterDriver("mysql", &MysqlDriver{})
}
