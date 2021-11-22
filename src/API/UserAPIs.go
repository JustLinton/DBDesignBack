package API

import (
	"db_design/src/entity"
	"db_design/src/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

func InitUsersApi2(err *error,db *gorm.DB,router *gin.Engine) {

	router.GET("/userinfo", func(c *gin.Context) {

		//permission number needed to use the API
		var permID=209

		//get user token
		usrToken,cerr := c.Cookie("user_token")
		if cerr!= nil{
			fmt.Printf("cookie not found\n")
			c.String(http.StatusOK,fmt.Sprintf("notlogged"))
			return
		}

		//check if have the permission to use the API.
		if !utils.CheckPerm(permID,usrToken,db){
			c.String(http.StatusOK, fmt.Sprintf("permission"))
			return
		}

		uid := c.Query("uid")
		type Result struct {
			Name  string
			Email string
			IdCard string
			PGname string
			Phone string
			Uid string
		}

		if usrToken !="-1"{
			var uu[]entity.User
			db.Find(&uu, "id=?", uid)
			//fmt.Printf("len: %d",len(uu) )

			if len(uu)==0{
				c.Status(http.StatusBadRequest)
				fmt.Printf("/userinfo : user not found\n")
				return
			}

			//look for pg name
			var pgname string
			var ppgg[]entity.PermGroup
			db.Find(&ppgg, "pg_id=?", uu[0].PGID)
			if len(ppgg)==0{
				pgname=""
			}else {
				pgname=ppgg[0].Name
			}

			//found.
			c.JSON(http.StatusOK,Result{
				Name: uu[0].Name,
				Email: uu[0].Email,
				PGname: pgname,
				Phone: uu[0].Phone,
				Uid: uu[0].ID,
				IdCard: uu[0].IDCard,
			})

		}else{
			c.String(http.StatusOK,fmt.Sprintf("notlogged"))
		}

	})


	router.GET("/userlist", func(c *gin.Context) {

		//permission number needed to use the API
		var permID=208

		//get user token
		usrToken,cerr := c.Cookie("user_token")
		if cerr!= nil{
			fmt.Printf("cookie not found\n")
			c.String(http.StatusOK,fmt.Sprintf("notlogged"))
			return
		}

		//check if have the permission to use the API.
		if !utils.CheckPerm(permID,usrToken,db){
			c.String(http.StatusOK, fmt.Sprintf("permission"))
			return
		}

		//page := c.Query("page")
		//pageLen := c.Query("pageLen")
		var tmpRow map[string]string
		tmpRow = make(map[string]string)
		var rows []map[string]string

		type Result struct {
			Rows []map[string]string
		}

		if usrToken !="-1"{
			var uu[]entity.User
			db.Find(&uu)
			//fmt.Printf("len: %d",len(uu) )

			for i:=0;i<len(uu);i++{
				theUsr:=uu[i]

				tmpRow["name"]=theUsr.Name
				tmpRow["idcard"]=theUsr.IDCard
				tmpRow["phone"]=theUsr.Phone
				tmpRow["email"]=theUsr.Email
				tmpRow["uid"]=theUsr.ID

				//look for pg name
				var ppgg[]entity.PermGroup
				db.Find(&ppgg, "pg_id=?", theUsr.PGID)
				if len(ppgg)==0{
					tmpRow["pgname"]=""
				}else {
					tmpRow["pgname"]=ppgg[0].Name
				}

				//fake data generator
				//for j:=0;j<=30;j++{
				//	tmpRow["name"]=strconv.Itoa(j+i*100)
				//	tmpRow["idcard"]=theUsr.IDCard
				//	tmpRow["phone"]=theUsr.Phone
				//	tmpRow["email"]=theUsr.Email
				//	rows = append(rows, tmpRow)
				//	tmpRow = make(map[string]string)
				//}

				rows = append(rows, tmpRow)
				tmpRow = make(map[string]string)
			}

			c.JSON(http.StatusOK,Result{
				Rows: rows,
			})

		}else{
			c.String(http.StatusOK,fmt.Sprintf("notlogged"))
		}

	})
}