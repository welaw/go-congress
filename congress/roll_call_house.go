package congress

import (
	"encoding/xml"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/welaw/go-congress/pkg/errs"
	"github.com/welaw/go-congress/pkg/utils"
	welawproto "github.com/welaw/welaw/proto"
)

const (
	baseURL   = "clerk.house.gov/evs"
	yearIndex = "http://clerk.house.gov/evs/2017/index.asp"
	roll0     = "http://clerk.house.gov/evs/2017/ROLL_000.asp"
	roll      = "http://clerk.house.gov/evs/2017/ROLL_500.asp"
)

func RollcallURL(year, page string) string {
	return fmt.Sprintf("http://clerk.house.gov/evs/%s/%s.xml", year, page)
}

func YearIndexURL(year string) string {
	return "http://clerk.house.gov/evs/%s/index.asp"
}

type RecordedVote struct {
	XMLName    xml.Name    `xml:"recorded-vote"`
	Legislator *Legislator `xml:"legislator"`
	//Party      string   `xml:"legislator>party,attr"`
	//State      string   `xml:"legislator>state,attr"`
	//Role       string   `xml:"legislator>role,attr"`
	Vote string `xml:"vote"`
}

type Legislator struct {
	XMLName    xml.Name `xml:"legislator"`
	BioguideID string   `xml:"name-id,attr"`
	Party      string   `xml:"party,attr"`
	State      string   `xml:"state,attr"`
	Role       string   `xml:"role,attr"`
}

type RollcallVote struct {
	XMLName      xml.Name        `xml:"rollcall-vote"`
	VoteMetadata *VoteMetadata   `xml:"vote-metadata"`
	Votes        []*RecordedVote `xml:"vote-data>recorded-vote"`
}

type VoteMetadata struct {
	Congress   string
	Session    string
	Chamber    string
	LegisNum   string
	VoteResult string
}

func GetHouseRollCall(url string) (map[string]string, *welawproto.VoteResult, error) {
	dat, err := utils.DownloadData(url)
	if err != nil {
		return nil, nil, err
	}
	if len(dat) == 0 {
		return nil, nil, errs.ErrNotFound
	}
	var rollcall RollcallVote
	err = xml.Unmarshal(dat, &rollcall)
	if err != nil {
		return nil, nil, err
	}
	votes := make(map[string]string)
	for _, v := range rollcall.Votes {
		votes[v.Legislator.BioguideID] = v.Vote
	}
	t := time.Now()
	s := int64(t.Unix())
	n := int32(t.Nanosecond())
	when := &timestamp.Timestamp{Seconds: s, Nanos: n}
	vr := &welawproto.VoteResult{
		UpstreamGroup: "house",
		Result:        rollcall.VoteMetadata.VoteResult,
		PublishedAt:   when,
	}
	return votes, vr, nil
}
