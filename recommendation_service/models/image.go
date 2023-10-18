package models

import "gorm.io/gorm"

type ImageSQL struct {
	gorm.Model
	ImageID    int      `gorm:"primary_key"`
	Tags       []TagSQL `gorm:"many2many:image_tag;"`
	Views      int
	Likes      int
	Engagement int
	// Tags       []string `json:"tags"`
}

type TagSQL struct {
	gorm.Model
	Tagname         string `gorm:"primary_key"`
	TotalEngagement int
	// ImageList       []Image `json:"imageList" gorm:"many2many:tag_image;"`
	// ID      int    `json:"id" gorm:"primary_key;autoIncrement:true"`
}

type Image struct {
	ImageID    int      `json:"imageID"`
	Tags       []string `json:"tags"`
	Views      int      `json:"views"`
	Likes      int      `json:"likes"`
	Engagement int      `json:"engagement"`
}

type Tag struct {
	Tagname         string `json:"tagname"`
	TotalEngagement int    `json:"totalEngagement"`
}
