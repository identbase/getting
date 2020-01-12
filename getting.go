package getting

import (
	"errors"
	"net/url"

	"github.com/identbase/getting/pkg/resource"
)

/*
Getting is the core client object. This is the starting point for working with
Getting. */
type Getting struct {
	// Bookmark is the default uri to use on all requests.
	bookmark string
}

/*
New creates a new Getting object. */
func New(b string) (*Getting, error) {
	if b == "" {
		return nil, errors.New("bookmark unspecified")
	}

	g := &Getting{
		bookmark: b,
	}

	return g, nil
}

/*
Follow is a shortcut for Go. */
func (g *Getting) Follow(rt string, v map[string]string) (*resource.Resource, error) {
	r, err := g.Go("")
	if err != nil {
		return nil, err
	}

	return r.Follow(rt, v)

}

/*
Go returns a resource by its uri. This function doesnt require a uri
if one is not specified, it will return the bookmark resource. */
func (g *Getting) Go(u string) (*resource.Resource, error) {

	ubuf, err := url.Parse(g.bookmark)
	if err != nil {
		return nil, err
	}

	uri, err := ubuf.Parse(u)
	if err != nil {
		return nil, err
	}

	// TODO: Use some sort of cache system to prevent rerequesting things

	return resource.New(g, uri), nil
}
