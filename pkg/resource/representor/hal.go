package representor

import (
	// "bytes"
	"encoding/json"
	// "errors"
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
	Body        string
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
					// TODO: This seems like there is a
					// better way to do this
					var link HALLink
					_, ok1 := lvbuf["name"]
					_, ok2 := lvbuf["type"]

					if ok1 {
						link = HALLink{
							HRef:  lvbuf["href"].(string),
							Name:  lvbuf["name"].(string),
							Title: lvbuf["title"].(string),
						}
					} else if ok2 {
						link = HALLink{
							HRef:  lvbuf["href"].(string),
							Title: lvbuf["title"].(string),
							Type:  lvbuf["type"].(string),
						}
					} else if ok1 && ok2 {
						link = HALLink{
							HRef:  lvbuf["href"].(string),
							Name:  lvbuf["name"].(string),
							Title: lvbuf["title"].(string),
							Type:  lvbuf["type"].(string),
						}
					} else {
						link = HALLink{
							HRef:  lvbuf["href"].(string),
							Title: lvbuf["title"].(string),
						}
					}
					// self: {
					links[lk] = []HALLink{link}
				} else if lvbuf, ok := lv.([]interface{}); ok {
					// item: [
					for i := 0; i < len(lvbuf); i++ {
						links[lk] = append(links[lk], HALLink{
							HRef:  lvbuf[i].(map[string]interface{})["href"].(string),
							Title: lvbuf[i].(map[string]interface{})["title"].(string),
						})
					}
				}
			}

			b.Links = links
		case "_embedded":
			// TODO: Unmarshal HalBody recursively
			// b.Embedded = v.([]*HALBody)
			fmt.Println("TODO: Unmarshal _links", bk, bv)
		default:
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

/*
parseLinks converts HALLink into link.Link. */
func (r *HALRepresentor) parseLinks(b string) []*link.Link {
	return []*link.Link{}
}

/*
GetBody */
func (r *HALRepresentor) GetBody() interface{} {
	return nil
}

/*
SetBody */
func (r *HALRepresentor) setBody(b interface{}) {}

/*
GetLink */
func (r *HALRepresentor) GetLink(rt string) *link.Link {
	return nil
}

/*
GetLinks */
func (r *HALRepresentor) GetLinks(rt string) []*link.Link {
	return []*link.Link{}
}

/*
HasLink */
func (r *HALRepresentor) HasLink(rt string) bool {
	return false
}
