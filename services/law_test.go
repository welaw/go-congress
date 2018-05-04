package services

//func TestOutdatedItems(t *testing.T) {

//c := MockClient{}
//db := client.NewDB()
//var l log.Logger
//usgpo := NewUpstream(c, db, l)

//_, err := usgpo.GetItems(false, "2016")
//require.NoError(t, err)

////for _, item := range items {
////b := ToBill(item.Loc)
////}

//}

//func MakeLawWithBranch() *lawv1.Law {
//branches := make([]*model.Branch, 1)
//branches[0] = &model.Branch{
//Name: "sponsor",
//}

//l := &model.Law{
//Ident:    "testident",
//Upstream: "usgpo",
//Creator:  "testcreator",
//Title:    "testtitle",
//Branches: branches,
//}

//return l
//}

//func TestInsertVersion(t *testing.T) {

//law := MakeLawWithBranch()

//v1 := &model.Version{
//LastMod: "2016-02-12T03:15:00.014Z",
//Body:    "test law body",
//}

//v2 := &model.Version{
//LastMod: "2016-01-13T09:55:00.227Z",
//Body:    "test law body",
//}

//v3 := &model.Version{
//LastMod: "2016-05-13T09:55:00.227Z",
//Body:    "test law body",
//}

//}
