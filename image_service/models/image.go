package models

type Image struct {
	ImageID     int      `json:"imageID"`
	Author      string   `json:"author"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Likes       int      `json:"likes"`
	Views       int      `json:"views"`
	ImageChunk  []byte   `json:"imageChunk"`
}
