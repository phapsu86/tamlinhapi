package models

import (
	"errors"
	"html"
	"strings"
	
	"github.com/jinzhu/gorm"
)

type Hung struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Title     string    `gorm:"size:255;not null;unique" json:"title"`
	Content   string    `gorm:"size:255;not null;" json:"content"`
	
}

func (p *Hung) Prepare() {
	p.ID = 0
	p.Title = html.EscapeString(strings.TrimSpace(p.Title))
	p.Content = html.EscapeString(strings.TrimSpace(p.Content))
	
}

func (p *Hung) Validate() error {

	if p.Title == "" {
		return errors.New("Required Title")
	}
	if p.Content == "" {
		return errors.New("Required Content")
	}
	
	return nil
}

func (t *Hung) FindAllTest(db *gorm.DB) (*[]Hung, error) {
	var err error
	posts := []Hung{}
	err = db.Debug().Model(&Hung{}).Limit(100).Find(&posts).Error
	if err != nil {
		return &[]Hung{}, err
	}
	
	return &posts, nil
}



