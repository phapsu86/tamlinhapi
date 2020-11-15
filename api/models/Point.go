package models

import (
	"errors"
	//"html"
	//"strings"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

type Point struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	UserID    uint64    `gorm:"not null;unique" json:"intro_image"`
	Amount    int64    `gorm:"not null;unique" json:"intro"`
	RefID     string    `gorm:"not null;unique" json:"ref_id"`
	Type      string    `gorm:"size:255;not null;unique" json:"type"`
	Note      string    `gorm:"size:255;not null;" json:"note"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (p *Point) Prepare() {
	p.ID = 0
	p.CreatedAt = time.Now()

}

func (p *Point) Validate() error {

	if p.UserID == 0 {
		return errors.New("Required UserID")
	}
	if p.Amount == 0 {
		return errors.New("Required Số Lượng")
	}
	if p.Type == "" {
		return errors.New("Required Author")
	}
	return nil
}

func (p *Point) SavePoint(db *gorm.DB) (*Point, error) {

	user := User{}
	var err error
	err = db.Debug().Model(&Point{}).Create(&p).Error
	if err != nil {
		return &Point{}, err
	}
	if p.UserID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.UserID).Take(&user).Error
		if err != nil {
			return &Point{}, err
		}
	}
	return p, nil
}

func (p *Point) FindPointByID(db *gorm.DB, pid uint64) (int64, error) {
	var err error

	points, err := db.Debug().Model(&Point{}).Select("sum(amount) as total_point").Where("user_id = ?", pid).Group("user_id").Rows()
	if err != nil {
		return 0, err
	}
	var total_point int64

	for points.Next() {

		points.Scan(&total_point)
		fmt.Println(total_point)

	}

	return total_point, nil

}
