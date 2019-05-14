package sitemap

import "encoding/xml"

// SiteMap is XML sitemap
type SiteMap struct {
	XMLName xml.Name     `xml:"urlset"`
	URLs    []SiteMapURL `xml:"url"`
}

// SiteMapURL is a sitemap URL
type SiteMapURL struct {
	Loc        string `xml:"loc"`
	LastMod    string `xml:"lastmod"`
	ChangeFreq string `xml:"changefreq,omitempty"`
}
