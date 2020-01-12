package link

import (
	"net/url"

	"github.com/yosida95/uritemplate"
)

type LinkSet map[string][]Link

func (s LinkSet) Add(k string, v Link) {
	s[k] = append(s[k], v)
}

func (s LinkSet) Set(k string, v []Link) {
	s[k] = v
}

func (s LinkSet) Get(k string) []Link {
	if !s.Has(k) {
		return []Link{}
	}

	return s[k]
}

func (s LinkSet) Values() []Link {
	l := []Link{}

	for _, v := range s {
		l = append(l, v...)
	}

	return l
}

func (s LinkSet) Has(k string) bool {
	_, ok := s[k]
	return ok
}

type Link struct {
	Context   string
	HRef      string
	Rel       string
	Name      string
	Templated bool
	Title     string
	Type      string
}

func (l Link) Expand(lv map[string]string) (string, error) {
	if !l.Templated {
		return l.Resolve()
	} else {
		t, err := uritemplate.New(l.Context)
		if err != nil {
			return "", err
		}

		tv := uritemplate.Values{}
		for k, v := range lv {
			tv.Set(k, uritemplate.String(v))
		}

		h, err := t.Expand(tv)
		if err != nil {
			return "", err
		}

		return h, nil
	}
}

func (l Link) Resolve() (string, error) {
	u, err := url.Parse(l.Context)
	if err != nil {
		return "", err
	}

	h, err := u.Parse(l.HRef)
	if err != nil {
		return "", err
	}

	return h.String(), nil
}
