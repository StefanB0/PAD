package controller

type getTagsResponse struct {
	Tags []string `json:"tags"`
}

type recommendRequest struct {
	Tag string `json:"tag"`
}

type addImageRequest struct {
	ID   int      `json:"id"`
	Tags []string `json:"tags"`
}

type updateImageRequest struct {
	ID   int      `json:"id"`
	Views int     `json:"views"`
	Likes int     `json:"likes"`
}