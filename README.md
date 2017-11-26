# disrelay

## Configuration

Create a file named `.env` that contains the following key->value pairs:

```shell
DISCORD_TOKEN=<discord token here>
CHANNEL_MAP=<discord channel id>:<irc channel name without hash>,269624426272653312:geek
IRC_NICK=Ryleth
IRC_USER=flagatn
IRC_SERVER=irc.ponychat.net:6697
IRC_PASS=hunter2
IRC_TLS=true
```

## Running

```shell
$ docker service create --name disrelay --env-file .env xena/disrelay:1.2
```
