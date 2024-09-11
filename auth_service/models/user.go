package models

type User struct {
	ID       int    `json:"id" gorm:"primary_key;autoIncrement:true"`
	Username string `json:"username"`
	Password string `json:"password"`
}
