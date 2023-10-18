package database

import (
	"errors"
	"padrecommendations/models"
)

// todo: implement postgresql database

type AnalyticsPostgresDBOld struct {
	tags   []*models.Tag
	images []*models.Image
}

func NewAnalyticsPostgresDB() *AnalyticsPostgresDBOld {
	return &AnalyticsPostgresDBOld{}
}

func (db *AnalyticsPostgresDBOld) CreateTaglist(taglist *models.Tag) {
	db.tags = append(db.tags, taglist)
}

func (db *AnalyticsPostgresDBOld) GetTags() []string {
	var tags []string
	for _, tag := range db.tags {
		tags = append(tags, tag.Tagname)
	}

	return tags
}

func (db *AnalyticsPostgresDBOld) GetTaglist(tagname string) (*models.Tag, error) {
	for _, tag := range db.tags {
		if tag.Tagname == tagname {
			return tag, nil
		}
	}

	return &models.Tag{}, errors.New("Tag not found")
}

func (db *AnalyticsPostgresDBOld) AddImage(image models.Image) {
	images := db.images
	for _, img := range images {
		if img.ImageID == image.ImageID {
			return
		}
	}

	for _, tag := range image.Tags {
		taglist, err := db.GetTaglist(tag)
		if err != nil {
			taglist = &models.Tag{Tagname: tag}
			db.CreateTaglist(taglist)
		}
		db.images = append(db.images, &image)

		taglist.TotalEngagement += image.Engagement
	}
}

func (db *AnalyticsPostgresDBOld) GetImage(imageID int) (*models.Image, error) {
	for _, image := range db.images {
		if image.ImageID == imageID {
			return image, nil
		}
	}

	return nil, errors.New("Image not found")
}

func (db *AnalyticsPostgresDBOld) UpdateImage(image models.Image) {
	oldimage, err := db.GetImage(image.ImageID)
	if err != nil {
		return
	}

	oldimage.Views = image.Views
	oldimage.Likes = image.Likes
	oldimage.Engagement = image.Engagement
}

func (db *AnalyticsPostgresDBOld) AddViews(imageID int, views int) {
	image, err := db.GetImage(imageID)
	if err != nil {
		return
	}

	image.Views += views
	image.Engagement += views

	for _, tag := range image.Tags {
		taglist, err := db.GetTaglist(tag)
		if err != nil {
			continue
		}

		taglist.TotalEngagement += views
	}
}

func (db *AnalyticsPostgresDBOld) AddLikes(imageID, likes, engagement int) {
	image, err := db.GetImage(imageID)
	if err != nil {
		return
	}

	image.Likes += likes
	image.Engagement += engagement

	for _, tag := range image.Tags {
		taglist, err := db.GetTaglist(tag)
		if err != nil {
			continue
		}

		taglist.TotalEngagement += engagement
	}
}

func (db *AnalyticsPostgresDBOld) DeleteImage(imageID int) {
	for i, image := range db.images {
		if image.ImageID == imageID {
			db.images = append(db.images[:i], db.images[i+1:]...)
			break
		}
	}
}
