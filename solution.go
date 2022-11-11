package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"strings"

	utopiago "github.com/Sagleft/utopialib-go"
	"github.com/badoux/goscraper"
	"github.com/sagleft/go-reddit/v2/reddit"
)

func main() {
	sol, err := newSolution()
	if err != nil {
		log.Fatalln(err)
	}

	err = sol.do()
	if err != nil {
		log.Fatalln(err)
	}
}

/*
           _       _   _
          | |     | | (_)
 ___  ___ | |_   _| |_ _  ___  _ __
/ __|/ _ \| | | | | __| |/ _ \| '_ \
\__ \ (_) | | |_| | |_| | (_) | | | |
|___/\___/|_|\__,_|\__|_|\___/|_| |_|

*/

func newSolution() (*solution, error) {
	sol := solution{}

	// parse args
	err := sol.parseArgs()
	if err != nil {
		return nil, err
	}

	// parse config file
	err = parseConfig(&sol)
	if err != nil {
		return nil, err
	}

	// get cache
	sol.Cache, err = NewCacheHandler(cacheFolderPath)
	if err != nil {
		return nil, err
	}

	// create utopia obj
	sol.Utopia = newUtopiaService().setToken(sol.Config.Utopia.Token).
		setHost(sol.Config.Utopia.Host).setPort(sol.Config.Utopia.Port).
		setHTTPS(sol.Config.Utopia.HTTPSEnabled)

	err = sol.Utopia.connect()
	if err != nil {
		return nil, err
	}

	return &sol, nil
}

func (sol *solution) parseArgs() error {
	fromSubreddits := flag.String("subreddit", "facepalm", "subbredit to crawl posts")
	channelID := flag.String("channel", "your channel id", "utopia channelID to export posts")
	flag.Parse()
	if fromSubreddits == nil {
		return errors.New("failed to get -subreddit arg")
	}

	sol.Config.FromSubreddits = strings.Split(*fromSubreddits, ",")
	if channelID == nil {
		return errors.New("failed to get -channel arg")
	}

	if *channelID == "" {
		return errors.New("-channel arg is empty")
	}
	sol.Config.UtopiaChannelID = *channelID
	return nil
}

func (sol *solution) isJoinedToChannel(channelID string) (bool, error) {
	channels, err := sol.Utopia.Client.GetChannels(utopiago.GetChannelsTask{
		SearchFilter: channelID,
		ChannelType:  utopiago.ChannelTypeJoined,
	})
	if err != nil {
		return false, err
	}

	return len(channels) > 0, nil
}

func (sol *solution) do() error {
	isJoined, err := sol.isJoinedToChannel(sol.Config.UtopiaChannelID)
	if err != nil {
		return err
	}
	if !isJoined {
		if _, err := sol.Utopia.Client.JoinChannel(sol.Config.UtopiaChannelID, ""); err != nil {
			return err
		}
	}

	credentials := reddit.Credentials{
		ID:       sol.Config.Reddit.APIKeyID,
		Secret:   sol.Config.Reddit.APISecret,
		Username: sol.Config.Reddit.User,
		Password: sol.Config.Reddit.Password,
	}
	client, err := reddit.NewClient(credentials)
	if err != nil {
		return errors.New("failed to connect to reddit: " + err.Error())
	}

	posts, _, err := client.Subreddit.TopPosts(
		context.Background(),
		GetRandomArrString(sol.Config.FromSubreddits),
		&reddit.ListPostOptions{
			ListOptions: reddit.ListOptions{
				Limit: subredditPostsLimit,
			},
			Time: getSubredditPostsBy,
		},
	)
	if err != nil {
		return err
	}

	postsUsedInQuery := 0
	for _, post := range posts {
		if sol.processPost(post) {
			postsUsedInQuery++
		}

		if postsUsedInQuery == sol.Config.PostsPerQuery {
			// all need posts used in this query
			return nil
		}
	}
	return nil
}

func (sol *solution) processPost(post *reddit.Post) bool {
	if sol.Cache.IsPostUsed(sol.Config.UtopiaChannelID, post.ID) {
		return false
	}

	var postImageURL string
	if isPhotoInURL(post.URL) {
		postImageURL = post.URL
	} else {
		// try find image in webpreview
		scraped, err := goscraper.Scrape(post.URL, 2)
		if err != nil {
			log.Println("failed to scrape webpreview for post " + post.ID)
			return false
		}
		scrapedImages := scraped.Preview.Images
		if len(scrapedImages) == 0 {
			return false
		}
		postImageURL = scrapedImages[0]
	}
	if postImageURL == "" {
		log.Println("post " + post.ID + " image is not recognized")
		return false
	}

	err := sol.Cache.MarkPostUsed(sol.Config.UtopiaChannelID, post.ID)
	if err != nil {
		log.Println("Failed to mark post used: " + err.Error())
		return false
	}

	if !isRemoteFileExists(postImageURL) {
		log.Println("remote image does not exists: " + postImageURL)
		return false
	}

	isDebug := false
	//sourceLink := html.A{Value: "[Source]", URL: "https://www.reddit.com" + post.Permalink}
	sourceLink := "https://www.reddit.com" + post.Permalink
	//postText := "<b>" + post.Title + "</b> " + sourceLink.Html()
	postText := post.Title
	if sol.Config.ShowSource {
		postText += "\n\n" + sourceLink
	}

	if !isDebug {

		err := sol.Utopia.postMedia(sol.Config.UtopiaChannelID, mediaPost{
			Text:         postText,
			ImageURL:     postImageURL,
			IsLocalImage: false,
		})
		if err != nil {
			log.Println("Failed to send photo to channel: " + err.Error())
		}

	} else {
		log.Println("debug, post ID: " + post.ID)
	}
	return true
}
