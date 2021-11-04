package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/satori/go.uuid"
	"log"
	"net/http"
)

// UserInfo 用户信息
type UserInfo struct {
	ID string
	Name string
	Gender string
	Hobby string
}

func main() {
	//连接和初始化数据库，进行连接错误处理
	db, err := connectDatabase()

	if err!= nil{
		panic(err)
	}
	defer db.Close()
	db.AutoMigrate(&UserInfo{})
	operateTest(db)

	//初始化API请求处理
	router := initAPIs(&err,db)

	// 开始端口监听（进入阻塞）。默认端口是8080,也可以指定端口 r.Run(":80")
	router.Run(":8080")
}


func initAPIs(err *error,db *gorm.DB) *gin.Engine {
	//start up the gin frame and implement the methods to respond the http requests.
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK,"Hello！欢迎来到GO世界！")
	})
	router.POST("form", func(c *gin.Context) {

		pname := c.PostForm("pname")
		pgender := c.PostForm("pgender")
		phobby := c.PostForm("phobby")
		userUid := uuid.Must(uuid.NewV4(),*err)

		fmt.Printf("uuid:%s\n", userUid)
		strUsrUid := fmt.Sprintf("%s",userUid)
		entry1 := UserInfo{strUsrUid,pname,pgender,phobby}
		db.Create(entry1)
	})
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