package fdsys

import (
	"encoding/xml"
	"fmt"
	"strings"

	welawproto "github.com/welaw/welaw/proto"
)

const (
	DisplayInline     = "yes-display-inline"
	DontDisplayInline = "no-display-inline"
)

var (
	mixedTextNames = []string{"quote", "header-in-text"}
	// TODO
	houseTypes  = []string{"hr", "hrres", "hjres", "hconres"}
	senateTypes = []string{"s", "sres", "sjres", "sconres"}
)

type LawMods struct {
	XMLName xml.Name `xml:"mods"`
	//SearchTitle string   `xml:"extension>searchTitle"`
	Date        string        `xml:"originInfo>dateIssued"`
	BillNumber  string        `xml:"extension>billNumber"`
	BillVersion string        `xml:"extension>billVersion"`
	CongMembers []*CongMember `xml:"extension>congMember"`
}

func (m *LawMods) GetSponsor() string {
	sponsor := ""
	for _, c := range m.CongMembers {
		if c.Role == "SPONSOR" {
			if sponsor != "" {
				panic(fmt.Errorf("more than one sponsor"))
			}
			sponsor = c.BioGuideID
		}
	}
	return sponsor
}

type CongMember struct {
	XMLName    xml.Name `xml:"congMember"`
	BioGuideID string   `xml:"bioGuideId,attr"`
	Chamber    string   `xml:"chamber,attr"`
	Congress   string   `xml:"congress,attr"`
	Role       string   `xml:"role,attr"`
	State      string   `xml:"state,attr"`
}

type LawXML interface {
	ToLaw() *welawproto.LawSet
	ParseItem(*Item)
}

type Sponsor struct {
	XMLName    xml.Name `xml:"sponsor"`
	BioGuideID string   `xml:"name-id,attr"`
}

type Title struct {
	XMLName xml.Name `xml:"title"`
	Enum    string   `xml:"enum"`
	Header  string   `xml:"header"`
}

func findIdent(loc string, items []*Item) int {
	ident := LocToIdent(loc)
	var found string
	for i, item := range items {
		if item.Loc == loc {
			continue
		}
		found = LocToIdent(item.Loc)
		if ident == found {
			return i
		}
	}
	return -1
}

func trimInnerWhitespace(s string) string {
	if len(s) == 0 {
		return ""
	}
	before := ""
	if string(s[0]) == " " {
		before = " "
	}
	after := ""
	if string(s[len(s)-1]) == " " {
		after = " "
	}
	n := strings.Join(strings.Fields(s), " ")
	n = before + n + after
	return n
}

//func trimWhitespace(in string) (out string) {
//white := false
//for _, c := range in {
//if unicode.IsSpace(c) {
//if !white {
//out = out + " "
//}
//white = true
//} else {
//out = out + string(c)
//white = false
//}
//}
//return
//}
