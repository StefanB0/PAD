package models

type Taglist struct {
	ID              int     `json:"id" gorm:"primary_key;autoIncrement:true"`
	Tagname         string  `json:"tagname"`
	ImageList       []Image `json:"imageList" gorm:"many2many:tag_image;"`
	TotalEngagement int     `json:"totalEngagement"`
}

type Image struct {
	ImageID    int64    `json:"imageID" gorm:"primary_key;autoIncrement:true"`
	Tags       []string `json:"tags"`
	Views      int      `json:"views"`
	Likes      int      `json:"likes"`
	Engagement int      `json:"engagement"`
}
