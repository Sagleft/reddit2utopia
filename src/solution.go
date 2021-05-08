package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/badoux/goscraper"
	"github.com/sagleft/go-reddit/reddit/v2"
	"github.com/tabvn/html"
)

const (
	configJSONPath      = "../config/config.json"
	cacheFolderPath     = "../cache"
	getSubredditPostsBy = "day"
	subredditPostsLimit = 24
	postsPerQuery       = 1
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

type solutionConfig struct {
	BotToken            string `json:"bot_token"`
	DisableNotification bool   `json:"disable_notification"`

	ChatID        string
	FromSubreddit string
}

type solution struct {
	Config solutionConfig
	TG     *telegramService
	Cache  *CacheHandler
}

func newSolution() (*solution, error) {
	sol := solution{}

	// parse args
	err := sol.parseArgs()
	if err != nil {
		return nil, err
	}

	// parse config file
	if _, err := os.Stat(configJSONPath); os.IsNotExist(err) {
		return nil, errors.New("failed to find config json")
	}
	jsonBytes, err := ioutil.ReadFile(configJSONPath)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(jsonBytes, &sol.Config)

	// get cache
	sol.Cache, err = NewCacheHandler(cacheFolderPath)
	if err != nil {
		return nil, err
	}

	//log.Println("[DEBUG] chat id: "+sol.Config.ChatID, "bot token: "+sol.Config.BotToken)

	// create tg obj
	sol.TG = newTelegramService().setBotToken(sol.Config.BotToken).
		setChatID(sol.Config.ChatID).setNotificationMode(sol.Config.DisableNotification)

	return &sol, nil
}

func (sol *solution) parseArgs() error {
	fromSubreddit := flag.String("subreddit", "facepalm", "subbredit to crawl posts")
	tgChannelLink := flag.String("channel", "yourchannellink", "telegram channel link to export posts")
	flag.Parse()
	if fromSubreddit == nil {
		return errors.New("failed to get -subreddit arg")
	}
	sol.Config.FromSubreddit = *fromSubreddit
	if tgChannelLink == nil {
		return errors.New("failed to get -channel arg")
	}
	if *tgChannelLink == "" {
		return errors.New("-channel arg is empty")
	}
	sol.Config.ChatID = *tgChannelLink
	return nil
}

func (sol *solution) do() error {
	posts, _, err := reddit.DefaultClient().Subreddit.TopPosts(
		context.Background(), sol.Config.FromSubreddit, &reddit.ListPostOptions{
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

		if postsUsedInQuery == postsPerQuery {
			// all need posts used in this query
			return nil
		}
	}
	return nil
}

func (sol *solution) processPost(post *reddit.Post) bool {
	if sol.Cache.IsPostUsed(sol.Config.ChatID, post.ID) {
		//log.Println("post " + post.ID + " already used")
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
			//log.Println("post " + post.ID + " not contains photo & webpreview")
			return false
		}
		postImageURL = scrapedImages[0]
	}
	if postImageURL == "" {
		log.Println("post " + post.ID + " image is not recognized")
		return false
	}

	err := sol.Cache.MarkPostUsed(sol.Config.ChatID, post.ID)
	if err != nil {
		log.Println("Failed to mark post used: " + err.Error())
		return false
	}

	if !isRemoteFileExists(postImageURL) {
		log.Println("remote image does not exists: " + postImageURL)
		return false
	}

	isDebug := false
	sourceLink := html.A{Value: "[Source]", URL: "https://www.reddit.com" + post.Permalink}
	postText := "<b>" + post.Title + "</b> " + sourceLink.Html()
	if !isDebug {
		err = sol.TG.sendPostImage(postImageURL, postText)
		if err != nil {
			log.Println(postImageURL)
			log.Println("Failed to send photo to channel: " + err.Error())
		}
	} else {
		log.Println("debug, post ID: " + post.ID)
	}
	return true
}

func isPhotoInURL(url string) bool {
	if strings.Contains(url, ".png") {
		return true
	}
	if strings.Contains(url, ".jpg") {
		return true
	}
	return false
}

/*
 _       _
| |     | |
| |_ ___| | ___  __ _ _ __ __ _ _ __ ___
| __/ _ \ |/ _ \/ _` | '__/ _` | '_ ` _ \
| ||  __/ |  __/ (_| | | | (_| | | | | | |
 \__\___|_|\___|\__, |_|  \__,_|_| |_| |_|
                 __/ |
                |___/
*/

type telegramService struct {
	BotToken            string
	ChatID              string
	DisableNotification bool
}

type telegramResponse struct {
	OK          bool                   `json:"ok"`
	Result      telegramResponseResult `json:"result"`
	Description string                 `json:"description"`
}

type telegramResponseResult struct {
	MessageID int64 `json:"message_id"`
	Date      int64 `json:"date"`
}

func newTelegramService() *telegramService {
	return &telegramService{}
}

func (tg *telegramService) setBotToken(token string) *telegramService {
	tg.BotToken = token
	return tg
}

func (tg *telegramService) setChatID(chatID string) *telegramService {
	tg.ChatID = chatID
	return tg
}

func (tg *telegramService) setNotificationMode(enabled bool) *telegramService {
	tg.DisableNotification = !enabled
	return tg
}

/*func (tg *telegramService) sendPost(postText string) error {
	tgAPIURL := "https://api.telegram.org/bot" + tg.BotToken +
		"/sendMessage?chat_id=" + tg.ChatID +
		"&text=" + url.QueryEscape(postText) +
		"&parse_mode=HTML"

	if tg.DisableNotification {
		tgAPIURL += "&disable_notification=true"
	}
	resp, err := http.Get(tgAPIURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	tResp := telegramResponse{}
	parseErr := json.Unmarshal(body, &tResp)
	if parseErr != nil {
		return parseErr
	}
	if !tResp.OK {
		return errors.New(tResp.Description)
	}
	return nil
}*/

func (tg *telegramService) sendPostImage(imageURL, postText string) error {
	tgAPIURL := "https://api.telegram.org/bot" + tg.BotToken +
		"/sendPhoto?chat_id=" + tg.ChatID +
		"&photo=" + url.QueryEscape(imageURL) +
		"&caption=" + url.QueryEscape(postText) +
		"&parse_mode=HTML"

	if tg.DisableNotification {
		tgAPIURL += "&disable_notification=true"
	}
	resp, err := http.Get(tgAPIURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	tResp := telegramResponse{}
	parseErr := json.Unmarshal(body, &tResp)
	if parseErr != nil {
		return parseErr
	}
	if !tResp.OK {
		return errors.New(tResp.Description)
	}
	return nil
}

func isRemoteFileExists(url string) bool {
	resp, err := http.Head(url)
	if err != nil {
		return false
	}
	if resp.StatusCode != http.StatusOK {
		return false
	}
	return true
}
