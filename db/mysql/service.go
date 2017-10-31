package mysql

import (
	"github.com/jinzhu/gorm"
)

type mysqlService struct {
	dn string
	DB *gorm.DB
}

func (ms *mysqlService) AllocateID() (result int64, err error) {

	reuslt, err := ms.DB.DB().Exec("insert into ids set value = 1")
	if err != nil {
		return 0, err
	}

	return reuslt.LastInsertId()
}
