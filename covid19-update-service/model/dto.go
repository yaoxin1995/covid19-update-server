package model

import (
	"encoding/json"
	"fmt"

	"github.com/pmoule/go2hal/hal"
)

type SubscriptionRequestDTO struct {
	Email          *string `json:"email"`
	TelegramChatID *string `json:"telegramChatId"`
}

type TopicRequestDTO struct {
	Position  GPSPosition `json:"position"`
	Threshold uint        `json:"threshold"`
}

func (t *TopicRequestDTO) UnmarshalJSON(data []byte) (err error) {
	required := struct {
		Position  *GPSPosition `json:"position"`
		Threshold *uint        `json:"threshold"`
	}{}
	all := struct {
		Position  GPSPosition `json:"position"`
		Threshold uint        `json:"threshold"`
	}{}
	err = json.Unmarshal(data, &required)
	if err != nil {
		return err
	} else if required.Threshold == nil {
		err = fmt.Errorf("threshold is missing")
	} else if required.Position == nil {
		err = fmt.Errorf("position is missing")
	} else {
		err = json.Unmarshal(data, &all)
		t.Position = all.Position
		t.Threshold = all.Threshold
	}
	return err
}

type IncidenceDTO struct {
	Incidence float64 `json:"incidence"`
}

func (i IncidenceDTO) ToHAL(path string) hal.Resource {
	root := hal.NewResourceObject()
	root.AddData(i)

	selfRel := hal.NewSelfLinkRelation()
	selfLink := &hal.LinkObject{Href: path}
	selfRel.SetLink(selfLink)
	root.AddLink(selfRel)

	return root
}
