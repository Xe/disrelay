# disrelay

`disrelay` is a dead simple IRC<->Discord relay bot. It takes in IRC messages and
batches them for discord and then relays every individual line in a Discord message
as an IRC line.

## Configuration

This application follows the [12 factor rules of configuration](https://12factor.net/config). 
All configuration for disrelay is done in the form of environment variables. 
For convenience, a helper library called `godotenv` will automagically load all
of the shell-style `KEY=value` pairs defined in `.env` of the current working
directory when the program starts if it exists. Usage is simple: create a file 
named `.env` that contains the following key->value pairs:

```shell
DISCORD_TOKEN=<discord token here>
CHANNEL_MAP=<discord channel id>:<irc channel name without hash>,269624426272653312:geek
IRC_NICK=Ryleth
IRC_USER=flagatn
IRC_SERVER=irc.ponychat.net:6697
IRC_PASS=hunter2
IRC_TLS=true
```

Below each variable here gets explained in detail:

### `DISCORD_TOKEN`

Specifies the Discord bot token that `disrelay` should use when authenticating
with the Discord bot API.

This setting is required.

### `CHANNEL_MAP`

Specifies the mapping of Discord channel IDs to irc channel names WITHOUT the 
leading hash. An example is:

```
CHANNEL_MAP=269624426272653312:geek
```

This setting is required.

### `IRC_NICK`

Specifies the nickname that will be used when `disrelay` connects to an IRC server.

This setting is required.

### `IRC_USER`

Specifies the username that will be used when `disrelay` connects to an IRC server.

If unset, the default IRC username will be `ad`.

### `IRC_SERVER`

Specifies the irc server (host+port) that `disrelay` will connect to.

If unset, the default IRC server will be `127.0.0.1:6667`.

### `IRC_PASS`

Specifies the server password (`PASS`, [see this section of the RFC](https://tools.ietf.org/html/rfc1459#section-4.1.1)) that will be used when `disrelay` connects to an IRC server.

This setting is not required, not setting it is fine.

### `IRC_TLS`

If set to `true`, `disrelay` will use TLS connections to connect to IRC servers.

If unset, the default value is `false` for convenience of local testing.

## Running

### Locally

Download the binaries for your favorite OS and put your `.env` file somewhere,
then run `disrelay` in a restart loop such as:

```shell
# ~/prefix/disrelay/run: mode 744

while true
do
  ./disrelay
  sleep 2 # just to be sure
done
```

### Docker

```shell
$ docker run -d --name disrelay --env-file .env xena/disrelay:1.2
```

### Docker Swarm

```shell
$ docker service create --name disrelay --env-file .env xena/disrelay:1.2
```

## License

```
Copyright (c) 2017 Christine Dodrill <me@christine.website>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
```

## Libraries Used

- [ln](https://github.com/Xe/ln) for logging (`key=value` pairs)
- [discordgo](https://github.com/bwmarrin/discordgo) for Discord bot functions
- [env](https://github.com/caarlos0/env) for parsing the configuration from the environment
- [godotenv](https://github.com/joho/godotenv) for reading envvars out of `.env`
- [irc](https://godoc.org/gopkg.in/irc.v1) for IRC bot functions
- [bundler](https://godoc.org/google.golang.org/api/support/bundler) for batching up IRC messages to send to discord in one big burst

## Support

To get help with this bot, please message me on Discord `Cadey~#1932` or [Telegram](https://t.me/miamorecadenza). I'm happy to help with installation, add/remove features and whatever else is needed to make the goal of message relaying from IRC to Discord a non-issue.
