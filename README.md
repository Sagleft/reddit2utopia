![logo](https://github.com/Sagleft/reddit2utopia/raw/main/assets/logo.png)

Bot for throwing posts from Reddit to Utopia. [Tutorial](https://info.crypton.click/articles/20-create-an-autofill-channel-in-utopia.html).

## concept

This solution is useful for you if you want to create a channel in Utopia, but do not want to manually fill it with content. This bot will do everything for you.

You choose the subreddit and the channel in Utopia to transfer these posts to. Further along the crown, call this bot to transfer 1 post from the last posts in 24 hours.
The bot processes only those posts to which a link or image is attached.

The finished build can be found on the [releases page.](releases)

## run without docker

TODO

## run in docker

You'll also need this: [utopia-api-docker](https://github.com/Sagleft/utopia-api-docker)

create new container:

```bash
docker create --name redditbot uto9234/reddit2utopia
```

## setup

File `config.json` contains the bot settings, fill in the data to connect to Utopia client.

The account to which you will connect via the API must be a member of the chat \ channel to which you want to send messages, and also have the rights to write messages.

## bot cross-platform build

just run:

```bash
bash make.sh
```

or build for current system:

```bash
go build
```

To do this, you must have Golang v1.16.3 + installed on your system

## usage example

In this example, the post from subreddit "anonim" will be taken and placed in channel in Utopia

```bash
./bot_linux-amd64 -subreddit=anonim -channel=16288010C39BD8D20C337FFC9684657F
```

## help me stay productive

<a href="https://www.buymeacoffee.com/sagleft" target="_blank"><img src="https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png" alt="Buy Me A Coffee" style="height: 41px !important;width: 174px !important;box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;-webkit-box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;" ></a>

* Crypton: F50AF5410B1F3F4297043F0E046F205BCBAA76BEC70E936EB0F3AB94BF316804

---

More useful information about Utopia on the website: [info.crypton.click](https://info.crypton.click)
