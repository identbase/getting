package representor

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/identbase/getting/pkg/link"
)

/*
Representor interface provides a way to handle many different types of
acceptable representations.*/
type Representor interface {
	parse(b []byte) (interface{}, error)
	GetBody() interface{}
	GetLink(rt string) (*link.Link, error)
	GetLinks(rt string) []link.Link
	HasLink(rt string) bool
}

/*
Create creates a Representor based on the Content-Type. */
func Create(u url.URL, t string, b []byte) (Representor, error) {
	switch true {
	case strings.Contains(t, "application/hal+json"):
		return NewHALRepresentor(u, t, b)
	default:
		return nil, errors.New("unsupported content-type")
	}
}

/*
CreateFromResponse pulls the necessary information from the request and passes
it on to the Create function. */
func CreateFromResponse(u url.URL, r http.Response, b []byte) (Representor, error) {

	ct := r.Header.Get("Content-Type")
	if ct == "" {
		return nil, errors.New("missing content-type")
	}

	// TODO: Check the header for Link

	return Create(u, ct, b)
}
