package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

//Tb UserRemind
type UserRemind struct {
	ID          uint64        `gorm:"primary_key;auto_increment" json:"id"`
	UserID      uint64        `gorm:"size:255;not null;unique" json:"user_id"`
	ObjectID    uint64        `gorm:"size:255;not null;" json:"object_id"`
	ObjectType  int8          `gorm:"size:255;not null;" json:"object_type"`
	BeginTime   time.Time     `gorm:"size:255;not null;" json:"begin_time"`
	EndTime     time.Time     `gorm:"size:255;not null;" json:"end_time"`
	EventDetail ReligionEvent `json:"event_detail"`
	CreatedAt   time.Time     `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time     `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Status      int           `gorm:"size:255;not null;" json:"status"`
}

// Person is a representation of a person
func (p *UserRemind) Prepare() {
	p.ID = 0
	//p.Name = html.EscapeString(strings.TrimSpace(p.Notes))
	//p.Image = html.EscapeString(strings.TrimSpace(p.Image))
	//p.Notes = p.Notes
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *UserRemind) Validate() error {

	// if p.Name == "" {
	// 	return errors.New("Required Title")
	// }
	// if p.Image == "" {
	// 	return errors.New("Required Content")
	// }
	// if p.Notes == "" {
	// 	return errors.New("Required Author")
	// }
	return nil
}

//Save
func (p *UserRemind) SaveUserRemind(db *gorm.DB) (*UserRemind, error) {
	var err error
	err = db.Debug().Model(&UserRemind{}).Create(&p).Error
	if err != nil {
		return &UserRemind{}, err
	}

	return p, nil
}

//FindALL
func (p *UserRemind) FindAllUserReminds(db *gorm.DB, uid uint64) ([]UserRemind, error) {
	var err error
	rows := []UserRemind{}
	err = db.Debug().Model(&UserRemind{}).Where("user_id = ?", uid).Limit(100).Find(&rows).Error
	if err != nil {
		return []UserRemind{}, err
	}
	if len(rows) > 0 {
		for i, _ := range rows {
			err := db.Debug().Model(&ReligionEvent{}).Where("id = ?", rows[i].ObjectID).Take(&rows[i].EventDetail).Error
			if err != nil {
				return []UserRemind{}, err
			}
		}
	}

	return rows, nil
}

func (p *UserRemind) FindUserRemindByEventID(db *gorm.DB, uid uint64, objectID uint64) (*UserRemind, error) {
	var err error
	rows := UserRemind{}
	err = db.Debug().Model(&UserRemind{}).Where("object_id = ? and user_id = ?", objectID, uid).Take(&rows).Error
	if err != nil {
		return &UserRemind{}, err
	}

	return &rows, nil

}

func (p *UserRemind) DeleteAUserRemind(db *gorm.DB, id uint64, uid uint64) (int64, error) {

	db = db.Debug().Model(&UserRemind{}).Where("object_id = ? and user_id = ?", id, uid).Take(&UserRemind{}).Delete(&UserRemind{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
