package controller

type getImageRequest struct {
	ImageID int `json:"imageID"`
}

type getImageResponse struct {
	ImageID     int      `json:"imageID"`
	Author      string   `json:"author"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

type uploadRequest struct {
	Token       string   `json:"token"`
	SagaID      string   `json:"sagaid"`
	Author      string   `json:"author"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	ImageBytes  []byte   `json:"image"`
}

type likeRequest struct {
	ImageID int `json:"imageID"`
}

type updateRequest struct {
	Token       string `json:"token"`
	ImageID     int    `json:"imageID"`
	Author      string `json:"author"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type deleteRequest struct {
	Token   string `json:"token"`
	ImageID int64  `json:"imageID"`
}
