package entity

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

type Permssion struct {
	PermID int
	Name string
}

func InitPermission(db *gorm.DB){
	db.AutoMigrate(&Permssion{})

	var uu[]Permssion
	db.Find(&uu, "perm_id=?", 0)
	if len(uu)==0 {
		tmp := Permssion{100,"dashboard.user"}
		db.Create(tmp)
		tmp = Permssion{101,"dashboard.waterstaff"}
		db.Create(tmp)
		tmp = Permssion{102,"dashboard.gasstaff"}
		db.Create(tmp)
		tmp = Permssion{200,"userman.root"}
		db.Create(tmp)
		tmp = Permssion{300,"waterman.rec"}
		db.Create(tmp)
	}
}

func InitPermissionsApi(err *error,db *gorm.DB,router *gin.Engine) {

	router.GET("/haveperm", func(c *gin.Context) {
		userToken,cerr := c.Cookie("user_token")
		permID := c.Query("permid")
		fmt.Printf("out: %s",permID )

		if cerr!= nil{
			c.Status(http.StatusBadRequest)
			fmt.Printf("haveperm : cookie not found\n")
			return
		}

		var uu[]User
		db.Find(&uu, "id=?", userToken)

		if len(uu)==0{
			c.Status(http.StatusBadRequest)
			fmt.Printf("haveperm : user not found\n")
			return
		}

		//found.
		var ggpp[]GroupPermssion
		db.Where("perm_id = ? AND pg_id = ?",permID,uu[0].PGID).Find(&ggpp)

		if len(ggpp)!=0{
			//have this permission
			c.String(http.StatusOK, fmt.Sprintf("ok"))
		}else{
			c.String(http.StatusOK, fmt.Sprintf("nok"))
		}

	})

}