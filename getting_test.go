package getting

import (
	"encoding/json"
	"errors"
	"fmt"
	// "io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/identbase/getting/pkg/resource"
	"github.com/identbase/getting/pkg/resource/representor"
)

// Unit tests
func Test_Getting_New(t *testing.T) {
	tests := []struct {
		name string
		b    string
		want *Getting
		err  error
	}{
		{
			"success",
			"localhost:8000",
			&Getting{
				bookmark: "localhost:8000",
			},
			nil,
		},
		{
			"error bookmark unspecified",
			"",
			nil,
			errors.New("bookmark unspecified"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, e := New(tt.b)

			if tt.err != nil && e == nil {
				t.Errorf("getting.New() should error with %v, got %v", tt.err, e)
			} else if tt.err == nil && e != nil {
				t.Errorf("getting.New() errored with %v when it shouldnt have", e)
			}

			if tt.want != nil && g == nil {
				t.Errorf("getting.New() should create Getting object, got nil")
			} else if tt.want == nil && g != nil {
				t.Errorf("getting.New() should not have created Getting object but got one")
			} else if tt.want != nil && g != nil {
				if g.bookmark != tt.want.bookmark {
					t.Errorf("getting.New() = %v, want %v", g, tt.want)
				}
			}
		})
	}
}

func Test_Getting_Go_UrlParse(t *testing.T) {
	endpoint, _ := url.Parse("http://localhost:8000/api/endpoint")
	client, _ := New("http://localhost:8000")
	tests := []struct {
		name string
		u    string
		want *resource.Resource
		err  error
	}{
		{
			"success",
			"/api/endpoint",
			&resource.Resource{
				Client: client,
				URI:    endpoint,
			},
			nil,
		},
		{
			"error",
			fmt.Sprintf("%c", rune(0x7f)),
			nil,
			errors.New("net/url: invalid control character in URL"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, e := client.Go(tt.u)

			if tt.err != nil && e == nil {
				t.Errorf("getting.Go() should error, got nil")
			} else if tt.err == nil && e != nil {
				t.Errorf("getting.Go() errored with %v when it shouldnt have", e)
			}

			if tt.want != nil && r == nil {
				t.Errorf("getting.Go() should create Resource object, got nil")
			} else if tt.want == nil && r != nil {
				t.Errorf("getting.Go() should not have created Resource object but got one")
			} else if tt.want != nil && r != nil {
				if (*r).URI.String() != (*tt.want).URI.String() {
					t.Errorf("getting.Go() = %v, want %v", r, tt.want)
				}
			}
		})
	}
}

// Acceptance tests
func Test_Getting_SimpleResourceGet(t *testing.T) {
	tests := []struct {
		name string
		want interface{}
		err  error
	}{
		{
			"success",
			representor.HALBody{
				Links: map[string][]representor.HALLink{
					"self": []representor.HALLink{
						representor.HALLink{
							HRef:  "/",
							Title: "Test",
						},
					},
				},
				Properties: map[string]interface{}{
					"foo": "bar",
				},
				Embedded: map[string]representor.HALBody{},
			},
			nil,
		},
	}
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"_links": map[string]interface{}{
				"self": map[string]string{
					"href":  "/",
					"title": "Test",
				},
				"test": map[string]string{
					"href":  "/",
					"title": "WINNING",
				},
			},
			"foo": "bar",
		}

		js, err := json.Marshal(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/hal+json")
		w.Write(js)
	}))
	defer s.Close()

	g, err := New(s.URL)
	if err != nil {
		t.Error(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := g.Follow("test")
			if err != nil {
				t.Error(err)
			}

			o, err := r.Get()
			if err != nil {
				t.Error(err)
			}

			for k, v := range tt.want.(representor.HALBody).Links {
				if _, ok := o.(representor.HALBody).Links[k]; !ok {
					t.Error("object mismatch")
				}

				if o.(representor.HALBody).Links[k][0] != v[0] {
					t.Error("object mismatch")
				}
			}
		})
	}

}
