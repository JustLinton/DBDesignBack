package entity

import "github.com/jinzhu/gorm"

type GroupPermssion struct {
	PermID int
	PGID int
}

func InitGroupPermission(db *gorm.DB){
	db.AutoMigrate(&GroupPermssion{})
	var uu[]GroupPermssion
	db.Find(&uu, "perm_id=?", 0)
	if len(uu)==0 {
		//group user(0)
		tmp := GroupPermssion{100,0}
		db.Create(tmp)
		//group waterStaff(2)
		tmp = GroupPermssion{101,2}
		db.Create(tmp)
		tmp = GroupPermssion{300,2}
		db.Create(tmp)
	}
}