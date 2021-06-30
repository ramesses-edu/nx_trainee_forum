package models

import (
	"encoding/xml"
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
