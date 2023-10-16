package models

type Image struct {
	ImageID     int64      `json:"imageID" gorm:"primary_key;autoIncrement:true"`
	Author      string   `json:"author"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	ImageChunk  []byte   `json:"imageChunk"`
}