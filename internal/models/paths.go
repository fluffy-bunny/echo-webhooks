package models

import (
	"echo-starter/internal/wellknown"
)

type Paths struct {
	Home  string
	About string
}

func NewPaths() *Paths {
	return &Paths{
		Home:  wellknown.HomePath,
		About: wellknown.AboutPath,
	}
}
