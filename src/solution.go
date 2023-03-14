package main

import (
	"context"
	"fmt"
	"log"

	swissknife "github.com/Sagleft/swiss-knife"
	"github.com/badoux/goscraper"
	"github.com/sagleft/go-reddit/v2/reddit"
	"gopkg.in/robfig/cron.v2"
)

func main() {
	swissknife.PrintIntroMessage(botLogName, donateAddress, coinTag)

	if err := runApp(); err != nil {
		log.Fatalln(err)
	}

	log.Println("bot started")
	swissknife.RunInBackground()
}

func runApp() error {
	var err error
	sol := solution{}

	log.Println("load config..")
	sol.Config, err = parseConfig()
	if err != nil {
		return fmt.Errorf("parse config: %w", err)
	}

	sol.Cache, err = NewCacheHandler(cacheFolderPath)
	if err != nil {
		return fmt.Errorf("create cache handler: %w", err)
	}

	log.Println("parse content routes..")
	routes := parseContentRoutes(sol.Config.Main.Routes)

	log.Println("get channels..")
	chats := getChats(routes)

	log.Println("connect to Utopia Network..")
	if sol.Utopia, err = utopiaConnect(
		sol.Config.Utopia,
		sol.Config.Main.BotNickname,
		chats,
	); err != nil {
		return fmt.Errorf("connect to utopia: %w", err)
	}

	if err := sol.Utopia.updateAccountName(); err != nil {
		return fmt.Errorf("update account name: %w", err)
	}

	if err := sol.Utopia.loadBotPubkey(); err != nil {
		return fmt.Errorf("load bot pubkey: %w", err)
	}

	log.Println("connect to reddit..")
	if sol.Reddit, err = redditConnect(sol.Config.Reddit); err != nil {
		return fmt.Errorf("reddit connect: %w", err)
	}

	log.Println("setup cron..")
	if err := sol.setupCron(); err != nil {
		return fmt.Errorf("setup cron: %w", err)
	}
	return nil
}

func (sol *solution) setupCron() error {
	c := cron.New()
	c.AddFunc(parseCronSpec(sol.Config.Main.Cron), func() {
		err := sol.findAndPlacePost()
		if err != nil {
			log.Fatalln(err)
		}
	})
	c.Start()

	return nil
}

func (sol *solution) markPostProcessing(isProcessing bool) {
	sol.IsProcessingPost = isProcessing
}

func (sol *solution) findAndPlacePost() error {
	if sol.IsProcessingPost {
		log.Println("prevent multiple post processing. skip")
		return nil
	}

	sol.markPostProcessing(true)
	defer sol.markPostProcessing(false)

	fmt.Println()

	subreddit := GetRandomArrString(sol.FromSubreddits)
	fmt.Println("use subreddit: " + subreddit)

	ctx, ctxCancel := context.WithTimeout(context.Background(), getPostsTimeout)
	defer ctxCancel()

	posts, _, err := client.Subreddit.TopPosts(
		ctx,
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

	if len(posts) == 0 {
		log.Println("posts not found. ignore")
		return nil
	}

	postsUsedInQuery := 0
	for _, post := range posts {
		isPostUsed, err := sol.processPost(post, subreddit)
		if err != nil {
			return fmt.Errorf("process post: %w", err)
		}
		if isPostUsed {
			postsUsedInQuery++
		}

		if postsUsedInQuery == sol.Config.Main.MaxPostsPerQuery ||
			postsUsedInQuery == sol.Config.Main.UsePostsPerQuery {
			fmt.Printf("relevant posts not found (ignored %v)\n", postsUsedInQuery)
			return nil
		}
	}
	return nil
}

func (sol *solution) processPost(post *reddit.Post, subreddit string) (bool, error) {
	if sol.Cache.IsPostUsed(sol.Config.Main.UtopiaChannelID, post.ID, subreddit) {
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

	err := sol.Cache.MarkPostUsed(sol.Config.Main.UtopiaChannelID, post.ID, subreddit)
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
	if sol.Config.Main.ShowSource {
		postText += "\n\n" + sourceLink
	}

	err = sol.Utopia.postMedia(sol.Config.Main.UtopiaChannelID, mediaPost{
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
