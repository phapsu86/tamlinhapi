package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Otp struct {
	ID        uint64 `gorm:"primary_key;auto_increment" json:"id"`
	Phone     string `gorm:"size:255;not null;unique" json:"phone"`
	Code      string `gorm:"size:255;not null;" json:"code"`
	Status    int
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *Otp) Prepare() {
	p.ID = 0
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	// p.Title = html.EscapeString(strings.TrimSpace(p.Title))
	// p.Content = html.EscapeString(strings.TrimSpace(p.Content))

}

func (p *Otp) Validate() error {

	// if p.Title == "" {
	// 	return errors.New("Required Title")
	// }
	// if p.Content == "" {
	// 	return errors.New("Required Content")
	// }

	return nil
}

func (t *Otp) FindOtpPhone(db *gorm.DB, phone string) (*Otp, error) {
	var err error
	item := Otp{}
	err = db.Debug().Model(&Otp{}).Where("phone = ?", phone).Take(&item).Error
	if err != nil {
		return &Otp{}, err
	}
	return &item, nil
}

func (p *Otp) CreateOtp(db *gorm.DB) error {
	var err error
	err = db.Debug().Model(&Otp{}).Create(&p).Error
	if err != nil {
		return err
	}
	return nil
}
