package fdsys

import (
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/welaw/go-congress/proto"
)

const (
	PrimarySitemapURL = "https://www.gpo.gov/smap/fdsys/sitemap.xml"
)

// Sitemap covers the "Primary" and "Year" sitemaps.
type Sitemap struct {
	Items []*Item `xml:"sitemap"`
}

// CollectionSitemap covers BILLS.
type CollectionSitemap struct {
	Items []*Item `xml:"url"`
}

type Item struct {
	Loc         string `xml:"loc"`
	LastMod     string `xml:"lastmod"`
	Ident       string
	PublishedAt string
	BioguideID  string
	User        *proto.User
	When        *timestamp.Timestamp
}

func (sm Sitemap) FindCollectionSitemapItem(y string, l string) *Item {
	var base string
	for _, i := range sm.Items {
		base = filepath.Base(i.Loc)
		if base == CollectionSitemapName(y, l) {
			return i
		}
	}
	return nil
}

func (sm Sitemap) FindSitemapItem(y string) *Item {
	var base string
	for _, i := range sm.Items {
		base = filepath.Base(i.Loc)
		if base == YearSitemapName(y) {
			return i
		}
	}
	return nil
}

func YearSitemapName(y string) string {
	return fmt.Sprintf("sitemap_%s.xml", y)
}

func CollectionSitemapName(y string, l string) string {
	return fmt.Sprintf("%s_%s_sitemap.xml", y, l)
}

type ByLastMod []*Item

func (s ByLastMod) Len() int      { return len(s) }
func (s ByLastMod) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s ByLastMod) Less(i, j int) bool {
	it, err := time.Parse(time.RFC3339Nano, s[i].LastMod)
	if err != nil {
		panic(err)
	}
	jt, err := time.Parse(time.RFC3339Nano, s[j].LastMod)
	if err != nil {
		panic(err)
	}
	_, _, in, _ := ParseLoc(s[i].Loc)
	_, _, jn, _ := ParseLoc(s[j].Loc)
	inn, err := strconv.Atoi(in)
	if err != nil {
		panic(err)
	}
	jnn, err := strconv.Atoi(jn)
	if err != nil {
		panic(err)
	}
	switch {
	case inn < jnn:
		return true
	case inn == jnn:
		if jt.Sub(it).Seconds() <= 0 {
			return false
		}
		return true
	default:
		return false
	}
}

type ByActionDate []*Item

func (s ByActionDate) Len() int      { return len(s) }
func (s ByActionDate) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s ByActionDate) Less(i, j int) bool {
	it, err := time.Parse(time.RFC3339Nano, s[i].LastMod)
	if err != nil {
		panic(err)
	}
	jt, err := time.Parse(time.RFC3339Nano, s[j].LastMod)
	if err != nil {
		panic(err)
	}
	_, _, in, _ := ParseLoc(s[i].Loc)
	_, _, jn, _ := ParseLoc(s[j].Loc)
	inn, err := strconv.Atoi(in)
	if err != nil {
		panic(err)
	}
	jnn, err := strconv.Atoi(jn)
	if err != nil {
		panic(err)
	}
	switch {
	case inn < jnn:
		return true
	case inn == jnn:
		if jt.Sub(it).Seconds() <= 0 {
			return false
		}
		return true
	default:
		return false
	}
}
