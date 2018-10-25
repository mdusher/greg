# Greg

This is Greg.

Greg is a Discord bot written in Go.

Greg isn't really doing much right now. 

Greg is just Greggin'

### Greg works like so
```
git clone git@github.com:mdusher/greg.git
cd greg
docker build -t greg .
docker run -ti --rm -e "BOT_TOKEN=notarealbottokenbecausethatwouldbesilly" -e "BOT_PREFIX=greg" greg
```

### Environment Variables
| Variable          | Description                                                       |
|-------------------|-------------------------------------------------------------------|
| BOT_TOKEN         | The Discord Bot Token                                             |
| BOT_PREFIX        | Comma separated list of prefixes for Greg to react to.            |
| BOT_WEBHOOK_ID    | (optional) Discord Webhook ID for logging to a Discord channel    |
| BOT_WEBHOOK_TOKEN | (optional) Discord Webhook Token for logging to a Discord channel |
