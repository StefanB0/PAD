package models

type DiscoveryResponse struct {
	Services []string `json:"services"`
}

type SubscribeRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type SubscribeResponse struct {
	Secretkey string `json:"secret_key"`
}