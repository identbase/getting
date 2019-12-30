package resource

import (
	"io/ioutil"
	"net/http"
	// "net/http/httputil"
	"net/url"

	"github.com/identbase/getting/pkg/link"
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
	Representor        Representor
	nextRefreshHeaders map[string]string
}

/*
Getting interface represents the getting client object. */
type Getting interface {
	Go(u string) (*Resource, error)
}

/*
Representor interface represents the representation of the body of the request
or response. */
type Representor interface {
	GetBody() interface{}
	GetLink(rt string) (*link.Link, error)
	GetLinks(rt string) []link.Link
}

/*
New creates a new Resource object. */
func New(c Getting, u *url.URL) *Resource {
	r := Resource{
		Client: c,
		URI:    u,
	}

	return &r
}

/*
Link returns a specific link based on its rel. */
func (r *Resource) Link(rt string) (*link.Link, error) {
	repr, err := r.representation()
	if err != nil {
		return nil, err
	}
	return repr.GetLink(rt)
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

	return r.Representor.GetBody(), nil
}

/*
representation returns the resource in the specified representation. */
func (r *Resource) representation() (Representor, error) {
	if r.Representor == nil {
		// TODO: Use Resource.refresh() once we figure out how to not
		// cause a race condition
		// r.refresh()
		c := &http.Client{}

		// TODO: Figure out if we should be setting the body (3rd) parameter
		req, err := http.NewRequest("GET", r.URI.String(), nil)
		if err != nil {
			return nil, err
		}

		// req.Header.Set("Accept", r.ContentType)

		for k, v := range r.nextRefreshHeaders {
			req.Header.Set(k, v)
		}

		// dump, err := httputil.DumpRequestOut(req, true)
		// if err != nil {
		// 	return nil, err
		// }

		resp, err := c.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		// fmt.Printf("REQUEST %v %q", req, dump)

		// dump, err = httputil.DumpResponse(resp, true)
		// if err != nil {
		// 	return nil, err
		// }

		// fmt.Printf("RESPONSE %v %q", resp, dump)

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		r.Representor, err = representor.CreateFromResponse(*r.URI, *resp, body)
		if err != nil {
			return nil, err
		}

	}

	return r.Representor, nil
}

/*
Get fetches the resource representation. */
func (r *Resource) Get() (interface{}, error) {
	repr, err := r.representation()
	if err != nil {
		return nil, err
	}

	return repr.GetBody(), nil
}

/*
Go resolves a new resource based on a relative URI. */
func (r *Resource) Go(u string) (*Resource, error) {
	h, err := r.URI.Parse(u)
	if err != nil {
		return nil, err
	}

	return r.Client.Go(h.String())
}

/*
Follow follows a relationship, based on its reltype. For example, this might be
'alternate', 'item', 'edit', or a custom url-based one. */
func (r *Resource) Follow(rt string) (*Resource, error) {
	l, err := r.Link(rt)
	if err != nil {
		return nil, err
	}

	h, err := l.Resolve()
	if err != nil {
		return nil, err
	}

	return r.Go(h)
}
