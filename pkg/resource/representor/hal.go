package representor

import (
	// "bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/identbase/getting/pkg/link"
)

/*
HALRepresentor is a basic template implementation of a Representor. It is
basically a 'body' of a request or response.

This should not be used, refer to hal.go for the HalRepresentor as a good
default option to use. */
type HALRepresentor struct {
	URI         url.URL
	ContentType string
	Body        interface{}
	Links       link.LinkSet
}

type HALLink struct {
	HRef  string `json:"href"`
	Name  string `json:"name,omitempty"`
	Title string `json:"title"`
	// Templated string
	Type string `json:"type,omitempty"`
}

type HALBody struct {
	// TODO: This does not handle multiple "item" links. Need to handle
	// either a single HALLink or an array of HALLinks
	Links map[string][]HALLink `json:"_links,omitempty"`
	// This should be only JSON acceptable types: string, int, float, bool
	Properties map[string]interface{} `json:"-"`
	Embedded   map[string]HALBody     `json:"_embedded,omitempty"`
}

func mapInterfaceToHALLink(i map[string]interface{}) HALLink {
	// TODO: This seems like there is a
	// better way to do this
	var link HALLink
	_, nok := i["name"]
	_, tok := i["type"]

	if !nok && !tok {
		link = HALLink{
			HRef:  i["href"].(string),
			Title: i["title"].(string),
		}
	} else {
		if nok && tok {
			link = HALLink{
				HRef:  i["href"].(string),
				Name:  i["name"].(string),
				Title: i["title"].(string),
				Type:  i["type"].(string),
			}
		} else {
			if nok {
				link = HALLink{
					HRef:  i["href"].(string),
					Name:  i["name"].(string),
					Title: i["title"].(string),
				}
			}

			if tok {
				link = HALLink{
					HRef:  i["href"].(string),
					Title: i["title"].(string),
					Type:  i["type"].(string),
				}
			}
		}
	}

	return link
}

/*
UnmarshalJSON will properly convert JSON back into a HALBody object. */
func (b *HALBody) UnmarshalJSON(d []byte) error {
	var r map[string]interface{}

	if err := json.Unmarshal(d, &r); err != nil {
		return err
	}

	for bk, bv := range r {
		switch bk {
		case "_links":
			links := map[string][]HALLink{}

			lbuf := bv.(map[string]interface{})

			// TODO: Unmarshal the other potential HALLink data
			// _links: {
			for lk, lv := range lbuf {
				if lvbuf, ok := lv.(map[string]interface{}); ok {

					// self: {
					links[lk] = []HALLink{
						mapInterfaceToHALLink(lvbuf),
					}
				} else if lvbuf, ok := lv.([]interface{}); ok {
					// item: [
					for i := 0; i < len(lvbuf); i++ {
						// TODO: This seems like there is a
						// better way to do this
						lvbufitem := lvbuf[i].(map[string]interface{})

						links[lk] = append(links[lk], mapInterfaceToHALLink(lvbufitem))
					}
				}
			}

			b.Links = links
		case "_embedded":
			// TODO: Unmarshal HalBody recursively
			// b.Embedded = v.([]*HALBody)
			fmt.Println("TODO: Unmarshal _links", bk, bv)
		default:
			if b.Properties == nil {
				b.Properties = map[string]interface{}{}
			}

			b.Properties[bk] = bv
		}
	}

	return nil
}

/*
MarshalJSON will properly convert a HALBody into JSON. */
func (b HALBody) MarshalJSON() ([]byte, error) {
	var r map[string]interface{}

	for k, v := range b.Properties {
		r[k] = v
	}

	r["_embedded"] = map[string]HALBody{}
	for k, v := range b.Embedded {
		// TODO: Marshal the _embedded stuff recursively
		fmt.Println("TODO: Marshal _embedded", k, v)
	}

	r["_links"] = map[string]HALLink{}
	for k, v := range b.Links {
		// TODO: Marshal the _links
		fmt.Println("TODO: Marshal _links", k, v)
	}

	buf, err := json.Marshal(&r)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

/*
New creates a new Representor object. */
func NewHALRepresentor(u url.URL, ct string, b []byte) (*HALRepresentor, error) {
	r := HALRepresentor{
		URI:         u,
		ContentType: ct,
	}

	if b != nil {
		p, err := r.parse(b)
		if err != nil {
			return nil, err
		}

		r.setBody(p)
	}

	return &r, nil
}

/*
parse converts a JSON byte string into a HALBody. */
func (r *HALRepresentor) parse(b []byte) (interface{}, error) {
	var h HALBody

	if err := json.Unmarshal(b, &h); err != nil {
		return nil, err
	}

	return h, nil
}

// TODO: This function should be called in parse (and there should also be a
// parseEmbedded called there too)
/*
parseLinks converts HALLink into link.Link. */
func (r *HALRepresentor) parseLinks(b interface{}) []*link.Link {
	h := b.(HALBody)
	l := []*link.Link{}

	// TODO: This doesnt account for links in the Header
	for k, v := range h.Links {
		for i := 0; i < len(v); i++ {
			l = append(l, &link.Link{
				Context: r.URI.String(),
				HRef:    v[i].HRef,
				Rel:     k,
				Name:    v[i].Name,
				Title:   v[i].Title,
				Type:    v[i].Type,
			})
		}
	}

	return l
}

/*
GetBody */
func (r *HALRepresentor) GetBody() interface{} {
	return r.Body
}

/*
setBody */
func (r *HALRepresentor) setBody(b interface{}) {
	r.Body = b.(HALBody)
	// TODO: Initialize this somewhere else
	r.Links = link.LinkSet{}
	// TODO: The links dont need to be reparsed here really
	l := r.parseLinks(r.Body)
	for _, v := range l {
		if r.Links.Has(v.Rel) {
			r.Links.Add(v.Rel, *v)
		} else {
			r.Links.Set(v.Rel, append([]link.Link{}, *v))
		}
	}
}

/*
GetLink */
func (r *HALRepresentor) GetLink(rt string) (*link.Link, error) {
	l := r.Links.Get(rt)

	if len(l) == 0 {
		return nil, errors.New("link not found")
	}

	return &l[0], nil
}

/*
GetLinks */
func (r *HALRepresentor) GetLinks(rt string) []link.Link {

	if rt == "" {
		return r.Links.Values()
	}

	return r.Links.Get(rt)
}

/*
HasLink */
func (r *HALRepresentor) HasLink(rt string) bool {
	return r.Links.Has(rt)
}
