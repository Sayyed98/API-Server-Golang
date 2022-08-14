package models

import (
	"errors"

	"github.com/jinzhu/gorm"
)

type Tag struct {
	ID      uint64 `gorm:"primary_key;auto_incrment" json:"id"`
	TagList string `gorm:"size:255" json:"tags"`
	List    User   `josn:"listuser`
	Post    Post   `json:"post`
}

func (t *Tag) CreateTag(db *gorm.DB) (*Tag, error) {
	var err error
	err = db.Debug().Model(&User{}).Create(&t).Error
	if err != nil {
		return &Tag{}, err
	}
	if t.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", t.ID).Take(&t.List).Error
		if err != nil {
			return &Tag{}, err
		}
	}
	return t, nil

}
func (t *Tag) AllTags(db *gorm.DB) (*[]Tag, error) {
	var err error
	tags := []Tag{}
	err = db.Debug().Model(&Tag{}).Limit(100).Find(&tags).Error
	if err != nil {
		return &[]Tag{}, err
	}
	if len(tags) > 0 {
		for i, _ := range tags {
			err := db.Debug().Model(&User{}).Where("id = ?", tags[i].List).Take(&tags[i].TagList).Error
			if err != nil {
				return &[]Tag{}, err
			}
		}
	}
	return &tags, nil
}

func (t *Tag) UpdateTags(db *gorm.DB) (*Tag, error) {

	var err error

	err = db.Debug().Model(&Tag{}).Where("id = ?", t.ID).Updates(Post{Title: t.TagList}).Error
	if err != nil {
		return &Tag{}, err
	}
	if t.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", t.ID).Take(&t.TagList).Error
		if err != nil {
			return &Tag{}, err
		}
	}
	return t, nil
}

func (t *Tag) DeleteAPost(db *gorm.DB, tid uint64) (int64, error) {

	db = db.Debug().Model(&Tag{}).Where("id = ? and author_id = ?", tid).Take(&Tag{}).Delete(&Tag{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Post not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
