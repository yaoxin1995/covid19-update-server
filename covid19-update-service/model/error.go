package model

import (
	"github.com/pmoule/go2hal/hal"
)

// Representation of an error with an error message.
type ErrorT struct {
	Error string `json:"error"`
}

// Creates a new ErrorT with the given message.
func NewError(message string) ErrorT {
	return ErrorT{Error: message}
}

// Represents the ErrorT with the JSON Hypertext Application Language.
// path is the relative URI of the ErrorT.
func (e ErrorT) ToHAL(path string) hal.Resource {

	root := hal.NewResourceObject()
	root.AddData(e)

	selfRel := hal.NewSelfLinkRelation()
	selfLink := &hal.LinkObject{Href: path}
	selfRel.SetLink(selfLink)
	root.AddLink(selfRel)

	return root
}
