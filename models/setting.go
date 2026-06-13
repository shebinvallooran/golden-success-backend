package models

import "time"

type NotificationSetting struct {
	ID                       uint      `json:"id" gorm:"primaryKey"`
	NotificationEmail        string    `json:"notification_email" gorm:"type:varchar(255)" validate:"omitempty,email"`
	EnableEmailNotifications bool      `json:"enable_email_notifications" gorm:"default:false"`
	SenderEmail              string    `json:"sender_email" gorm:"type:varchar(255)"`
	SMTPHost                 string    `json:"smtp_host" gorm:"type:varchar(255)"`
	SMTPPort                 int       `json:"smtp_port" gorm:"default:587"`
	SMTPUsername             string    `json:"smtp_username" gorm:"type:varchar(255)"`
	SMTPPassword             string    `json:"smtp_password" gorm:"type:varchar(255)"`
	SMTPSecure               bool      `json:"smtp_secure" gorm:"default:false"`
	UpdatedAt                time.Time `json:"updated_at"`
}
