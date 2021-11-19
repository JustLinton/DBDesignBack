package entity

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

type Permission struct {
	PermID int
	Name string
}

func InitPermission(db *gorm.DB){
	db.AutoMigrate(&Permission{})

	var uu[]Permission
	db.Find(&uu, "perm_id=?", 100)
	if len(uu)==0 {
		tmp := Permission{100,"dashboard.user"}
		db.Create(tmp)
		tmp = Permission{101,"dashboard.waterstaff"}
		db.Create(tmp)
		tmp = Permission{102,"dashboard.gasstaff"}
		db.Create(tmp)

		tmp = Permission{200,"function.essentials.dashboard"}
		db.Create(tmp)
		//--
		tmp = Permission{201,"function.essentials.report"}
		db.Create(tmp)
		//--
		tmp = Permission{202,"function.essentials.things"}
		db.Create(tmp)
		//--

		tmp = Permission{203,"function.waterman.root"}
		db.Create(tmp)
		tmp = Permission{204,"function.waterman.rec"}
		db.Create(tmp)
		tmp = Permission{205,"function.waterman.userman"}
		db.Create(tmp)

		tmp = Permission{206,"function.essentials.nmap"}
		db.Create(tmp)
		//--
		tmp = Permission{207,"function.userman.root"}
		db.Create(tmp)
		//--
		tmp = Permission{208,"function.userman.overview"}
		db.Create(tmp)
		//--
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
		var ggpp[]GroupPermission
		db.Where("perm_id = ? AND pg_id = ?",permID,uu[0].PGID).Find(&ggpp)

		if len(ggpp)!=0{
			//have this permission
			c.String(http.StatusOK, fmt.Sprintf("ok"))
		}else{
			c.String(http.StatusOK, fmt.Sprintf("nok"))
		}

	})

}