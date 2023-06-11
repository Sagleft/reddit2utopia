![logo](logo.jpg)

## How it works

This solution is useful for you if you want to create a channel in Utopia, but do not want to manually fill it with content. This bot will do everything for you.

You choose the subreddit and the channel in Utopia to transfer these posts to. Further along the crown, call this bot to transfer 1 post from the last posts in 24 hours.
The bot processes only those posts to which a link or image is attached.

The finished build can be found on the [releases page.](releases)

## Get started

1. Simply change the parameters in `docker-compose.yml` file.
2. Put `account.db` Utopia account file in app directory.
3. Run:

```bash
docker-compose up -d
```

:information_source: view bot logs:

```bash
docker container logs reddit2utopia_bot_1
```

:information_source: view Utopia client logs:

```bash
docker container logs reddit2utopia_utopia-api_1
```

:information_source: If there is no container by this name, then you can find it through:

```bash
docker ps -a | grep reddit2utopia
```

:information_source: stop app:

```bash
docker-compose down
```

## useful links

```
TODO
```
