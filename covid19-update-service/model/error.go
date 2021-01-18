package model

import (
	"github.com/pmoule/go2hal/hal"
)

type ErrorT struct {
	Error string `json:"error"`
}

func NewError(message string) ErrorT {
	return ErrorT{Error: message}
}

func (e ErrorT) ToHAL(path string) hal.Resource {

	root := hal.NewResourceObject()
	root.AddData(e)

	selfRel := hal.NewSelfLinkRelation()
	selfLink := &hal.LinkObject{Href: path}
	selfRel.SetLink(selfLink)
	root.AddLink(selfRel)

	return root
}
