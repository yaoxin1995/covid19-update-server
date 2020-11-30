package model

import (
	"github.com/jinzhu/gorm"
)

type Incidence struct {
	CommonModelFields
	Cases7Per100k float64 `json:"cases7_per_100k"`
}

func NewIncidence(iID uint, cases7per100k float64) (Incidence, error) {
	i := Incidence{
		CommonModelFields: CommonModelFields{ID: iID},
		Cases7Per100k:     cases7per100k,
	}
	err := i.Store()
	return i, err
}

func (i *Incidence) Store() error {
	return db.Save(&i).Error
}

func GetIncidence(iID uint) (*Incidence, error) {
	i := &Incidence{}
	err := db.Where("id = ?", iID).First(i).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return i, nil
}
