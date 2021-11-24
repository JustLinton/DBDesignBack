package entity

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"strings"
)

// User UserInfo 用户信息
type User struct {
	ID string
	Name string
	Phone string
	Email string
	PassSHA string
	PGID int
	IDCard string
	Gender string
}

func InitUsers(db *gorm.DB){
	db.AutoMigrate(&User{})
}

func InitUsersApi(err *error,db *gorm.DB,router *gin.Engine) {

	//domain := "localhost"
	domain := "nesto.cupof.beer"

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
			PGID int
			PGName string
			SidebarC map[string]int
			DashboardC map[string]int
			FunctionC map[string]int
		}

		type ResultVerbose struct {
			Name  string
			Email string
		}

		if usrToken !="-1"{
			var uu[]User
			db.Find(&uu, "id=?", usrToken)
			//fmt.Printf("len: %d",len(uu) )

			if len(uu)==0{
				c.Status(http.StatusBadRequest)
				fmt.Printf("profile : user not found\n")
				return
			}

			var ppgg[]PermGroup
			db.Find(&ppgg, "pg_id=?", uu[0].PGID)

			if len(ppgg)==0{
				c.Status(http.StatusBadRequest)
				fmt.Printf("profile : PermGroup not found\n")
				return
			}

			sidebarC,dashboardC,functionC:=assignPerms(db,uu)

			//found.
			if verbose=="false"{
				c.JSON(http.StatusOK,Result{
					Name: uu[0].Name,
					Email: uu[0].Email,
					PGID: uu[0].PGID,
					PGName: ppgg[0].Name,
					SidebarC: sidebarC,
					DashboardC: dashboardC,
					FunctionC: functionC,
				})
			}else{
				c.JSON(http.StatusOK,ResultVerbose{
					Name: uu[0].Name,
				})
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
		c.SetCookie("user_token", "-1", 1000, "/", domain, false, true)
		c.Status(http.StatusOK)
	})

	router.POST("/login", func(c *gin.Context) {

		passwd := c.DefaultPostForm("passwd","nil")
		phone := c.DefaultPostForm("phone","nil")

		if passwd=="nil"||phone=="nil" {
			c.String(http.StatusNotAcceptable, fmt.Sprintln("network"))
			return
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
			c.SetCookie("user_token", uu.ID, 1000, "/", domain, false, true)
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
			newUser := User{strUsrUid, name, phone, email,passwd,0,"","保密"}
			db.Create(newUser)
			//register success
			c.String(http.StatusOK,fmt.Sprintf("ok"))
			//c.SetCookie("user_token", strUsrUid, 1000, "/", "domain", false, true)
		}else {
			//phone num is used!
			c.String(http.StatusOK, fmt.Sprintf("phone"))
		}

	})
}

func assignPerms(db *gorm.DB,uu []User) (map[string]int,map[string]int,map[string]int){
	//user sidebar content
	var sidebarC map[string]int
	sidebarC = make(map[string]int)
	var dashboardC map[string]int
	dashboardC = make(map[string]int)
	var functionC map[string]int
	functionC = make(map[string]int)

	var ggpp[] GroupPermission
	var permSet[] Permission
	db.Find(&permSet)

	for i:=0;i<len(permSet);i++{
		ggpp = []GroupPermission{}
		db.Where("perm_id = ? AND pg_id = ?",permSet[i].PermID,uu[0].PGID).Find(&ggpp)
		if strings.Contains(permSet[i].Name,"dashboard."){
			dashboardC[permSet[i].Name]=len(ggpp)
		}
		if strings.Contains(permSet[i].Name,"function."){
			functionC[permSet[i].Name]=len(ggpp)
			assignFatherSidebar(permSet[i].Name,&sidebarC,len(ggpp))
		}
	}
	return sidebarC,dashboardC,functionC
}

func assignFatherSidebar(name string,sidebarC * map[string]int,queryRes int){
	(*sidebarC)["sidebar.overview"]=1
	if strings.Contains(name,"waterman.") || strings.Contains(name,"userman.") {
		if(*sidebarC)["sidebar.manage"]<=0{
			(*sidebarC)["sidebar.manage"]=queryRes
		}
	}
}