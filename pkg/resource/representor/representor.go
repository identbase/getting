package representor

import (
	"github.com/identbase/getting/pkg/link"
)

/*
Representor interface provides a way to handle many different types of
acceptable representations.*/
type Representor interface {
	parse(b string) string
	parseLinks(b string) []*link.Link
	GetBody() string
	SetBody(b string)
	GetLink(rt string) *link.Link
	GetLinks(rt string) []*link.Link
	HasLink(rt string) bool
}
