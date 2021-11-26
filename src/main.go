package main

import (
	"db_design/src/API"
	"db_design/src/entity"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"net/http"
)


func main() {
	//连接和初始化数据库，进行连接错误处理
	db, err := connectDatabase()

	if err!= nil{
		panic(err)
	}
	defer db.Close()
	//db.AutoMigrate(&entity.User{})
	//operateTest(db)

	//初始化API请求处理
	router := initAPIs(&err,db)

	//设置跨域
	router.Use(cors.Default())

	// 开始端口监听（进入阻塞）。默认端口是8080,也可以指定端口 r.Run(":80")
	router.Run(":8081")

}


func initAPIs(err *error,db *gorm.DB) *gin.Engine {
	//start up the gin frame and implement the methods to respond the http requests.
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK,"Hello！欢迎来到GO世界！")
	})

	entity.InitUsersApi(err,db,router)
	API.InitUsersApi2(err,db,router)
	entity.InitPermissionsApi(err,db,router)
	API.InitNotificationApi(err,db,router)

	return router
}

func connectDatabase() (*gorm.DB, error) {
	//connect the database.
	db, err := gorm.Open("mysql", "backend:123456@(127.0.0.1:3306)/dbdesign?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		log.Printf("open database error:%v", err)
		return nil, err
	}
	initEntities(db)
	return db, nil
}

func initEntities(db *gorm.DB){
	entity.InitUsers(db)
	entity.InitPermGroup(db)
	entity.InitPermission(db)
	entity.InitGroupPermission(db)
	entity.InitNotificationEntity(db)
	entity.InitUserAndNotificationEntity(db)
}