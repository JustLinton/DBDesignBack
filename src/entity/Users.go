package entity

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

// User UserInfo 用户信息
type User struct {
	ID string
	Name string
	Phone string
	Email string
	PassSHA string
	PGID int
}

func InitUsers(db *gorm.DB){
	db.AutoMigrate(&User{})
}

func InitUsersApi(err *error,db *gorm.DB,router *gin.Engine) {
	router.GET("/profile", func(c *gin.Context) {
		//get user token
		usrToken,cerr := c.Cookie("user_token")
		if cerr!= nil{
			fmt.Printf("cookie not found\n")
			c.String(http.StatusOK,fmt.Sprintf("notlogged"))
			return
		}

		verbose := c.DefaultQuery("verbose","true")

		type Result struct {
			Name  string
			Email string
		}

		type ResultVerbose struct {
			Name  string
			Email string
		}

		if usrToken !="-1"{
			var uu[]User
			db.Find(&uu, "id=?", usrToken)
			//fmt.Printf("len: %d",len(uu) )

			if len(uu)!=0 {
				//found.
				if verbose=="false"{
					c.JSON(http.StatusOK,Result{
						Name: uu[0].Name,
						Email: uu[0].Email,
					})
				}else{
					c.JSON(http.StatusOK,ResultVerbose{
						Name: uu[0].Name,
					})
				}
			}else{
				//fmt.Printf("not found\n" )
				c.String(http.StatusOK,fmt.Sprintf("not found"))
			}
		}else{
			c.String(http.StatusOK,fmt.Sprintf("notlogged"))
		}

	})

	router.GET("/checkLoggedIn", func(c *gin.Context) {
		cookie,cerr := c.Cookie("user_token")
		if cerr!= nil{
			fmt.Printf("cookie not found\n")
			c.String(http.StatusOK,fmt.Sprintf("notlogged"))
			return
		}
		//fmt.Printf("%s", cookie)
		if cookie!="-1"{
			c.String(http.StatusOK,fmt.Sprintf("ok"))
		}else{
			c.String(http.StatusOK,fmt.Sprintf("notlogged"))
		}

	})

	router.GET("/logout", func(c *gin.Context) {
		c.SetCookie("user_token", "-1", 1000, "/", "localhost", false, true)
		c.Status(http.StatusOK)
	})

	router.POST("/login", func(c *gin.Context) {

		passwd := c.DefaultPostForm("passwd","nil")
		phone := c.DefaultPostForm("phone","nil")

		if passwd=="nil"||phone=="nil" {
			c.String(http.StatusNotAcceptable, fmt.Sprintln("network"))
		}

		//sha256 check
		//passwdBYTE:=[]byte(passwd)
		//passwdSHA:=GetSHA256HashCode(passwdBYTE)
		//fmt.Println(passwdSHA)

		var uu = new(User)
		db.Find(&uu, "phone=?", phone)

		if uu.PassSHA==passwd {
			//fmt.Printf("pass_sha:%s , pass:%s", uu.PassSHA,passwd)
			//passwd is correct
			//login success
			c.SetCookie("user_token", uu.ID, 1000, "/", "localhost", false, true)
			c.String(http.StatusOK,fmt.Sprintf("ok"))

		}else {
			//not correct passwd!
			c.String(http.StatusOK, fmt.Sprintf("passwd"))
		}
	})

	router.POST("/register", func(c *gin.Context) {

		email := c.DefaultPostForm("email","nil")
		passwd := c.DefaultPostForm("passwd","nil")
		name := c.DefaultPostForm("name","nil")
		phone := c.DefaultPostForm("phone","nil")

		if email=="nil"||passwd=="nil"||phone=="nil"||name=="nil" {
			c.String(http.StatusNotAcceptable, fmt.Sprintln("network"))
		}

		//sha256 check
		//passwdBYTE:=[]byte(passwd)
		//passwdSHA:=GetSHA256HashCode(passwdBYTE)
		//fmt.Println(passwdSHA)

		userUid := uuid.Must(uuid.NewV4(),*err)

		var uu = new(User)
		db.Find(&uu, "phone=?", phone)
		//fmt.Printf("phone:%s", uu.Phone)
		if uu.Phone!=phone {
			//phone num haven't used yet
			fmt.Printf("uuid:%s\n", userUid)
			strUsrUid := fmt.Sprintf("%s",userUid)
			newUser := User{strUsrUid, name, phone, email,passwd,0}
			db.Create(newUser)
			//register success
			c.String(http.StatusOK,fmt.Sprintf("ok"))
			//c.SetCookie("user_token", strUsrUid, 1000, "/", "localhost", false, true)
		}else {
			//phone num is used!
			c.String(http.StatusOK, fmt.Sprintf("phone"))
		}

	})
}