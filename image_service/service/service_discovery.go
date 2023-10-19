package service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"padimage/models"
)

type DiscoveryService struct {
	address   string
	selfaddr  string
	secretkey string
}

func NewDiscoveryService(selfaddr, address string) *DiscoveryService {
	return &DiscoveryService{
		address:  address,
		selfaddr: selfaddr,
	}
}

func (s *DiscoveryService) Subscribe() error {
	payload := new(bytes.Buffer)

	json.NewEncoder(payload).Encode(models.SubscribeRequest{
		Address: s.selfaddr,
		Name:    "image_service",
	})

	res, err := http.Post(s.address+"/service", "application/json", payload)
	if err != nil {
		return err
	}

	var response models.SubscribeResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return err
	}

	s.secretkey = response.Secretkey

	return nil
}

func (s *DiscoveryService) Unsubscribe() error {
	return nil
}

func (s *DiscoveryService) GetServiceAddress(serviceName string) ([]string, error) {
	url := s.address + "/service/" + serviceName

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	var serviceList models.DiscoveryResponse
	err = json.NewDecoder(res.Body).Decode(&serviceList)
	if err != nil {
		return nil, err
	}

	return serviceList.Services, nil
}
