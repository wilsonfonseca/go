package hal

import (
	"net/url"
	"strconv"
)

// BasePage represents the simplest page: one with no links and only embedded records.
// Can be used to build custom page-like resources
type BasePage struct {
	BaseURL  *url.URL `json:"-"`
	Embedded struct {
		Records []Pageable `json:"records"`
	} `json:"_embedded"`
}

// Add appends the provided record onto the page
func (p *BasePage) Add(rec Pageable) {
	p.Embedded.Records = append(p.Embedded.Records, rec)
}

// Init initialized the Records slice.  This ensures that an empty page
// renders its records as an empty array, rather than `null`
func (p *BasePage) Init() {
	if p.Embedded.Records == nil {
		p.Embedded.Records = make([]Pageable, 0, 1)
	}
}

// Page represents the common page configuration (i.e. has self, next, and prev
// links) and has a helper method `PopulateLinks` to automate their
// initialization.
type Page struct {
	Links struct {
		Self Link `json:"self"`
		Next Link `json:"next"`
		Prev Link `json:"prev"`
	} `json:"_links"`

	BasePage
	BasePath    string     `json:"-"`
	Order       string     `json:"-"`
	Limit       uint64     `json:"-"`
	Cursor      string     `json:"-"`
	QueryParams url.Values `json:"-"`
}

func ModifyUrlParam(u *url.URL, key string, val string) *url.URL {
	q := u.Query()
	q.Del(key)
	q.Add(key, val)
	u.RawQuery = q.Encode()
	return u
}

// PopulateLinks sets the common links for a page.
func (p *Page) PopulateLinks() {
	p.Init()

	rec := p.Embedded.Records

	newUrl := p.BaseURL

	//verify paging params
	ModifyUrlParam(newUrl, "cursor", p.Cursor)
	ModifyUrlParam(newUrl, "order", p.Order)
	ModifyUrlParam(newUrl, "limit", strconv.FormatInt(int64(p.Limit), 10))

	//self: re-encode existing query params
	p.Links.Self = NewLink(newUrl.String())

	//next: update cursor to last record (if any)
	if len(rec) > 0 {
		ModifyUrlParam(newUrl, "cursor", rec[len(rec)-1].PagingToken())
	}
	p.Links.Next = NewLink(newUrl.String())

	//prev: inverse order and update cursor to first record (if any)
	ModifyUrlParam(newUrl, "order", p.InvertedOrder())
	if len(rec) > 0 {
		ModifyUrlParam(newUrl, "cursor", rec[0].PagingToken())
	}
	p.Links.Prev = NewLink(newUrl.String())
}

// InvertedOrder returns the inversion of the page's current order. Used to
// populate the prev link
func (p *Page) InvertedOrder() string {
	switch p.Order {
	case "asc":
		return "desc"
	case "desc":
		return "asc"
	default:
		return "asc"
	}
}
