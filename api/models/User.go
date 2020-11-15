package models

import (
	"errors"
	"html"
	"log"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id" `
	Nickname  string    `gorm:"size:255;not null;unique" json:"nickname"`
	Mobile    string    `gorm:"size:255;not null;unique" json:"mobile"`
	Address   string    `gorm:"size:255;not null" json:"address"`
	Birthday  string    `gorm:"size:255;not null" json:"birthday"`
	Avartar   string    `gorm:"size:255;not null" json:"avatar"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	Point     int64     `json:"point" `
	Password  string    `gorm:"size:100;not null;" json:"password"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Otp       string
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (u *User) BeforeSave() error {
	hashedPassword, err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	//u.Password = "aaaaatest"
	return nil
}

func (u *User) BeforeUpdate() error {

	return nil
}

func (u *User) Prepare() {
	u.ID = 0
	u.Nickname = html.EscapeString(strings.TrimSpace(u.Nickname))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

func (u *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if u.Nickname == "" {
			return errors.New("Required Nickname")
		}

		if u.Email != "" {
			if err := checkmail.ValidateFormat(u.Email); err != nil {
				return errors.New("Invalid Email")
			}
		}

		return nil
	case "login":
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Mobile == "" {
			return errors.New("Required Mobile")
		}
		// if err := checkmail.ValidateFormat(u.Email); err != nil {
		// 	return errors.New("Invalid Email")
		// }
		return nil
	case "create":
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Mobile == "" {
			return errors.New("Required Mobile")
		}
		if u.Email != "" {
			if err := checkmail.ValidateFormat(u.Email); err != nil {
				return errors.New("Invalid Email")
			}
		} else {
			return errors.New("Required Email")

		}

		return nil

	default:
		if u.Mobile == "" {
			return errors.New("Required Mobile")
		}

		if u.Nickname == "" {
			return errors.New("Required Nickname")
		}
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil
	}
}

func (u *User) SaveUser(db *gorm.DB) (*User, error) {

	var err error
	err = db.Debug().Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) FindAllUsers(db *gorm.DB) (*[]User, error) {
	var err error
	users := []User{}
	err = db.Debug().Model(&User{}).Limit(100).Find(&users).Error
	if err != nil {
		return &[]User{}, err
	}
	return &users, err
}

func (u *User) FindUserByID(db *gorm.DB, uid uint64) (*User, error) {
	var err error
	err = db.Debug().Model(User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	pointModel := Point{}
	point, err := pointModel.FindPointByID(db, u.ID)
	err = db.Debug().Model(User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}

	u.Point = point
	if gorm.IsRecordNotFoundError(err) {
		return &User{}, errors.New("User Not Found")
	}
	return u, err
}

func (u *User) UpdatePasswordUser(db *gorm.DB, uid uint32) (*User, error) {

	// To hash the password
	err := u.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}
	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"password":  u.Password,
			"update_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &User{}, db.Error
	}
	// This is the display the updated user
	err = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) UpdateAUser(db *gorm.DB, uid uint64) (*User, error) {

	// To hash the password
	// err := u.BeforeSave()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"nickname": u.Nickname,
			"address":  u.Address,
			//	"password":  u.Password,
			"email": u.Email,
			//"avartar":   u.Avartar,
			"birthday":  u.Birthday,
			"update_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &User{}, db.Error
	}
	// This is the display the updated user
	err := db.Debug().Model(&User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) UpdateAvartarUser(db *gorm.DB, uid uint64) (*User, error) {

	// To hash the password
	//err := u.BeforeSave()
	var err error
	if err != nil {
		log.Fatal(err)
	}
	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			//"address":   u.Address,
			//"password":  u.Password,
			///"email":     u.Email,
			"avartar": u.Avartar,
			//"birthday":  u.Birthday,
			"update_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &User{}, db.Error
	}
	// This is the display the updated user
	err = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) DeleteAUser(db *gorm.DB, uid uint32) (int64, error) {

	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).Delete(&User{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
