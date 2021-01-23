package model

import (
	"github.com/jinzhu/gorm"
	"github.com/pmoule/go2hal/hal"
)

type Covid19Region struct {
	PersistentModel
	Incidence float64 `json:"incidence"`
}

type Incidence struct {
	Incidence float64 `json:"incidence"`
}

func (i Incidence) ToHAL(path string) hal.Resource {
	root := hal.NewResourceObject()
	root.AddData(i)

	selfRel := hal.NewSelfLinkRelation()
	selfLink := &hal.LinkObject{Href: path}
	selfRel.SetLink(selfLink)
	root.AddLink(selfRel)

	return root
}

func NewCovid19Region(cID uint, incidence float64) (Covid19Region, error) {
	i := Covid19Region{
		PersistentModel: PersistentModel{ID: cID},
		Incidence:       incidence,
	}
	err := i.Store()
	return i, err
}

func (c *Covid19Region) Store() error {
	return db.Save(&c).Error
}

func GetCovid19Region(cID uint) (*Covid19Region, error) {
	i := &Covid19Region{}
	err := db.Where("id = ?", cID).First(i).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return i, nil
}

func (c *Covid19Region) GetIncidence() Incidence {
	return Incidence{
		c.Incidence,
	}
}
