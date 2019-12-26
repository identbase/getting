package representor

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/identbase/getting/pkg/link"
)

/*
Representor interface provides a way to handle many different types of
acceptable representations.*/
type Representor interface {
	parse(b string) interface{}
	parseLinks(b string) []*link.Link
	GetBody() interface{}
	SetBody(b interface{})
	GetLink(rt string) *link.Link
	GetLinks(rt string) []*link.Link
	HasLink(rt string) bool
}

/*
Create creates a Representor based on the Content-Type. */
func Create(u url.URL, t string, b []byte) (*Representor, error) {
	return nil, errors.New("not implemented")
}

/*
CreateFromResponse pulls the necessary information from the request and passes
it on to the Create function. */
func CreateFromResponse(u url.URL, r http.Response, b []byte) (*Representor, error) {
	return nil, errors.New("not implemented")
}
