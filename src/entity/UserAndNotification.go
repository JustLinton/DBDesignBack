package entity

import "github.com/jinzhu/gorm"

type UserAndNotification struct {
	ReceiverID string
	NtID string
	HaveRead string
}

func InitUserAndNotificationEntity(db *gorm.DB){
	db.AutoMigrate(&UserAndNotification{})
}