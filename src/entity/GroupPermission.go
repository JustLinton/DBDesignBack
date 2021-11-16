package entity

import "github.com/jinzhu/gorm"

type GroupPermission struct {
	PermID int
	PGID int
}

func InitGroupPermission(db *gorm.DB){
	db.AutoMigrate(&GroupPermission{})
	var uu[]GroupPermission
	db.Find(&uu, "perm_id=?", 100)
	if len(uu)==0 {
		tmp := GroupPermission{100,0}
		db.Create(tmp)
		//--
		tmp = GroupPermission{200,0}
		db.Create(tmp)

		tmp = GroupPermission{101,2}
		db.Create(tmp)
		//--
		tmp = GroupPermission{200,2}
		db.Create(tmp)
		//--
		tmp = GroupPermission{201,2}
		db.Create(tmp)
		//--
		tmp = GroupPermission{202,2}
		db.Create(tmp)
		//--
		tmp = GroupPermission{203,2}
		db.Create(tmp)
		tmp = GroupPermission{204,2}
		db.Create(tmp)
		tmp = GroupPermission{205,2}
		db.Create(tmp)
	}
}