package entity

import (
	"fmt"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Notification struct {
	NtID string `gorm:"primary_key"`
	Title string
	Content string
	Type string
	SenderID string
	IsSystem string
	SystemTitle string
	SendTime time.Time
}

func InitNotificationEntity(db *gorm.DB){
	db.AutoMigrate(&Notification{})
}

func CreateAndSendNotification(db *gorm.DB,err *error,receiverUid_ string,title_ string,content_ string,type_ string,senderID_ string,isSystem_ string,systemTitle_ string){
	nn := CreateNotification(db,err,title_,content_,type_,senderID_,isSystem_,systemTitle_)
	var nuan = new(UserAndNotification)
	nuan.NtID=nn.NtID
	nuan.HaveRead="false"
	nuan.ReceiverID=receiverUid_
	db.Create(nuan)
}

func CreateNotification(db *gorm.DB,err *error,title_ string,content_ string,type_ string,senderID_ string,isSystem_ string,systemTitle_ string) Notification{
	ntUid := uuid.Must(uuid.NewV4(),*err)
	ntStrUid := fmt.Sprintf("%s",ntUid)
	var nn = new(Notification)
	nn.Title=title_
	nn.Content=content_
	nn.Type=type_
	nn.NtID=ntStrUid
	nn.SenderID=senderID_
	nn.IsSystem=isSystem_
	nn.SystemTitle=systemTitle_
	nn.SendTime=time.Now()
	db.Create(nn)
	return *nn
}