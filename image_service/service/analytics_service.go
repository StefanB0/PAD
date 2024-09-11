package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type AnalyticsService struct {
	sd     *DiscoveryService
	ls     *LoadBalancer
	client *http.Client
}

func NewAnalyticsService(sd *DiscoveryService) *AnalyticsService {
	return &AnalyticsService{
		sd: sd,
		ls: NewLoadBalancerService(),
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (s *AnalyticsService) UpdateLoadBalancer() {
	urls, err := s.sd.GetServiceAddress("analytics_service")
	if err != nil {
		log.Error().Err(err).Msg("Error getting analytics service address")
		return
	}
	s.ls.SetItems(urls)
}

func (s *AnalyticsService) AddEngagement(imageID, views, likes int) error {
	s.UpdateLoadBalancer()
	url := s.ls.GetItem()

	if url == "" {
		log.Error().Msg("No available analytics servers")
		return errors.New("No available analytics servers")
	}

	body := fiber.Map{
		"id":    imageID,
		"views": views,
		"likes": likes,
	}

	buffer := new(bytes.Buffer)
	if err := json.NewEncoder(buffer).Encode(body); err != nil {
		log.Error().Err(err).Msg("Error encoding body")
		return err
	}

	s.client.Post(url+"/updateImage", "application/json", buffer)

	return nil
}

func (s *AnalyticsService) AddImage(imageID int, sagaID string, tags []string) error {
	s.UpdateLoadBalancer()
	url := s.ls.GetItem()

	if url == "" {
		log.Error().Msg("No available analytics servers")
		return errors.New("No available analytics servers")
	}

	body := fiber.Map{
		"id":   imageID,
		"tags": tags,
		"sagaid": sagaID,
	}

	buffer := new(bytes.Buffer)
	if err := json.NewEncoder(buffer).Encode(body); err != nil {
		log.Error().Err(err).Msg("Error encoding body")
		return err
	}

	s.client.Post(url+"/addImage", "application/json", buffer)

	return nil
}