package controller

type getImageRequest struct {
	ImageID int64 `json:"imageID"`
}

type getImageResponse struct {
	ImageID     int64    `json:"imageID"`
	Author      string   `json:"author"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

type uploadRequest struct {
	Token       string   `json:"token"`
	Author      string   `json:"author"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	ImageBytes  []byte   `json:"image"`
}

type updateRequest struct {
	Token       string `json:"token"`
	ImageID     int64  `json:"imageID"`
	Author      string `json:"author"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type deleteRequest struct {
	Token   string `json:"token"`
	ImageID int64  `json:"imageID"`
}
