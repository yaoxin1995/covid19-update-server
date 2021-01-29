package model

import (
	"time"

	"github.com/pmoule/go2hal/hal"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

var db *gorm.DB

// General fields that are required to persist a model.
type PersistentModel struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"-"`
}

// Interface for models that can be represented using the Hypertext Application Language.
// https://tools.ietf.org/html/draft-kelly-json-hal-08
type HALCompatibleModel interface {
	ToHAL(path string) hal.Resource
}

// Connection setup for the database.
// dbType is the type of the used database, e.g. sqlite3
// dbSource is the source of the used database, e.g. the file path
func SetupDB(dbType, dbSource string) error {
	var err error
	db, err = gorm.Open(dbType, dbSource)
	if err != nil {
		return err
	}
	db.LogMode(false)
	db.DB().SetMaxOpenConns(0)

	// Add domain models to DB
	err = db.AutoMigrate(&Subscription{}, &Topic{}, &Covid19Region{}, &Event{}).Error
	if err != nil {
		return err
	}

	// Use UTC time
	gorm.NowFunc = func() time.Time {
		return time.Now().UTC()
	}

	return err
}
