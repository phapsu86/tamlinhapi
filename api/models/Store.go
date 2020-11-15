package models

import (

	"time"

	"github.com/jinzhu/gorm"
)

type Store struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Name     uint64    `gorm:"not null;unique" json:"name"`
	Address    string    `gorm:"not null" json:"address"`
	Lat 	float32  		`gorm:"not null" json:"lat"`
	Long     float32    `gorm:"not null" json:"long"`
	//Note   string    `gorm:"size:255;not null;" json:"note"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	Status int  		`gorm:"not null;unique" json:"status"`
}

func (p *Store) Prepare() {
	p.ID = 0
	p.CreatedAt = time.Now()

}

func (p *Store) Validate() error {

	// if p.Store == 0 {
	// 	return errors.New("Required UserID")
	// }
	// if p.Amount == 0 {
	// 	return errors.New("Required Số Lượng")
	// }
	// if p.Type =="" {
	// 	return errors.New("Required Author")
	// }
	return nil
}

func (p *Store) SaveStore(db *gorm.DB) (*Store, error) {

	user:= User{}
	var err error
	err = db.Debug().Model(&Point{}).Create(&p).Error
	if err != nil {
		return &Store{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.ID).Take(&user).Error
		if err != nil {
			return &Store{}, err
		}
	}
	return p, nil
}



func (p *Store) FindStoreByID(db *gorm.DB, pid uint64) (*Store, error) {
	var err error
	err = db.Debug().Model(&Store{}).Where("id = ? and status = 1",pid ).Take(&p).Error
	if err != nil {
		return &Store{}, err
	}
	

	return p, nil
	
}

