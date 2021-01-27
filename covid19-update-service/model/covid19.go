package model

import (
	"github.com/jinzhu/gorm"
	"github.com/pmoule/go2hal/hal"
)

// Covid19Region represents the region and its current incidence value. A region is identified by its id which corresponds to the OBJECTID that is used by the RKI Corona Landkreise API.
type Covid19Region struct {
	PersistentModel
	Incidence float64 `json:"incidence"`
}

// Wrapper object to represent the incidence value of a Covid19Region.
type Incidence struct {
	Incidence float64 `json:"incidence"`
}

// Represents the Incidence with the JSON Hypertext Application Language.
// path is the relative URI of the Incidence.
func (i Incidence) ToHAL(path string) hal.Resource {
	root := hal.NewResourceObject()
	root.AddData(i)

	selfRel := hal.NewSelfLinkRelation()
	selfLink := &hal.LinkObject{Href: path}
	selfRel.SetLink(selfLink)
	root.AddLink(selfRel)

	return root
}

// Creates a new Covid19Region.
// cID is the OBJECTID used by the RKI.
// incidence is the current COVID19 7-day incidence value.
func NewCovid19Region(cID uint, incidence float64) (Covid19Region, error) {
	i := Covid19Region{
		PersistentModel: PersistentModel{ID: cID},
		Incidence:       incidence,
	}
	err := i.Store()
	return i, err
}

// Persists the Covid19Region
func (c *Covid19Region) Store() error {
	return db.Save(&c).Error
}

// Gets a Covid19Region by its id or nil, if it was not found.
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

// Transforms a Covid19Region to a Incidence.
func (c *Covid19Region) GetIncidence() Incidence {
	return Incidence{
		c.Incidence,
	}
}
