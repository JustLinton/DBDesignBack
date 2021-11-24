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

	router.GET("/deluser", func(c *gin.Context) {

		//if the user doesn't pass the validation check
		if !utils.UserValidation(209,db,err,router,c){
			return
		}

		uid := c.Query("uid")
		var uu[] entity.User
		db.Find(&uu, "id=?", uid)

		if len(uu)==0{
			c.Status(http.StatusBadRequest)
			fmt.Printf("/deluser : user not found\n")
			return
		}

		db.Delete(uu[0])
		c.String(http.StatusOK,"ok")
	})

	router.POST("/userinfo", func(c *gin.Context) {

		//if the user doesn't pass the validation check
		if !utils.UserValidation(209,db,err,router,c){
			return
		}

		uid := c.PostForm("Uid")
		var uu[] entity.User
		db.Find(&uu, "id=?", uid)

		if len(uu)==0{
			c.Status(http.StatusBadRequest)
			fmt.Printf("/userinfo post : user not found\n")
			return
		}

		db.Delete(uu[0])
		uu[0].Name = c.DefaultPostForm("Name",uu[0].Name)
		uu[0].Email = c.DefaultPostForm("Email",uu[0].Email)
		uu[0].IDCard = c.DefaultPostForm("IdCard",uu[0].IDCard)
		pgname := c.PostForm("PGname")
		uu[0].Phone = c.DefaultPostForm("Phone",uu[0].Phone)
		uu[0].ID = c.DefaultPostForm("Uid",uu[0].ID)
		uu[0].Gender = c.DefaultPostForm("Gender",uu[0].Gender)
		pgid:=uu[0].PGID

		//use pgname to find the pgid
		var ppgg[]entity.PermGroup
		db.Find(&ppgg, "name=?", pgname)
		if len(ppgg)!=0{
			pgid=ppgg[0].PGID
		}
		uu[0].PGID=pgid

		db.Create(uu[0])
	})

	router.GET("/userinfo", func(c *gin.Context) {

		//permission number needed to use the API
		var permID=209

		//get user token and check if logged in
		usrToken,cerr := c.Cookie("user_token")
		if cerr!= nil{
			fmt.Printf("cookie not found\n")
			c.String(http.StatusOK,fmt.Sprintf("notlogged"))
			return
		}

		if usrToken =="-1"{
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
			PGNameList []string
			Gender string
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

			//list all pgnames.
			var ppggNameList []string
			//clear the array.
			ppgg=[]entity.PermGroup{}
			db.Find(&ppgg)
			for i:=0;i<len(ppgg);i++{
				ppggNameList = append(ppggNameList,ppgg[i].Name)
			}

			//found.
			c.JSON(http.StatusOK,Result{
				Name: uu[0].Name,
				Email: uu[0].Email,
				PGname: pgname,
				Phone: uu[0].Phone,
				Uid: uu[0].ID,
				IdCard: uu[0].IDCard,
				PGNameList:ppggNameList,
				Gender: uu[0].Gender,
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
				//	tmpRow["phone"]=strconv.Itoa(j+i*100)
				//	tmpRow["idcard"]=theUsr.IDCard
				//	tmpRow["name"]=strconv.Itoa(j+i*100)
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