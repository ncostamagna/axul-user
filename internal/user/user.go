package user

import (
	"time"

	"gorm.io/gorm"
)

//User model
type User struct {
	ID           string         `gorm:"size:40;primaryKey" json:"id"`
	UserName     string         `gorm:"size:70;unique" json:"username"`
	FirstName    string         `json:"firstname"`
	LastName     string         `json:"lastname"`
	Password     string         `json:"password"`
	Email        string         `json:"email"`
	Phone        string         `json:"phone"`
	ClientID     string         `json:"client_id"`
	ClientSecret string         `json:"client_secret"`
	Token        string         `json:"token"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
