package gormsample

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"testing"
	"time"
)

type STAdExpFlow struct {
	ID	int32 `gorm:"column:id; primary_key"`
	APK string `gorm:"column:apk"`
	AdSource int32 `gorm:"column:adsource"`
	CreateTime time.Time `gorm:"column:create_time"`
	DeleteFlag int32 `gorm:"column:delete_flag"`
}

func (*STAdExpFlow) TableName() string {
	return "t_ad_exp_flow"
}

func TestFirst(t *testing.T) {
	var err error
	dbConnArgs := fmt.Sprintf("ivanadmin:123456@(127.0.0.1:3381)/db_ad_monitor?charset=utf8&parseTime=True&loc=Local")
	DB, err := gorm.Open("mysql", dbConnArgs)
	if err != nil {
		fmt.Printf("fail connect:%s", err.Error())
		panic("failed to connect database")
	}
	DB.DB().SetMaxIdleConns(10)
	DB.DB().SetMaxOpenConns(100)

	var flowList []STAdExpFlow
	//err := db.Order("create_time asc").Where("create_time>?", lastTime).Find(&flowList).Error
	err = DB.Debug().Find(&flowList).Error
	if err != nil {
		fmt.Printf("GetTaskByLastTime|db find err:%s", err.Error())
	}
}