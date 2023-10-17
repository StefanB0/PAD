package database

import "padrecommendations/models"

type AnalyticsPostgresDB struct{}

func NewAnalyticsPostgresDB() *AnalyticsPostgresDB {
	return &AnalyticsPostgresDB{}
}

func (db *AnalyticsPostgresDB) CreateTaglist(tagname string) {

}

func (db *AnalyticsPostgresDB) GetTaglist(tagname string) models.Taglist {
	return models.Taglist{}
}

func (db *AnalyticsPostgresDB) AddImage(image models.Image) {

}

func (db *AnalyticsPostgresDB) GetImage(imageID int64) models.Image {
	return models.Image{}
}

func (db *AnalyticsPostgresDB) UpdateImage(image models.Image) {

}

func (db *AnalyticsPostgresDB) DeleteImage(imageID int64) {

}