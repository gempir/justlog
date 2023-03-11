# justlog [![Build Status](https://github.com/gempir/justlog/workflows/ci/badge.svg)](https://github.com/gempir/justlog/actions?query=workflow%3Aci)

### What is this?
Justlog is a twitch irc bot. It focuses on logging and providing an api for the logs.

### Optout

Click the X icon on the web ui to find an explanation for how to opt out.

### API

API documentation can be viewed via the justlog frontend by clicking the "docs" symbol:  
![image](https://user-images.githubusercontent.com/1629196/159481078-0de98f01-2816-49bd-8e17-ba7cf66cb064.png)

### Docker

```
mkdir logs
docker run -p 8025:8025 --restart=unless-stopped -v $PWD/config.json:/etc/justlog.json -v $PWD/logs:/logs ghcr.io/gempir/justlog
```

### Commands

Only admins can use these commands

- `!justlog status` will respond with uptime.
- `!justlog join gempir pajlada` will join the channels and append them to the config.
- `!justlog part gempir pajlada` will part the channels and remove them from the config.
- `!justlog optout gempir gempbot` will opt out users of message logging or querying previous logs of that user. the same applies to a user's own channel.
- `!justlog optin gempir gempbot` will revert the opt out.

### Config

```
{
    "admins": ["gempir"], // will only respond to commands executed by these users
    "logsDirectory": "./logs", // the directory to log into
    "adminAPIKey": "noshot", // your secret api key to access the admin api, can be any string, used in api request to admin endpoints
    "username": "gempbot", // bot username (can be justinfan123123 if you don't want to use an account)
    "oauth": "oauthtokenforchat", // bot token can be anything if justinfan123123
    "botVerified": true, // increase ratelimits if you have a verified bot, so the bot can join faster, false by default
    "clientID": "mytwitchclientid", // your client ID, needed for fetching userids or usernames etc
    "clientSecret": "mysecret", // your twitch client secret
    "logLevel": "info", // the log level, keep this to info probably, all options are: trace, debug, info, warn, error, fatal, and panic, logs output is stdout
    "channels": ["77829817", "11148817"], // the channels (userids) you want to log
    "archive": true // probably keep to true, will disable gzipping of old logs if false, useful if you setup compression on your own
}
```

### Development

Development requires [yarn](https://classic.yarnpkg.com/) and [go-swagger](https://goswagger.io/)

Run `go build && ./justlog` and `yarn start` in the web folder. 

Or run `make container` and `make run_container`
