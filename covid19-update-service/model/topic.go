package model

type Topic struct {
	CommonModelFields
	Position GPSPosition `gorm:"embedded;embeddedPrefix:position_" ,json:"position"`
}

type GPSPosition struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
