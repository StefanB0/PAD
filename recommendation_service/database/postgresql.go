package database

import (
	"errors"
	"padrecommendations/models"
)

// todo: implement postgresql database

type AnalyticsPostgresDB struct {
	tags   []*models.Taglist
	images []*models.Image
}

func NewAnalyticsPostgresDB() *AnalyticsPostgresDB {
	return &AnalyticsPostgresDB{}
}

func (db *AnalyticsPostgresDB) CreateTaglist(taglist *models.Taglist) {
	db.tags = append(db.tags, taglist)
}

func (db *AnalyticsPostgresDB) GetTags() []string {
	var tags []string
	for _, tag := range db.tags {
		tags = append(tags, tag.Tagname)
	}

	return tags
}

func (db *AnalyticsPostgresDB) GetTaglist(tagname string) (*models.Taglist, error) {
	for _, tag := range db.tags {
		if tag.Tagname == tagname {
			return tag, nil
		}
	}

	return &models.Taglist{}, errors.New("Tag not found")
}

func (db *AnalyticsPostgresDB) AddImage(image models.Image) {
	for _, tag := range image.Tags {
		taglist, err := db.GetTaglist(tag)
		if err != nil {
			taglist = &models.Taglist{Tagname: tag}
			db.CreateTaglist(taglist)
		}
		taglist.ImageList = append(taglist.ImageList, &image)
		db.images = append(db.images, &image)
	}
}

func (db *AnalyticsPostgresDB) GetImage(imageID int) (models.Image, error) {
	for _, image := range db.images {
		if image.ImageID == imageID {
			return *image, nil
		}
	}

	return models.Image{}, errors.New("Image not found")
}

func (db *AnalyticsPostgresDB) UpdateImage(image models.Image) {
	oldimage, err := db.GetImage(image.ImageID)
	if err != nil {
		return
	}

	oldimage.Views = image.Views
	oldimage.Likes = image.Likes
	oldimage.Engagement = image.Engagement
}

func (db *AnalyticsPostgresDB) AddViews(imageID int, views int) {
	image, err := db.GetImage(imageID)
	if err != nil {
		return
	}

	image.Views += views
	image.Engagement += views
}

func (db *AnalyticsPostgresDB) AddLikes(imageID int, likes int) {
	image, err := db.GetImage(imageID)
	if err != nil {
		return
	}

	image.Likes += likes
	image.Engagement += likes * 50
}

func (db *AnalyticsPostgresDB) DeleteImage(imageID int) {
	for i, image := range db.images {
		if image.ImageID == imageID {
			db.images = append(db.images[:i], db.images[i+1:]...)
			break
		}
	}

	for _, tag := range db.tags {
		for i, image := range tag.ImageList {
			if image.ImageID == imageID {
				tag.ImageList = append(tag.ImageList[:i], tag.ImageList[i+1:]...)
				break
			}
		}
	}
}
