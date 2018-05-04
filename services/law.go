package services

import (
	"context"
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/golang/protobuf/ptypes"
	"github.com/welaw/go-congress/congress"
	"github.com/welaw/go-congress/fdsys"
	"github.com/welaw/go-congress/pkg/errs"
	"github.com/welaw/go-congress/proto"
	welawproto "github.com/welaw/welaw/proto"
)

const (
	numWorkers = 10
)

func (svc *service) SendLaw(ctx context.Context, req *proto.ItemRange) (*proto.SendLawReply, error) {
	var resp proto.SendLawReply

	if req.Limit < 1 {
		return &resp, nil
	}

	if req.StartDate == "" {
		req.StartDate = time.Now().AddDate(0, -6, 0).UTC().String()
	}
	if req.EndDate == "" {
		req.EndDate = time.Now().UTC().String()
	}
	startTime, err := dateparse.ParseAny(req.StartDate)
	if err != nil {
		return &resp, err
	}
	endTime, err := dateparse.ParseAny(req.EndDate)
	if err != nil {
		return &resp, err
	}
	if endTime.Sub(startTime).Seconds() < 0 {
		return &resp, fmt.Errorf("end_time must be after start_time: start_time=%v, end_time=%v", startTime, endTime)
	}

	// get all sitemap items matching criteria
	allItems, err := svc.getItems(startTime, endTime, int(req.Limit), req.Ident)
	if err != nil {
		return &resp, err
	}

	// group laws by ident
	m := make(map[string][]*fdsys.Item)
	for _, i := range allItems {
		ident := fdsys.LocToShortIdent(i.Loc)
		m[ident] = append(m[ident], i)
	}

	items, err := svc.updateItemsModData(m)
	if err != nil {
		return nil, err
	}

	// limit the items
	if len(items) > int(req.Limit) {
		items = items[:req.Limit]
	}

	users, err := svc.ensureUsersByItems(ctx, items)
	if err != nil {
		return nil, err
	}
	// assign users to items
	for _, u := range users {
		for _, g := range items {
			for _, i := range g {
				if i.BioguideID == u.BioguideId {
					i.User = u
				}
			}
		}
	}

	responses, err := svc.sendLawItems(ctx, items)
	resp.NewItems = responses
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

func (svc service) Status(ctx context.Context, req *proto.ItemRange) (status *proto.StatusReply, err error) {
	status = &proto.StatusReply{}

	if req.StartDate == "" {
		req.StartDate = time.Now().AddDate(-1, 0, 0).UTC().String()
	}
	if req.EndDate == "" {
		req.EndDate = time.Now().UTC().String()
	}
	startTime, err := dateparse.ParseAny(req.StartDate)
	if err != nil {
		return status, err
	}
	endTime, err := dateparse.ParseAny(req.EndDate)
	if err != nil {
		return status, err
	}
	if endTime.Sub(startTime).Seconds() < 0 {
		return status, fmt.Errorf("end_time must be after start_time: start_time=%v, end_time=%v", startTime, endTime)
	}

	items, err := svc.getItems(startTime, endTime, int(req.Limit), req.Ident)

	var c, ty, n, v, longIdent string
	var newItems []string
	for _, item := range items {
		c, ty, n, v = fdsys.ParseLoc(item.Loc)
		longIdent = fdsys.ToLongIdent(c, ty, n, v)
		do, err := svc.shouldDownload(item)
		if err != nil {
			return nil, err
		}
		if !do {
			continue
		}
		newItems = append(newItems, longIdent)
	}
	status.NewItems = newItems
	return status, nil
}

func (svc *service) getCollectionSitemap(url string) (*fdsys.CollectionSitemap, error) {
	var sm fdsys.CollectionSitemap
	f, err := svc.getFile(url)
	if err != nil {
		return nil, err
	}
	err = xml.Unmarshal(f, &sm)
	if err != nil {
		return nil, err
	}
	return &sm, nil
}

func (svc *service) getSitemap(url string) (*fdsys.Sitemap, error) {
	var sm fdsys.Sitemap
	f, err := svc.getFile(url)
	if err != nil {
		return nil, err
	}
	err = xml.Unmarshal(f, &sm)
	if err != nil {
		return nil, err
	}
	return &sm, nil
}

func (svc *service) getItems(startTime, endTime time.Time, limit int, ident string) (items []*fdsys.Item, err error) {
	if limit == 0 {
		return items, nil
	}
	primary, err := svc.getSitemap(fdsys.PrimarySitemapURL)
	if err != nil {
		return
	}
	var year *fdsys.Sitemap
	var collection *fdsys.CollectionSitemap
	var item *fdsys.Item
	for y := startTime.Year(); y <= endTime.Year(); y++ {
		item = primary.FindSitemapItem(strconv.Itoa(y))
		if item == nil {
			continue
		}
		if do, err := svc.shouldDownload(item); err != nil {
			return nil, err
		} else if !do {
			continue
		}
		year, err = svc.getSitemap(item.Loc)
		if err != nil {
			return nil, err
		}
		item = year.FindCollectionSitemapItem(strconv.Itoa(y), "BILLS")
		if item == nil {
			continue
		}
		if do, err := svc.shouldDownload(item); err != nil {
			return nil, err
		} else if !do {
			continue
		}
		collection, err = svc.getCollectionSitemap(item.Loc)
		if err != nil {
			return nil, err
		}
		for _, item = range collection.Items {
			c, t, n, _ := fdsys.ParseLoc(item.Loc)
			if ident != "" && ident != fdsys.ToIdent(c, t, n) {
				continue
			}
			if t != congress.HouseKey && t != congress.SenateKey {
				svc.logger.Log("method", "get_items", "skipping bill type", t)
				continue
			}
			if do, err := svc.shouldDownload(item); err != nil {
				return nil, err
			} else if !do {
				continue
			}
			items = append(items, item)
		}
	}
	return items, nil
}

func (svc *service) getLawMods(url string) (*fdsys.LawMods, error) {
	var sm fdsys.LawMods
	f, err := svc.getFile(url)
	if err != nil {
		return nil, err
	}
	err = xml.Unmarshal(f, &sm)
	if err != nil {
		return nil, err
	}
	return &sm, nil
}

func (svc *service) getBillRes(url string) (fdsys.LawXML, error) {
	var sm fdsys.BillRes
	f, err := svc.getFile(url)
	if err != nil {
		return nil, err
	}
	err = xml.Unmarshal(f, &sm)
	if err != nil {
		return nil, err
	}
	reader := strings.NewReader(sm.LegisBody.InnerXML)
	sm.Content = fdsys.DecodeLaw(reader)
	return &sm, nil
}

func (svc *service) getAmendmentDoc(url string) (fdsys.LawXML, error) {
	var sm fdsys.AmendmentDoc
	f, err := svc.getFile(url)
	if err != nil {
		return nil, err
	}
	err = xml.Unmarshal(f, &sm)
	if err != nil {
		return nil, err
	}
	reader := strings.NewReader(sm.Body.InnerXML)
	sm.Content = fdsys.DecodeLaw(reader)
	return &sm, nil
}

func (svc *service) updateGroupMetadata(items []*fdsys.Item) (err error) {
	for _, item := range items {
		//svc.logger.Log("method", "update_group_metadata", "item", fmt.Sprintf("%+v", item))

		item.Ident = fdsys.LocToShortIdent(item.Loc)
		modsURL := fdsys.ModsURL(item.Loc)
		lawMods, err := svc.getLawMods(modsURL)
		if err != nil {
			return err
		}
		t, err := time.Parse("2006-01-02", lawMods.Date)
		if err != nil {
			return err
		}
		when, err := ptypes.TimestampProto(t)
		if err != nil {
			return err
		}
		bioguideID := lawMods.GetSponsor()
		item.When = when
		item.BioguideID = bioguideID
	}
	return
}

func (svc *service) getGroup(ctx context.Context, items []*fdsys.Item) (laws []*welawproto.LawSet, err error) {
	for _, item := range items {
		c, ty, n, v := fdsys.ParseLoc(item.Loc)
		longIdent := fdsys.ToLongIdent(c, ty, n, v)
		loc := fdsys.BillURL(longIdent)

		var lawXML fdsys.LawXML
		lawXML, err = svc.getBillRes(loc)
		if err != nil {
			lawXML, err = svc.getAmendmentDoc(loc)
			if err != nil {
				return nil, err
			}
		}
		lawXML.ParseItem(item)

		set := lawXML.ToLaw()
		set.Version.PublishedAt = item.When

		// some items dont have user info, use default 'committee' user
		username := defaultUser().Username
		if item.BioguideID != "" {
			username = item.BioguideID
		}
		set.Author = &welawproto.Author{
			Username: username,
		}

		set.Law.Upstream = svc.upstream.Ident

		switch ty {
		case congress.HouseKey:
			set.Version.UpstreamGroup = congress.HouseIdent
		case congress.SenateKey:
			set.Version.UpstreamGroup = congress.SenateIdent
		default:
			return nil, fmt.Errorf("bill type not found: %v", v)
		}

		set.Version.Tags = map[string]string{
			"Congress":     c,
			"Bill Version": strings.ToUpper(v),
		}

		laws = append(laws, set)
	}
	return
}

// ensureUsers
func (svc *service) ensureUsersByItems(ctx context.Context, items [][]*fdsys.Item) (users []*proto.User, err error) {
	// get users that need to be created
	groups, err := svc.groupUsersFromItems(items)
	if err != nil {
		return nil, err
	}

	type Result struct {
		users []*proto.User
		err   error
	}

	input := make(chan []string, len(groups))
	output := make(chan *Result)
	donech := make(chan struct{})

	for w := 0; w < numWorkers; w++ {
		go func() {
			for group := range input {
				users, err := svc.ensureUsers(ctx, group)
				if err != nil {
					svc.logger.Log("error ensuring users", err)
					output <- &Result{err: err}
					break
				}
				output <- &Result{users: users}
			}
			donech <- struct{}{}
		}()
	}
	for _, group := range groups {
		input <- group
	}
	close(input)

	var errs []error
	for i := 0; i < numWorkers; {
		select {
		case r := <-output:
			if r.err != nil {
				errs = append(errs, r.err)
			} else {
				for _, u := range r.users {
					users = append(users, u)
				}
			}
		case <-donech:
			i++
		}
	}

	if len(errs) > 0 {
		return users, fmt.Errorf("%+v", errs)
	}
	return users, nil
}

func (svc *service) updateItemsModData(itemMap map[string][]*fdsys.Item) (results [][]*fdsys.Item, err error) {

	input := make(chan []*fdsys.Item, len(itemMap))
	type Result struct {
		items []*fdsys.Item
		err   error
	}
	output := make(chan *Result)
	donech := make(chan struct{})

	for i := 0; i < numWorkers; i++ {
		go func() {
			for group := range input {
				err = svc.updateGroupMetadata(group)
				if err != nil {
					svc.logger.Log("error", fmt.Errorf("error grouping metadata --- %v", err))
					output <- &Result{err: err}
					break
				}
				output <- &Result{items: group}
			}
			donech <- struct{}{}
		}()
	}
	for _, items := range itemMap {
		input <- items
	}
	close(input)

	var errs []error
	for i := 0; i < numWorkers; {
		select {
		case r := <-output:
			if r.err != nil {
				errs = append(errs, err)
			} else if r.items != nil {
				results = append(results, r.items)
			}
		case <-donech:
			i++
		}
	}

	if len(errs) > 0 {
		return results, fmt.Errorf("%+v", errs)
	}
	return results, nil
}

func (svc *service) sendLawItems(ctx context.Context, items [][]*fdsys.Item) (results []string, err error) {

	type Result struct {
		msgs []string
		err  error
	}

	input := make(chan []*fdsys.Item, len(items))
	output := make(chan *Result)
	donech := make(chan struct{})

	for w := 0; w < numWorkers; w++ {
		go func() {
			var done []string
			for group := range input {
				laws, err := svc.getGroup(ctx, group)
				if err != nil {
					svc.logger.Log("error getting items", err)
					output <- &Result{err: err}
					break
				}

				_, err = svc.client.CreateLaws(ctx, laws, &welawproto.CreateLawsOptions{})
				if err != nil {
					svc.logger.Log("error creating laws", err)
					output <- &Result{err: err}
					break
				}
				for _, i := range group {
					err = svc.db.CreateItem(i)
					if err != nil {
						svc.logger.Log("error creating item in db", err)
						output <- &Result{err: err}
						donech <- struct{}{}
						return
					}
					done = append(done, i.Loc)
				}
				output <- &Result{msgs: done}
			}
			donech <- struct{}{}
		}()
	}
	for _, i := range items {
		input <- i
	}
	close(input)

	var errs []error
	for i := 0; i < numWorkers; {
		select {
		case r := <-output:
			if r.err != nil {
				errs = append(errs, err)
			} else {
				results = append(results, r.msgs...)
			}
		case <-donech:
			i++
		}
	}

	if len(errs) > 0 {
		return results, fmt.Errorf("%+v", errs)
	}
	return results, nil
}

// shouldDownload checks for an item in the db to see if it needs
// to be downloaded
func (svc *service) shouldDownload(i *fdsys.Item) (bool, error) {
	if i == nil {
		return false, fmt.Errorf("should_download: item is nil")
	}
	item, err := svc.db.GetItemByLoc(i.Loc)
	switch {
	case err == errs.ErrNotFound:
		return true, nil
	case err != nil:
		return false, err
	}
	t1, err := time.Parse(time.RFC3339Nano, item.LastMod)
	if err != nil {
		return false, err
	}
	t2, err := time.Parse(time.RFC3339Nano, i.LastMod)
	if err != nil {
		return false, err
	}
	if t1 == t2 {
		return false, nil
	}
	return false, fmt.Errorf("should_download: case not handled: item=%v", i)
}

func (svc *service) updateSitemaps(force bool) ([]string, error) {
	return nil, nil
}

type Items [][]*fdsys.Item

func (l Items) Len() int      { return len(l) }
func (l Items) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
func (l Items) Less(i, j int) bool {
	var newest int64
	for _, item := range l[i] {
		if item.When.Seconds > newest {
			newest = item.When.Seconds
		}
	}
	for _, item := range l[j] {
		if item.When.Seconds > newest {
			return false
		}
	}
	return true
}
