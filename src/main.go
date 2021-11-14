package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	uuid "github.com/satori/go.uuid"
	"log"
	"net/http"
)

// UserInfo 用户信息
type User struct {
	ID string
	Name string
	Phone string
	Email string
	PassSHA string
}

func main() {
	//连接和初始化数据库，进行连接错误处理
	db, err := connectDatabase()

	if err!= nil{
		panic(err)
	}
	defer db.Close()
	db.AutoMigrate(&User{})
	operateTest(db)

	//初始化API请求处理
	router := initAPIs(&err,db)

	//设置跨域
	router.Use(cors.Default())

	// 开始端口监听（进入阻塞）。默认端口是8080,也可以指定端口 r.Run(":80")
	router.Run(":8080")
}


func initAPIs(err *error,db *gorm.DB) *gin.Engine {
	//start up the gin frame and implement the methods to respond the http requests.
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK,"Hello！欢迎来到GO世界！")
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
		db.Find(&uu, "pass_sha=?", passwd)
		//fmt.Printf("pass_sha:%s", uu.Phone)
		if uu.PassSHA==passwd {
			//passwd is correct
			userUid := uuid.Must(uuid.NewV4(),*err)
			strUsrUid := fmt.Sprintf("%s",userUid)
			//login success
			c.Status(http.StatusOK)
			c.SetCookie("user_token", strUsrUid, 1000, "/", "localhost", false, true)

		}else {
			//not correct passwd!
			c.String(http.StatusNotAcceptable, fmt.Sprintln("passwd"))
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
			newUser := User{strUsrUid, name, phone, email,passwd}
			db.Create(newUser)
			//register success
			c.Status(http.StatusOK)
			c.SetCookie("user_token", strUsrUid, 1000, "/", "localhost", false, true)
		}else {
			//phone num is used!
			c.String(http.StatusNotAcceptable, fmt.Sprintln("phone"))
		}

	})

	//router.POST("form", func(c *gin.Context) {
	//
	//	pname := c.DefaultPostForm("pname","nil")
	//	pgender := c.PostForm("pgender")
	//	phobby := c.PostForm("phobby")
	//	fmt.Printf("pname:%s\n",pname)
	//	userUid := uuid.Must(uuid.NewV4(),*err)
	//	c.String(http.StatusOK, fmt.Sprintln("post 请求  string 格式话"))
	//
	//	fmt.Printf("uuid:%s\n", userUid)
	//	strUsrUid := fmt.Sprintf("%s",userUid)
	//	entry1 := UserInfo{strUsrUid,pname,pgender,phobby}
	//	db.Create(entry1)
	//})
	return router
}

func connectDatabase() (*gorm.DB, error) {
	//connect the database.
	db, err := gorm.Open("mysql", "backend:123456@(127.0.0.1:3306)/dbdesign?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		log.Printf("open database error:%v", err)
		return nil, err
	}
	return db, nil
}

// GetSHA256HashCode SHA256生成哈希值
func GetSHA256HashCode(message []byte)string{
	//方法一：
	//创建一个基于SHA256算法的hash.Hash接口的对象
	hash := sha256.New()
	//输入数据
	hash.Write(message)
	//计算哈希值
	bytes := hash.Sum(nil)
	//将字符串编码为16进制格式,返回字符串
	hashCode := hex.EncodeToString(bytes)
	//返回哈希值
	return hashCode
}

func operateTest(db *gorm.DB){

	//// 自动迁移(automatically creates the table.)
	//db.AutoMigrate(&UserInfo{})
	//
	//u1 := UserInfo{3, "枯藤", "男", "篮球"}
	//u2 := UserInfo{2, "topgoer.com", "女", "足球"}
	//// 创建记录
	//db.Create(&u1)
	//db.Create(&u2)
	//// 查询
	//var u = new(UserInfo)
	//db.First(u)
	//fmt.Printf("%#v\n", u)
	//var uu UserInfo
	//db.Find(&uu, "hobby=?", "足球")
	//fmt.Printf("%#v\n", uu)
	//// 更新
	//db.Model(&uu).Update("hobby", "双色球")
	//fmt.Printf("%#v\n", uu)
	//// 删除
	//db.Delete(&u)

}