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

// TODO: Below really doesnt need to exist, its just for easy copy pasting

/*
BaseRepresentor is a basic template implementation of a Representor. It is
basically a 'body' of a request or response.

This should not be used, refer to hal.go for the HalRepresentor as a good
default option to use. */
type BaseRepresentor struct {
	uri         string
	contentType string
	body        string
}

/*
New creates a new Representor object */
func New(u string, ct string, b string) *BaseRepresentor {
	r := BaseRepresentor{
		uri:         u,
		contentType: ct,
		body:        b,
	}

	return &r
}

/*
parse */
func (r *BaseRepresentor) parse(b string) string {
	return ""
}

/*
parseLinks */
func (r *BaseRepresentor) parseLinks(b string) []*link.Link {
	return []*link.Link{}
}

/*
GetBody */
func (r *BaseRepresentor) GetBody() string {
	return ""
}

/*
SetBody */
func (r *BaseRepresentor) setBody(b string) {}

/*
GetLink */
func (r *BaseRepresentor) GetLink(rt string) *link.Link {
	return nil
}

/*
GetLinks */
func (r *BaseRepresentor) GetLinks(rt string) []*link.Link {
	return []*link.Link{}
}

/*
HasLink */
func (r *BaseRepresentor) HasLink(rt string) bool {
	return false
}
