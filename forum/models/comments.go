package models

import (
	"encoding/xml"
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
