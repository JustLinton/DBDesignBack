package utils

import (
	"db_design/src/entity"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

func UserValidation(permID int,db *gorm.DB,err *error,router *gin.Engine,c *gin.Context) bool{

	//get user token
	usrToken,cerr := c.Cookie("user_token")

	//not logged
	if cerr!= nil{
		fmt.Printf("cookie not found\n")
		c.String(http.StatusOK,fmt.Sprintf("notlogged"))
		return false
	}

	//not logged
	if usrToken =="-1"{
		c.String(http.StatusOK,fmt.Sprintf("notlogged"))
		return false
	}

	//check if have the permission to use the API.
	if !CheckPerm(permID,usrToken,db){
		c.String(http.StatusOK, fmt.Sprintf("permission"))
		return false
	}

	return true
}

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
