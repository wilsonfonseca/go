package url

import (
	gUrl "net/url"
)

type URL gUrl.URL

// SetParam
func (u URL) SetParam(key string, val string) URL {
	gu := gUrl.URL(u)
	q := gu.Query()
	q.Del(key)
	q.Add(key, val)
	gu.RawQuery = q.Encode()
	return URL(gu)
}

func (u URL) String() string {
	gu := gUrl.URL(u)
	return gu.String()
}

func Parse(s string) (u URL, err error) {
	gu, err := gUrl.Parse(s)
	if err != nil {
		return
	}
	return URL(*gu), nil
}
