package resource

import (
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/identbase/getting/pkg/resource/representor"
)

/*
Resource represents endpoint on a server. The endpoint has a uri, you  might for
example be able to GET its presentation. A resource may also have a list of
links on them, pointing to other resources. */
type Resource struct {
	Client             Getting
	URI                *url.URL
	ContentType        string
	Representor        *representor.Representor
	nextRefreshHeaders map[string]string
}

/*
Getting interface represents the getting client object. */
type Getting interface {
}

/*
New creates a new Resource object. */
func New(c Getting, u *url.URL) *Resource {
	r := Resource{
		Client: &c,
		URI:    u,
	}

	return &r
}

/*
refresh fetches the resource representation. */
func (r *Resource) refresh() (interface{}, error) {
	c := &http.Client{}

	// TODO: Figure out if we should be setting the body (3rd) parameter
	req, err := http.NewRequest("GET", r.URI.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", r.ContentType)

	for k, v := range r.nextRefreshHeaders {
		req.Header.Set(k, v)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	r.Representor, err = representor.CreateFromResponse(*r.URI, *resp, body)
	if err != nil {
		return nil, err
	}

	return (*r.Representor).GetBody(), nil
}

/*
representation returns the resource in the specified representation. */
func (r *Resource) representation() *representor.Representor {
	// TODO: Check if we should refresh here
	if r.Representor == nil {
		r.refresh()
	}

	return r.Representor
}

/*
Get fetches the resource representation. */
func (r *Resource) Get() interface{} {
	repr := r.representation()

	return (*repr).GetBody()
}

/*
Follow follows a relationship, based on its reltype. For example, this might be
'alternate', 'item', 'edit', or a custom url-based one. */
func (r *Resource) Follow(rt string) *Resource {
	return nil
}
