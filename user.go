package insta

import (
	"bytes"
	"encoding/json"
	"github.com/valyala/fasthttp"
	"insta/models"
	"io/ioutil"
	"log"
	"os"
)

func hook(n models.PurpleNode) {

	img := n.DisplayURL
	log.Println(n.Owner.Username, n.Owner.ID, n.ID, img)
	getIMG(n.Owner.Username, n.Owner.ID, n.ID, img)

	//res := fmt.Sprintf("%s\t%s\t%s\t%d\t%d\t%d\t%s\t%s\t%v\t%s\t%s",
	//	n.Typename, // GraphVideo
	//	n.ID,
	//	"https://www.instagram.com/p/"+n.Shortcode+"/", // B8Ed2ghhMco // https://www.instagram.com/p/B8Ed2ghhMco/
	//	n.EdgeMediaToComment.Count,                     // 10
	//	n.TakenAtTimestamp,                             // 1580657992
	//	n.EdgeLikedBy.Count, // 120
	//	//n.EdgeMediaPreviewLike.Count,                   // 0
	//	n.Owner.ID, // 7062024874
	//	n.Owner.Username,                               //
	//	n.IsVideo, // false
	//	//n.Location,                                     //
	//	//n.GatingInfo,                                   //
	//	n.AccessibilityCaption, // Image may contain: 1 person
	//	n.EdgeMediaToCaption.Edges[0].Node.Text, // She’s teaching daddy a thing or two Snowboard session at Absolut Park in Austria! 🏂 		#rodeoand5th #travel 📹 by: @grilo
	//)

	f, _ := os.OpenFile("insta_detail.json", os.O_APPEND|os.O_CREATE, 0666)
	defer f.Close()
	body, _ := json.Marshal(n)
	f.Write(body)
	f.Write([]byte("\n"))
	log.Println(string(body))

}

func GetPostsByUserList(list string) {

	body, err := ioutil.ReadFile(list)
	if err != nil {
		log.Println("[E]", err)
		return
	}
	lines := bytes.Split(body, []byte("\r\n"))
	if len(lines) == 1 {
		lines = bytes.Split(body, []byte("\n"))
	}

	for _, v := range lines {
		v = bytes.Trim(v, "@")
		v = bytes.TrimSpace(v)
		GetPostsByUser(string(v))
	}
}

func GetPostsByUser(user string) {

	queryHash, insta := getFirstUserPage(user)
	gql := insta.EntryData.ProfilePage[0].Graphql
	src := gql.User
	getUserPageByScroll(queryHash, src, 0, 7000)
}

func getFirstUserPage(user string) (string, *models.Insta) {

	u := END_POINT + user
	code, body, err := fasthttp.Get(nil, u)
	if err != nil || code != 200 {
		log.Panicln("[E]", code, err)
	}
	queryHash := getQueryHash(body, Page)

	jsonBody := getJSONFromBody(body)
	var insta models.Insta
	err = json.Unmarshal(jsonBody, &insta)
	if err != nil {
		log.Panicln("[E]", err)
	}
	return queryHash, &insta
}

func getUserPageByScroll(queryHash string, o models.User, count, total int) {

	src := o.EdgeOwnerToTimelineMedia
	{ // doing something
		for _, v := range src.Edges {
			n := v.Node
			hook(n)
			//log.Println("[D] Node info:",
			//	n.Typename, // GraphVideo
			//	//n.EdgeMediaToCaption.Edges[0].Node.Text, // She’s teaching daddy a thing or two Snowboard session at Absolut Park in Austria! 🏂 		#rodeoand5th #travel 📹 by: @grilo
			//	"https://www.instagram.com/p/"+n.Shortcode+"/", // B8Ed2ghhMco // https://www.instagram.com/p/B8Ed2ghhMco/
			//	n.EdgeMediaToComment.Count,                     // 10
			//	//n.TakenAtTimestamp,                             // 1580657992
			//	n.EdgeLikedBy.Count, // 120
			//	//n.EdgeMediaPreviewLike.Count,                   // 0
			//	n.Owner.ID, // 7062024874
			//	//n.Owner.Username,                               //
			//	n.IsVideo, // false
			//	//n.Location,                                     //
			//	//n.GatingInfo,                                   //
			//	n.AccessibilityCaption, // Image may contain: 1 person
			//)

			//img := n.DisplayURL
			//log.Println(n.Owner.Username, n.Owner.ID, n.ID, img)
			//getIMG(n.Owner.Username, n.Owner.ID, n.ID, img)
		}
	}

	//NextScroll
	if src.PageInfo.HasNextPage && (total == -1 || (count+12) < total) {
		count += 50
		next := GetNextScroll(queryHash, "id", src.Edges[0].Node.Owner.ID, count, *src.PageInfo.EndCursor)
		getUserPageByScroll(queryHash, next.Data.User, count, total)
	}
}
