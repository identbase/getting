package representor

import (
	"bytes"
	"testing"
)

func Test_HALBody_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		d    []byte
		want *HALBody
		err  error
	}{
		{
			"success simple",
			bytes.NewBufferString(`{"_links": {"self": {"href": "/","title": "Test"}},"foo": "bar"}`).Bytes(),
			&HALBody{
				Links: map[string][]HALLink{
					"self": []HALLink{
						HALLink{
							HRef:  "/",
							Title: "Test",
						},
					},
				},
			},
			nil,
		},
		{
			"success full",
			bytes.NewBufferString(`{"_links": {"self": {"href": "/","title": "Test", "name": "alternate", "type": "blob"}},"foo": "bar"}`).Bytes(),
			&HALBody{
				Links: map[string][]HALLink{
					"self": []HALLink{
						HALLink{
							HRef:  "/",
							Name:  "alternate",
							Title: "Test",
							Type:  "blob",
						},
					},
				},
			},
			nil,
		},
		{
			"success many links",
			bytes.NewBufferString(`{"_links": {"self": {"href": "/","title": "Test"}, "test1": {"href": "/test/1", "title": "Test 1"}, "test2": {"href": "/test/2", "title": "Test 2"}},"foo": "bar"}`).Bytes(),
			&HALBody{
				Links: map[string][]HALLink{
					"self": []HALLink{
						HALLink{
							HRef:  "/",
							Title: "Test",
						},
					},
					"test1": []HALLink{
						HALLink{
							HRef:  "/test/1",
							Title: "Test 1",
						},
					},
					"test2": []HALLink{
						HALLink{
							HRef:  "/test/2",
							Title: "Test 2",
						},
					},
				},
			},
			nil,
		},
		{
			"success many item links",
			bytes.NewBufferString(`{"_links": {"self": {"href": "/","title": "Test"}, "item": [{"href": "/test/1", "title": "Test 1"}, {"href": "/test/2", "title": "Test 2"}]},"foo": "bar"}`).Bytes(),
			&HALBody{
				Links: map[string][]HALLink{
					"self": []HALLink{
						HALLink{
							HRef:  "/",
							Title: "Test",
						},
					},
					"item": []HALLink{
						HALLink{
							HRef:  "/test/1",
							Title: "Test 1",
						},
						HALLink{
							HRef:  "/test/2",
							Title: "Test 2",
						},
					},
				},
			},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := HALBody{
				Links:      map[string][]HALLink{},
				Properties: map[string]interface{}{},
				Embedded:   map[string]HALBody{},
			}
			if err := b.UnmarshalJSON(tt.d); err != nil {
				t.Errorf("HALBody.UnmarshalJSON() should not error, got %v", err)
			}

			for k, v := range b.Links {
				for i := 0; i < len(v); i++ {
					if tt.want.Links[k][i].HRef != v[i].HRef {
						t.Errorf("HALBody.UnmarshalJSON() Link %v.HRef expected '%v', got '%v'", k, tt.want.Links[k][i].HRef, v[i].HRef)
					}
					if tt.want.Links[k][i].Name != v[i].Name {
						t.Errorf("HALBody.UnmarshalJSON() Link %v.Name expected '%v', got '%v'", k, tt.want.Links[k][i].Name, v[i].Name)
					}
					if tt.want.Links[k][i].Title != v[i].Title {
						t.Errorf("HALBody.UnmarshalJSON() Link %v.Title expected '%v', got '%v'", k, tt.want.Links[k][i].Title, v[i].Title)
					}
					if tt.want.Links[k][i].Type != v[i].Type {
						t.Errorf("HALBody.UnmarshalJSON() Link %v.Type expected '%v', got '%v'", k, tt.want.Links[k][i].Type, v[i].Type)
					}
				}

			}
		})
	}
}
