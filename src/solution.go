package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
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

	subreddit := GetRandomArrString(sol.Config.FromSubreddits)
	fmt.Println("use subreddit: " + subreddit)

	posts, _, err := client.Subreddit.TopPosts(
		context.Background(),
		subreddit,
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
		isPostUsed, err := sol.processPost(post, subreddit)
		if err != nil {
			log.Fatalln(err)
		}
		if isPostUsed {
			postsUsedInQuery++
		}

		if postsUsedInQuery == sol.Config.MaxPostsPerQuery ||
			postsUsedInQuery == sol.Config.UsePostsPerQuery {
			fmt.Printf("relevant posts not found (ignored %v)\n", postsUsedInQuery)
			return nil
		}
	}
	return nil
}

func getRedditURL(url string) string {
	if strings.Contains(url, "http://") || strings.Contains(url, "https://") {
		return url
	}

	return redditHost + url
}

func (sol *solution) processPost(post *reddit.Post, subreddit string) (bool, error) {
	if sol.Cache.IsPostUsed(sol.Config.UtopiaChannelID, post.ID, subreddit) {
		fmt.Printf("post %q already used\n", post.Title)
		return false, nil
	}

	postImageURL := ""
	postResourceURL := getRedditURL(post.URL)

	if isPhotoInURL(postResourceURL) {
		postImageURL = postResourceURL
	} else {
		// try find image in webpreview
		scraped, err := goscraper.Scrape(postResourceURL, 2)
		if err != nil {
			log.Printf("failed to scrape webpreview for post %v: %s\n", post.ID, err.Error())
			return false, nil
		}
		scrapedImages := scraped.Preview.Images
		if len(scrapedImages) == 0 {
			fmt.Printf("ignore post %q without images. trying URL: %q\n", post.Title, postResourceURL)
			return false, nil
		}
		postImageURL = scrapedImages[0]
	}
	if postImageURL == "" {
		log.Println("post " + post.ID + " image is not recognized")
		return false, nil
	}

	err := sol.Cache.MarkPostUsed(sol.Config.UtopiaChannelID, post.ID, subreddit)
	if err != nil {
		return false, fmt.Errorf("failed to mark post used: %w", err)
	}

	if !isRemoteFileExists(postImageURL) {
		return false, fmt.Errorf("remote image does not exists: %s", postImageURL)
	}

	//sourceLink := html.A{Value: "[Source]", URL: "https://www.reddit.com" + post.Permalink}
	sourceLink := redditHost + post.Permalink
	//postText := "<b>" + post.Title + "</b> " + sourceLink.Html()
	postText := post.Title
	if sol.Config.ShowSource {
		postText += "\n\n" + sourceLink
	}

	err = sol.Utopia.postMedia(sol.Config.UtopiaChannelID, mediaPost{
		Text:         postText,
		ImageURL:     postImageURL,
		IsLocalImage: false,
	})
	if err != nil {
		return false, fmt.Errorf("failed to send photo to channel: %w", err)
	}

	fmt.Printf("mark post %q as used\n", post.Title)
	return true, nil
}
