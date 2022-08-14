package models

import (
	"html"
	"strings"

	"github.com/jinzhu/gorm"
)

type Comment struct {
	ID       uint64 `gorm:"primary_key;auto_increment" json:"id"`
	Content  string `gorm:"size:255" json:"content"`
	Author   User   `json:"author"`
	AuthorID Post   `json:"author_id"`
}

func (c *Comment) Prepare() {
	c.ID = 0
	c.Content = html.EscapeString(strings.TrimSpace(c.Content))
	c.Author = User{}

}

func (c *Comment) SaveComment(db *gorm.DB) (*Comment, error) {
	var err error
	err = db.Debug().Model(&Comment{}).Create(&c).Error
	if err != nil {
		return &Comment{}, err
	}
	if c.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", c.AuthorID).Take(&c.Author).Error
		if err != nil {
			return &Comment{}, err
		}
	}
	return c, nil
}

func (c *Comment) CommentAll(db *gorm.DB) (*[]Comment, error) {
	var err error
	comments := []Comment{}
	err = db.Debug().Model(&Comment{}).Limit(100).Find(&comments).Error
	if err != nil {
		return &[]Comment{}, err
	}
	if len(comments) > 0 {
		for i, _ := range comments {
			err := db.Debug().Model(&User{}).Where("id = ?", comments[i].AuthorID).Take(&comments[i].Author).Error
			if err != nil {
				return &[]Comment{}, err
			}
		}
	}
	return &posts, nil
}
