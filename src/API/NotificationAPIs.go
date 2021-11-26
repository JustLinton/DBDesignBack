package API

import (
	"db_design/src/entity"
	"db_design/src/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"time"
)

func InitNotificationApi(err *error,db *gorm.DB,router *gin.Engine) {

	router.POST("/notification", func(c *gin.Context) {

		//if the user doesn't pass the validation check
		if !utils.UserValidation(209,db,err,router,c){
			return
		}

		var uu[] entity.User
		db.Find(&uu, "phone=?", c.DefaultPostForm("Phone",""))

		if len(uu)==0{
			c.Status(http.StatusBadRequest)
			fmt.Printf("/deluser : user not found\n")
			return
		}

		entity.CreateAndSendNotification(
			db,err,
			uu[0].ID,
			c.DefaultPostForm("Title",""),
			c.DefaultPostForm("Content",""),
			c.DefaultPostForm("Type",""),
			c.DefaultPostForm("SenderID",""),
			c.DefaultPostForm("IsSystem",""),
			c.DefaultPostForm("SystemTitle",""),
		)

		c.String(http.StatusOK,"ok")
	})


	router.GET("/notilist", func(c *gin.Context) {

		//no permissions needed.
		//get user token
		usrToken,cerr := c.Cookie("user_token")
		if cerr!= nil || usrToken=="-1"{
			fmt.Printf("cookie not found or not logged\n")
			c.String(http.StatusOK,fmt.Sprintf("notlogged"))
			return
		}

		type Row struct {
			Title string
			Content string
			Type string
			SenderName string
			SenderAva string
			HaveRead string
			RowNum int
			SendTimePassed string
			SendTime string
		}

		var tmpRow Row
		var rows []Row

		type Result struct {
			Rows []Row
			TotalNum int
		}

		var uan[]entity.UserAndNotification
		db.Find(&uan)
		//fmt.Printf("len: %d",len(uu) )

		rowNum:=0

		for i:=0;i<len(uan);i++{
			theUan:=uan[i]
			if theUan.ReceiverID!=usrToken{
				continue
			}

			if theUan.HaveRead=="true"{
				//dont list notifications that is read
				continue
			}

			var nn[]entity.Notification
			db.Find(&nn,"nt_id=?", theUan.NtID)

			if len(nn)==0{
				//notification not found
				fmt.Printf("notification not found.\n")
				continue
			}

			if nn[0].IsSystem=="true" {
				tmpRow.SenderName=nn[0].SystemTitle
				tmpRow.SenderAva="2"
			}else{
				var sender[]entity.User
				db.Find(&sender,"id=?", nn[0].SenderID)

				if len(sender)==0{
					//notification not found
					fmt.Printf("notification sender not found.\n")
					continue
				}
				tmpRow.SenderName=sender[0].Name
				tmpRow.SenderAva="linton"
			}

			sendTimePassed := "0 分钟前"
			secsPassed64 := time.Now().Unix()-nn[0].SendTime.Unix()
			secsPassed64Str := strconv.FormatInt(secsPassed64, 10)
			secsPassed16 ,_ := strconv.Atoi(secsPassed64Str)
			dayp,hrp,minp := utils.ResolveTime(secsPassed16)
			sendTimePassed = fmt.Sprintf("%d",minp) + " 分钟前"
			if hrp>0{
				sendTimePassed = fmt.Sprintf("%d",hrp) + " 小时前"
			}
			if dayp>0{
				sendTimePassed = fmt.Sprintf("%d",dayp) + " 天前"
			}

			tmpRow.Title=nn[0].Title
			tmpRow.Content=nn[0].Content
			tmpRow.Type=nn[0].Type
			tmpRow.SendTime=nn[0].SendTime.Format("01月02日 15:04")
			tmpRow.SendTimePassed=sendTimePassed
			tmpRow.HaveRead=theUan.HaveRead
			tmpRow.RowNum=rowNum
			rowNum+=1

			//only show the recent 8 messages.
			if rowNum<=8 {
				rows = append(rows, tmpRow)
			}
		}

		c.JSON(http.StatusOK,Result{
			Rows: rows,
			TotalNum: rowNum,
		})

	})
}
