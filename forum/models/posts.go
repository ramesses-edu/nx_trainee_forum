package models

import (
	"encoding/xml"

	"gorm.io/gorm"
)

type Posts struct { //structure for response array of posts in xml format
	XMLName xml.Name `xml:"posts" json:"-" gorm:"-"`
	Posts   []Post   `xml:"post"`
}

/////////////////////////////////////////////////////////////////////////////////////////
type Post struct {
	UserID   int       `json:"userId" gorm:"column:userId"`
	ID       int       `json:"id" gorm:"column:id;primaryKey"`
	Title    string    `json:"title" gorm:"column:title;type:VARCHAR(256)"`
	Body     string    `json:"body" gorm:"column:body;type:VARCHAR(256)"`
	Comments []Comment `xml:"-" json:"-" gorm:"foreignKey:PostID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (p *Post) GetPost(db *gorm.DB, param map[string]interface{}) *gorm.DB {
	return db.Where(param).First(&p)
}

func (p *Post) CreatePost(db *gorm.DB) *gorm.DB {
	return db.Select("UserID", "Title", "Body").Create(&p)
}
func (p *Post) UpdatePost(db *gorm.DB) *gorm.DB {
	return db.Model(&p).Updates(Post{Title: p.Title, Body: p.Body})
}
func (p *Post) DeletePost(db *gorm.DB) *gorm.DB {
	return db.Where("userId = ?", p.UserID).Delete(&p)
}

func ListPosts(db *gorm.DB, param map[string]interface{}) ([]Post, *gorm.DB) {
	pp := []Post{}
	tx := db.Where(param).Find(&pp)
	return pp, tx
}
