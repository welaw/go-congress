package congress

import (
	"encoding/xml"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/welaw/go-congress/pkg/utils"
	welawproto "github.com/welaw/welaw/proto"
)

const (
	SenateMemberDataURL = "https://www.senate.gov/legislative/LIS_MEMBER/cvc_member_data.xml"
	rollcallMenuURL     = "https://www.senate.gov/legislative/LIS/roll_call_lists/vote_menu_%s_%s.xml"
	rollcallURL         = "https://www.senate.gov/legislative/LIS/roll_call_votes/vote%s%s/vote_%s_%s_%s.xml"
)

func GetSenateRollCallURL(congress, session, voteNumber string) string {
	return fmt.Sprintf(rollcallURL, congress, session, congress, session, voteNumber)
}

func GetSenateRollCallMenuURL(congress, session string) string {
	return fmt.Sprintf(rollcallMenuURL, congress, session)
}

type VoteSummary struct {
	XMLName      xml.Name `xml:"vote_summary"`
	Congress     string   `xml:"congress"`
	Session      string   `xml:"session"`
	CongressYear string   `xml:"congress_year"`
	Votes        []*Vote  `xml:"votes>vote"`
}

type Vote struct {
	XMLName    xml.Name `xml:"vote"`
	VoteNumber string   `xml:"vote_number"`
	VoteDate   string   `xml:"vote_date"`
	Issue      string   `xml:"issue"`
	Question   string   `xml:"question"`
	Result     string   `xml:"result"`
	Yays       string   `xml:"vote_tally>yeas"`
	Nays       string   `xml:"vote_tally>nays"`
	Title      string   `xml:"title"`
}

// vote page
type RollCallVote struct {
	XMLName              xml.Name  `xml:"roll_call_vote"`
	Congress             string    `xml:"congress"`
	Session              string    `xml:"session"`
	CongressYear         string    `xml:"congress_year"`
	VoteNumber           string    `xml:"vote_number"`
	VoteDate             string    `xml:"vote_date"`
	ModifyDate           string    `xml:"modify_date"`
	VoteQuestionText     string    `xml:"vote_question_text"`
	VoteDocumentText     string    `xml:"vote_document_text"`
	VoteResultText       string    `xml:"vote_result_text"`
	Question             string    `xml:"question"`
	VoteTitle            string    `xml:"vote_title"`
	MajorityRequirements string    `xml:"majority_requirements"`
	VoteResult           string    `xml:"vote_result"`
	Yays                 string    `xml:"count>yeas"`
	Nays                 string    `xml:"count>nays"`
	Present              string    `xml:"count>present"`
	Absent               string    `xml:"count>absent"`
	Members              []*Member `xml:"members>member"`
}

type Member struct {
	MemberFull  string `xml:"member_full"`
	LastName    string `xml:"last_name"`
	FirstName   string `xml:"first_name"`
	VoteCast    string `xml:"vote_cast"`
	LisMemberID string `xml:"lis_member_id"`
}

type Senators struct {
	XMLName        xml.Name   `xml:"senators"`
	LastUpdateDate string     `xml:"lastUpdate>date"`
	LastUpdateTime string     `xml:"lastUpdate>time"`
	Senators       []*Senator `xml:"senator"`
}

type Senator struct {
	XMLName    xml.Name     `xml:"senator"`
	LisID      string       `xml:"lis_member_id,attr"`
	FirstName  string       `xml:"name>first"`
	LastName   string       `xml:"name>last"`
	State      string       `xml:"state"`
	BioGuideID string       `xml:"bioguideId"`
	Committees []*Committee `xml:"committees>committee"`
}

type Committee struct {
	XMLName xml.Name `xml:"committee"`
	Code    string   `xml:"code,attr"`
	T       string   `xml:",chardata"`
}

func (s *Senators) GetBioguideByLisID(lisID string) string {
	for _, senator := range s.Senators {
		if senator.LisID == lisID {
			return senator.BioGuideID
		}
	}
	return ""
}

func GetSenateRollCall(congress, session, number string) (map[string]string, *welawproto.VoteResult, error) {
	// TODO try all sessions, should be only 2
	url := GetSenateRollCallURL(congress, session, number)
	dat, err := utils.DownloadData(url)
	if err != nil {
		return nil, nil, err
	}
	if len(dat) == 0 {
		err = fmt.Errorf("downloaded data length is 0")
		return nil, nil, err
	}

	var rollCallVote RollCallVote
	err = xml.Unmarshal(dat, &rollCallVote)
	if err != nil {
		return nil, nil, err
	}
	votes := make(map[string]string)
	for _, m := range rollCallVote.Members {
		votes[m.LisMemberID] = m.VoteCast
	}
	t := time.Now()
	s := int64(t.Unix())
	n := int32(t.Nanosecond())
	when := &timestamp.Timestamp{Seconds: s, Nanos: n}
	vr := &welawproto.VoteResult{
		UpstreamGroup: "senate",
		Result:        rollCallVote.VoteResult,
		PublishedAt:   when,
	}
	return votes, vr, nil
}
