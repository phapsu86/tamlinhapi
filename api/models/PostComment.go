package models

import (
	"errors"
	//	"html"
	//	"strings"

	"time"

	"github.com/jinzhu/gorm"
)

// type result struct {
// 	TotalJoin int
// 	TotalShare int
// 	TotalFollow int
//   }

//Tb PostLsfc
type PostComment struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	PostID    uint64    `gorm:"size:11;not null;unique" json:"post_id"`
	UserID    uint64    `gorm:"size:11;not null;" json:"user_id"`
	Status    int       `gorm:"size:1;not null;" json:"status"`
	User      User      `json:"user"`
	Contents  string    `gorm:"size:255;not null;" json:"contents"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Person is a representation of a person
func (p *PostComment) Prepare() {
	p.ID = 0
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *PostComment) Validate() error {

	if p.PostID == 0 {
		return errors.New("Required POST_ID")
	}

	if p.Contents == "" {
		return errors.New("Required Comment")
	}

	return nil
}

//Save
func (p *PostComment) SavePostComment(db *gorm.DB) (*PostComment, error) {
	var err error
	item := Post{}
	if p.PostID != 0 {
		err = db.Debug().Model(&Post{}).Where("id = ?", p.PostID).Take(&item).Error
		if err != nil {
			return &PostComment{}, err
		}
		// truong hop co roi thi update
		err = db.Debug().Model(&PostLsfc{}).Create(&p).Error
		if err != nil {
			return &PostComment{}, err

		}
	}
	return p, nil
}

//FindALL
func (p *PostComment) FindAllPostComment(db *gorm.DB, pid uint64, page uint64) ([]PostComment, error) {
	var err error
	rows := []PostComment{}
	err = db.Debug().Model(&PostComment{}).Limit(10).Offset(page * 10).Find(&rows).Error
	if err != nil {
		return []PostComment{}, err
	}
	if len(rows) > 0 {
		for i, _ := range rows {
			err := db.Debug().Model(&User{}).Where("id = ?", rows[i].UserID).Take(&rows[i].User).Error
			if err != nil {
				return []PostComment{}, err
			}
		}
	}

	return rows, nil
}

// func (p *PostLsfc) FindReligionPostLsfjByID(db *gorm.DB, pid uint64,event_id uint64) (*PostLsfc, error) {
// 	var err error
// 	err = db.Debug().Model(&PostLsfc{}).Where("user_id = ? and post_id = ?", pid,event_id).Take(&p).Error
// 	if err != nil {
// 		return &PostLsfc{}, err
// 	}

// 	return p, nil
// }
