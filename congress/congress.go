package congress

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	welawproto "github.com/welaw/welaw/proto"
)

const (
	HouseKey           = "hr"
	HouseIdent         = "house"
	SenateKey          = "s"
	SenateIdent        = "senate"
	lawURL             = "https://www.congress.gov/bill/%sth-congress/%s-bill/%s"
	actionsURL         = "https://www.congress.gov/bill/%sth-congress/%s-bill/%s/actions"
	BioguideURL        = "http://bioguide.congress.gov/scripts/biodisplay.pl?index=%s"
	BioguidePictureURL = "http://bioguide.congress.gov/bioguide/photo/%c/%s.jpg"
)

func ActionsURL(congress, chamber, billNumber string) string {
	return fmt.Sprintf(actionsURL, congress, chamber, billNumber)
}

func LawURL(congress, chamber, billNumber string) string {
	return fmt.Sprintf(lawURL, congress, chamber, billNumber)
}

func GetRollCalls(congress, chamber, session, billNumber string) (map[string]string, *welawproto.VoteResult, error) {
	url := ActionsURL(congress, chamber, billNumber)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, nil, err
	}

	m := make(map[string]string)
	var votes map[string]string
	var result *welawproto.VoteResult
	doc.Find(".item_table tr td:nth-child(2)").Each(func(i int, s *goquery.Selection) {
		t := s.Text()
		t = strings.TrimSpace(t)
		f := strings.Fields(t)
		// TODO
		if len(f) == 0 || f[0] != "Passed/agreed" {
			return
		}
		s.Find("a").Each(func(i int, s *goquery.Selection) {
			l, ok := s.Attr("href")
			if !ok {
				return
			}
			votes, result, err = handleChamberVote(congress, session, l)
			if err != nil {
				return
			}
			for k, v := range votes {
				m[k] = v
			}

		})
	})

	return m, result, nil
}

func handleChamberVote(congress, session, link string) (map[string]string, *welawproto.VoteResult, error) {
	senate := `\bvote\b=(\d{5})`
	r := regexp.MustCompile(senate)
	match := r.FindAllStringSubmatch(link, -1)
	if len(match) > 0 {
		return GetSenateRollCall(congress, session, match[0][1])
	}
	return GetHouseRollCall(link)
}
