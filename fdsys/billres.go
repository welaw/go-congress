package fdsys

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	welawproto "github.com/welaw/welaw/proto"
)

type BillRes struct {
	BillLoc string
	//Publisher string `xml:"publisher"`
	//Date string `xml:"date"`
	XMLName xml.Name `xml:"bill"`
	// Example: 115
	//Date           string    `xml:"metadata>dublinCore>date"`
	Congress       string    `xml:"form>congress"`
	Session        string    `xml:"form>session"`
	LegisNum       string    `xml:"form>legis-num"`
	LegisType      string    `xml:"form>legis-type"`
	CurrentChamber string    `xml:"form>current-chamber"`
	OfficialTitle  string    `xml:"form>official-title"`
	Actions        []*Action `xml:"form>action"`
	// hr, s, hjres, sjres, hconres, sconres, hres, sres
	BillType string
	// Example: 1027
	BillNumber string
	// as, cps, fph, lth, ppv, rds, rhv, rhuc, ash, eah, fps, lts, pap, rev, rih,
	// sc, eas, hdh, nat, pwah, reah, ris, ath, eh, hds, oph, rah, res, rsv, ats,
	// eph, ihv, ops, ras, renr, rth, cdh, enr, iph, pav, rch, rfh, rts, cds, esv,
	// ips, pch, rcs, rfs, s_p, cph, fah, isv, pcs, rdh, rft, sas, mostrecent
	BillVersion string
	// LastMod
	When time.Time
	//LegisNum      string     `xml:"form>legis-num"`
	//Body      *LegisBody `xml:"legis-body,omit-empty"`
	LegisBody *LegisBody `xml:"legis-body"`
	Content   string
	//AttestationDate *AttestationDate `xml:"attestation>attestation-group>attestation-date"`
}

func (br *BillRes) GetLastAction() *Action {
	if br.Actions == nil {
		return nil
	}
	if len(br.Actions) == 0 {
		return nil
	}
	return br.Actions[len(br.Actions)-1]
}

type Action struct {
	XMLName       xml.Name     `xml:"action"`
	Stage         string       `xml:"stage,attr"`
	ActionDate    *ActionDate  `xml:"action-date"`
	Sponsor       *Sponsor     `xml:"action-desc>sponsor"`
	Cosponsors    []*Cosponsor `xml:"action-desc>cosponsor"`
	CommitteeName *Committee   `xml:"action-desc>committee-name"`
}

type ActionDate struct {
	XMLName xml.Name `xml:"action-date"`
	Date    string   `xml:"date,attr"`
	Value   string   `xml:",chardata"`
}

//type AttestationDate struct {
//XMLName xml.Name `xml:"attestation-date"`
//Date    string   `xml:"date,attr"`
//Chamber string   `xml:"chamber,attr"`
//Value   string   `xml:",innerchar"`
//}

type LegisBody struct {
	XMLName  xml.Name `xml:"legis-body"`
	Title    *Title   `xml:"title"`
	InnerXML string   `xml:",innerxml"`
}

type Cosponsor struct {
	XMLName xml.Name `xml:"cosponsor"`
	NameID  string   `xml:"name-id,attr"`
	T       string   `xml:",chardata"`
}

type Committee struct {
	XMLName xml.Name `xml:"committee-name"`
	T       string   `xml:",chardata"`
}

// move this elsewhere
func (br *BillRes) ParseItem(item *Item) {
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

func (br BillRes) Ident() string {
	return ToIdent(br.Congress, br.BillType, br.BillNumber)
}

func (br BillRes) LongIdent() string {
	return ToLongIdent(br.Congress, br.BillType, br.BillNumber, br.BillVersion)
}

func (br BillRes) Title() string {
	return fmt.Sprintf("%s %s", strings.ToUpper(br.BillType), br.BillNumber)
}

func (br *BillRes) String() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("## %sth Congress\n", br.Congress))
	buf.WriteString(fmt.Sprintf("## %s\n", br.Session))
	buf.WriteString(fmt.Sprintf("# %s\n", br.LegisNum))
	buf.WriteString(fmt.Sprintf("%s\n", trimInnerWhitespace(br.OfficialTitle)))
	buf.WriteString("\n---\n")
	// chamber
	buf.WriteString(fmt.Sprintf("### %s\n", br.CurrentChamber))
	ad := br.GetActionDate()
	if ad != "" {
		buf.WriteString(fmt.Sprintf("### %s\n", ad))
	}
	buf.WriteString("\n---\n")
	buf.WriteString(fmt.Sprintf("### %s\n", br.LegisType))
	buf.WriteString(fmt.Sprintf("%s\n", trimInnerWhitespace(br.OfficialTitle)))

	return buf.String()
}

func (br *BillRes) Summary() string {
	return br.Congress
}

func (br *BillRes) GetActionDate() string {
	lastAction := br.GetLastAction()
	if lastAction == nil {
		return ""
	}
	if lastAction.ActionDate == nil {
		return ""
	}
	return lastAction.ActionDate.Value
}

func (br *BillRes) GetSponsor() string {
	lastAction := br.GetLastAction()
	if lastAction == nil {
		return ""
	}
	if lastAction.Sponsor == nil {
		return ""
	}
	return lastAction.Sponsor.BioGuideID
}

func (br *BillRes) ToLaw() *welawproto.LawSet {
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
