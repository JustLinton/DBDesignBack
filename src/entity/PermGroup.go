package entity

import "github.com/jinzhu/gorm"

type PermGroup struct {
	PGID int
	Name string
}

func InitPermGroup(db *gorm.DB){
	db.AutoMigrate(&PermGroup{})

	var uu[]PermGroup
	db.Find(&uu, "pg_id=?", 0)
	if len(uu)==0 {
		tmp := PermGroup{0,"业主用户"}
		db.Create(tmp)
		tmp = PermGroup{1,"超级管理员"}
		db.Create(tmp)
		tmp = PermGroup{2,"水务员"}
		db.Create(tmp)
	}
}
