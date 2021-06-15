package models

import (
	"encoding/xml"
	"regexp"

	"gorm.io/gorm"
)

type Comments struct { //structure for response array of comments in xml format
	XMLName  xml.Name  `xml:"comments" json:"-" gorm:"-"`
	Comments []Comment `xml:"comment"`
}

////////////////////////////////////////////////////////////////////////////////////////////////
type Comment struct {
	PostID int    `json:"postId" gorm:"column:postId"`
	UserID int    `json:"userId" gorm:"column:userId"`
	ID     int    `json:"id" gorm:"column:id;primaryKey"`
	Name   string `json:"name" gorm:"column:name;type:VARCHAR(256)"`
	Email  string `json:"email" gorm:"column:email;type:VARCHAR(256)"`
	Body   string `json:"body" gorm:"column:body;type:VARCHAR(256)"`
}

func (c *Comment) GetComment(db *gorm.DB, param map[string]interface{}) *gorm.DB {
	return db.Where(param).First(&c)
}

func (c *Comment) CreateComment(db *gorm.DB) *gorm.DB {
	reEmail := regexp.MustCompile(`^[^@]+@[^@]+\.\w{1,5}$`)
	if c.Email != "" && !reEmail.Match([]byte(c.Email)) {
		return &gorm.DB{Error: gorm.ErrInvalidValue}
	}
	return db.Select("PostID", "UserID", "Name", "Email", "Body").Create(&c)
}
func (c *Comment) UpdateComment(db *gorm.DB) *gorm.DB {
	reEmail := regexp.MustCompile(`^[^@]+@[^@]+\.\w{1,5}$`)
	if c.Email != "" && !reEmail.Match([]byte(c.Email)) {
		return &gorm.DB{Error: gorm.ErrInvalidValue}
	}
	return db.Model(&c).Updates(Comment{Name: c.Name, Email: c.Email, Body: c.Body})
}
func (c *Comment) DeleteComment(db *gorm.DB) *gorm.DB {
	return db.Where("userId = ?", c.UserID).Delete(&c)
}

func ListComments(db *gorm.DB, param map[string]interface{}) ([]Comment, *gorm.DB) {
	cc := []Comment{}
	tx := db.Where(param).Find(&cc)
	return cc, tx
}
