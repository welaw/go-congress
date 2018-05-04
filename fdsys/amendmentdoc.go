package fdsys

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	welawproto "github.com/welaw/welaw/proto"
)

type EngrossedAmendmentBody struct {
	XMLName  xml.Name `xml:"engrossed-amendment-body"`
	InnerXML string   `xml:",innerxml"`
}

type AmendmentDoc struct {
	XMLName xml.Name                `xml:"amendment-doc"`
	Body    *EngrossedAmendmentBody `xml:"engrossed-amendment-body"`

	BillLoc string
	//Publisher string `xml:"publisher"`
	// Example: 115
	Congress string // `xml:"form>congress"`
	//Session  string `xml:"form>session"`
	// hr, s, hjres, sjres, hconres, sconres, hres, sres
	BillType string
	// Example: 1027
	BillNumber string
	// as, cps, fph, lth, ppv, rds, rhv, rhuc, ash, eah, fps, lts, pap, rev, rih,
	// sc, eas, hdh, nat, pwah, reah, ris, ath, eh, hds, oph, rah, res, rsv, ats,
	// eph, ihv, ops, ras, renr, rth, cdh, enr, iph, pav, rch, rfh, rts, cds, esv,
	// ips, pch, rcs, rfs, s_p, cph, fah, isv, pcs, rdh, rft, sas, mostrecent
	BillVersion string
	ActionDate  string `xml:"amendment-doc>engrossed-amendment-form>action>action-date"`
	// LastMod
	When time.Time
	//LegisNum      string     `xml:"form>legis-num"`
	//Chamber       string     `xml:"form>current-chamber"`
	Sponsor       *Sponsor
	OfficialTitle string
	Content       string
}

func (br *AmendmentDoc) ParseItem(item *Item) {
	when, err := ParseLastMod(item.LastMod)
	if err != nil {
		panic(err)
	}
	br.When = when
	br.BillLoc = item.Loc
	c, t, n, v := ParseLoc(item.Loc)
	br.Congress = c
	br.BillType = t
	br.BillNumber = n
	br.BillVersion = v
	return
}

func (br AmendmentDoc) Ident() string {
	return ToIdent(br.Congress, br.BillType, br.BillNumber)
}

func (br AmendmentDoc) LongIdent() string {
	return ToLongIdent(br.Congress, br.BillType, br.BillNumber, br.BillVersion)
}

func (br AmendmentDoc) Title() string {
	return fmt.Sprintf("%s %s", strings.ToUpper(br.BillType), br.BillNumber)
}

func (br *AmendmentDoc) Summary() string {
	return br.Congress
}

func (br *AmendmentDoc) String() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("%s\n", br.ActionDate))

	return buf.String()
}

func (br *AmendmentDoc) ToLaw() *welawproto.LawSet {
	return &welawproto.LawSet{
		Law: &welawproto.Law{
			Title:      br.OfficialTitle,
			ShortTitle: br.Title(),
			Ident:      br.Ident(),
		},
		Branch: &welawproto.Branch{},
		Version: &welawproto.Version{
			Body: br.String() + br.Content,
			Msg:  br.Summary(),
		},
	}
}
