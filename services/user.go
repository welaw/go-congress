package services

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/net/html"

	"github.com/welaw/go-congress/congress"
	"github.com/welaw/go-congress/fdsys"
	"github.com/welaw/go-congress/pkg/errs"
	"github.com/welaw/go-congress/proto"
	welawproto "github.com/welaw/welaw/proto"
)

var (
	lisRegex      = regexp.MustCompile(`^[a-zA-Z]{1}\d{3}$`)
	bioguideRegex = regexp.MustCompile(`^[a-zA-Z]{1}\d{6}$`)
)

func (svc service) ensureUsers(ctx context.Context, bioguides []string) (users []*proto.User, err error) {
	for _, bioguide := range bioguides {
		user, err := svc.getUserByBioguideID(bioguide)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	// create user
	var welawUsers []*welawproto.User
	for _, u := range users {
		welawUsers = append(welawUsers, temp(u))
	}
	// TODO
	_, err = svc.client.CreateUsers(context.Background(), welawUsers)
	fmt.Printf("sent create users: %v\n", len(welawUsers))
	if err != nil {
		fmt.Printf("err: %v\n", err)
		if err.Error() != "409 Conflict" && strings.HasPrefix(err.Error(), "[409 Conflict]") {
			return nil, err
		}
	}
	for _, u := range users {
		if u.BioguideId != "" {
			url := fmt.Sprintf(congress.BioguidePictureURL, u.BioguideId[0], u.BioguideId)
			base := filepath.Base(url)
			pictureBytes, err := svc.getFile(url)
			if err != nil {
				return nil, err
			}
			err = svc.client.UploadAvatar(context.Background(), &welawproto.UploadAvatarOptions{
				Image:    pictureBytes,
				Filename: base,
				Username: u.Username,
			})
			if err != nil {
				return nil, err
			}
		}
		fmt.Println("saving user in db")
		err = svc.db.CreateUser(u)
		if err != nil {
			return nil, err
		}
	}
	return users, nil
}

// groupUsersFromItems takes groups of items and selects the users not in the database
func (svc service) groupUsersFromItems(groups [][]*fdsys.Item) (g [][]string, err error) {
	var group []string
	var done []string
	count := 0
	size := 20
	for _, items := range groups {
		for _, item := range items {
			bioguide := item.BioguideID
			if contains(done, bioguide) {
				continue
			}
			done = append(done, bioguide)

			// add to results and create new group if over size
			if count >= size {
				g = append(g, group)
				group = []string{}
				count = 0
			}

			// if user doesn't exist, add to the results
			_, err := svc.db.GetUserByBioguideId(bioguide)
			if err == errs.ErrNotFound {
				group = append(group, bioguide)
				count++
			} else if err != nil {
				return nil, err
			}
		}
	}

	if count > 0 {
		g = append(g, group)
	}
	return g, nil
}

func (svc service) getSenators() (*congress.Senators, error) {
	var s congress.Senators
	f, err := svc.getFile(congress.SenateMemberDataURL)
	if err != nil {
		return nil, err
	}
	err = xml.Unmarshal(f, &s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (svc service) getUser(id string) (*proto.User, error) {
	// lisRegex := regexp.MustCompile(`^[a-zA-Z]{1}\d{3}$`)
	// bioguideRegex := regexp.MustCompile(`^[a-zA-Z]{1}\d{6}$`)

	// check db

	switch {
	case lisRegex.MatchString(id):
		return svc.getUserByLisID(id)
	case bioguideRegex.MatchString(id):
		return svc.getUserByBioguideID(id)
	}
	return nil, fmt.Errorf("unknown user_id: %v", id)
}

func (svc service) getUserByLisID(lisID string) (*proto.User, error) {
	senators, err := svc.getSenators()
	if err != nil {
		return nil, err
	}
	bioguideID := senators.GetBioguideByLisID(lisID)
	if bioguideID == "" {
		return nil, errs.ErrNotFound
	}
	user, err := svc.getUserByBioguideID(bioguideID)
	if err != nil {
		return nil, err
	}
	user.LisId = lisID
	return user, nil
}

func (svc service) getUserByBioguideID(bioguideID string) (user *proto.User, err error) {
	if bioguideID == "" {
		return defaultUser(), nil
	}
	url := fmt.Sprintf(congress.BioguideURL, bioguideID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	var fullName, biography string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "font" && len(n.Attr) > 1 && n.NextSibling != nil {
			if n.Attr[0].Val == "4" {
				names := strings.Split(n.FirstChild.Data, ",")
				firstName := strings.Trim(names[1], " ")
				lastName := strings.Trim(names[0], " ")
				lastName = strings.Title(strings.ToLower(lastName))
				fullName = fmt.Sprintf("%s %s", firstName, lastName)
				biography = strings.Trim(n.NextSibling.Data, "\n")
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return &proto.User{
		Username:     bioguideID,
		FullName:     fullName,
		Email:        fmt.Sprintf("congress-%s@welaw.org", bioguideID),
		EmailPrivate: false,
		Biography:    biography,
		PictureUrl:   fmt.Sprintf("https://storage.googleapis.com/welaw-static/avatar/congress/%s.jpg", bioguideID),
		Upstream:     upstreamIdent,
		Provider:     upstreamProvider,
		BioguideId:   bioguideID,
	}, nil
}

func defaultUser() *proto.User {
	return &proto.User{
		BioguideId: "",
		Username:   "committee",
		Email:      "congress-committee@welaw.org",
		FullName:   "Committee",
		Biography:  "",
		PictureUrl: upstreamPictureURL,
		Upstream:   upstreamIdent,
		Provider:   upstreamProvider,
	}
}

func temp(u *proto.User) *welawproto.User {
	return &welawproto.User{
		Username:        u.Username,
		Email:           u.Email,
		EmailPrivate:    true,
		FullName:        u.FullName,
		FullNamePrivate: false,
		Biography:       u.Biography,
		PictureUrl:      u.PictureUrl,
		Upstream:        u.Upstream,
		Provider:        u.Provider,
	}
}

func contains(arr []string, tar string) bool {
	for _, s := range arr {
		if s == tar {
			return true
		}
	}
	return false
}
