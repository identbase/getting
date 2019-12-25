package getting

import (
	"errors"
	"fmt"
	"net/url"
	"testing"

	"github.com/identbase/getting/pkg/resource"
)

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
	client, _ := New("localhost:8000")
	endpoint, _ := url.Parse("localhost:8000/api/endpoint")
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

			fmt.Println("result", r, e)

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
