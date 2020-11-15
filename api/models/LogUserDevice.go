package models

import (
	"time"
	"errors"

	"github.com/jinzhu/gorm"
)

type LogUserDevice struct {
	ID        uint64 `gorm:"primary_key;auto_increment" json:"id"`
	UserID     uint64 `json:"user_id"`
	DeviceID      string `gorm:"size:255;not null;" json:"device_id"`
	DeviceToken   string `gorm:"size:255;not null;" json:"device_token"`
	LogType    int `json:"log_type"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *LogUserDevice) Prepare() {
	p.ID = 0
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	

}

func (p *LogUserDevice) Validate() error {

	// if p.Title == "" {
	// 	return errors.New("Required Title")
	// }
	// if p.Content == "" {
	// 	return errors.New("Required Content")
	// }

	return nil
}

func (t *LogUserDevice) FindLogUserDeviceID(db *gorm.DB, device_id string) (*LogUserDevice, error) {
	var err error
	item := LogUserDevice{}
	err = db.Debug().Model(&LogUserDevice{}).Where("device_id = ?", device_id).Take(&item).Error
	if err != nil {
		return &LogUserDevice{}, err
	}
	return &item, nil
}

func (p *LogUserDevice) CreateLogUserDevice(db *gorm.DB) error {
	var err error
	err = db.Debug().Model(&LogUserDevice{}).Create(&p).Error
	if err != nil {
		return err
	}
	return nil
}


func (p *LogUserDevice) DeleteLogDevice(db *gorm.DB, device_id string,uid uint64) (int64, error) {

	db = db.Debug().Model(&LogUserDevice{}).Where("device_id = ? and user_id = ? ", device_id,uid).Take(&LogUserDevice{}).Delete(&LogUserDevice{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Device not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}