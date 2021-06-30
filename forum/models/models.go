package models

import (
	"regexp"

	"gorm.io/gorm"
)

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

func (m *Models) CreatePost(db *gorm.DB, p *Post) *gorm.DB {
	return db.Select("UserID", "Title", "Body").Create(&p)
}

func (m *Models) UpdatePost(db *gorm.DB, p *Post) *gorm.DB {
	return db.Model(&p).Updates(Post{Title: p.Title, Body: p.Body})
}

func (m *Models) DeletePost(db *gorm.DB, p *Post) *gorm.DB {
	return db.Where("userId = ?", p.UserID).Delete(&p)
}

//////////////////////////////////////////////////////////////////

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

func (m *Models) CreateComment(db *gorm.DB, c *Comment) *gorm.DB {
	reEmail := regexp.MustCompile(`^[^@]+@[^@]+\.\w{1,5}$`)
	if c.Email != "" && !reEmail.Match([]byte(c.Email)) {
		return &gorm.DB{Error: gorm.ErrInvalidValue}
	}
	return db.Select("PostID", "UserID", "Name", "Email", "Body").Create(&c)
}
func (m *Models) UpdateComment(db *gorm.DB, c *Comment) *gorm.DB {
	reEmail := regexp.MustCompile(`^[^@]+@[^@]+\.\w{1,5}$`)
	if c.Email != "" && !reEmail.Match([]byte(c.Email)) {
		return &gorm.DB{Error: gorm.ErrInvalidValue}
	}
	return db.Model(&c).Updates(Comment{Name: c.Name, Email: c.Email, Body: c.Body})
}
func (m *Models) DeleteComment(db *gorm.DB, c *Comment) *gorm.DB {
	return db.Where("userId = ?", c.UserID).Delete(&c)
}
