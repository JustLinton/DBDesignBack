package utils

import (
	"db_design/src/entity"
	"github.com/jinzhu/gorm"
)

func CheckPerm(permID int,usrToken string,db *gorm.DB) bool{

	//check if have the permission to use the API
	var uuf []entity.User
	db.Find(&uuf, "id=?", usrToken)

	if len(uuf) == 0 {
		return false
	}

	//found.
	var ggppf []entity.GroupPermission
	db.Where("perm_id = ? AND pg_id = ?", permID, uuf[0].PGID).Find(&ggppf)

	if len(ggppf) != 0 {
		//have this permission
		return true
	} else {
		return false
	}

}
