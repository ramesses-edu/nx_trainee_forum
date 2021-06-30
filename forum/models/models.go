package models

import "gorm.io/gorm"

type Models struct {
}

func (m *Models) GetPost(db *gorm.DB, param map[string]interface{}) (Post, *gorm.DB) {
	p := Post{}
	tx := db.Where(param).First(&p)
	return p, tx
}

func (m *Models) ListPosts(db *gorm.DB, param map[string]interface{}) ([]Post, *gorm.DB) {
	pp := []Post{}
	tx := db.Where(param).Find(&pp)
	return pp, tx
}

func (m *Models) GetComment(db *gorm.DB, param map[string]interface{}) (Comment, *gorm.DB) {
	c := Comment{}
	tx := db.Where(param).First(&c)
	return c, tx
}

func (m *Models) ListComments(db *gorm.DB, param map[string]interface{}) ([]Comment, *gorm.DB) {
	cc := []Comment{}
	tx := db.Where(param).Find(&cc)
	return cc, tx
}
