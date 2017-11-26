package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/Xe/ln"
	"github.com/bwmarrin/discordgo"
	"github.com/caarlos0/env"
	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/api/support/bundler"
	irc "gopkg.in/irc.v1"
)

var cfg struct {
	DiscordToken string   `env:"DISCORD_TOKEN,required"`
	ChannelMap   []string `env:"CHANNEL_MAP"`
	IRCNick      string   `env:"IRC_NICK,required"`
	IRCUser      string   `env:"IRC_USER" envDefault:"ad"`
	IRCServer    string   `env:"IRC_SERVER" envDefault:"127.0.0.1:6667"`
	IRCPass      string   `env:"IRC_PASS"`
	IRCTLS       bool     `env:"IRC_TLS" envDefault:"false"`
}

func main() {
	ctx := context.Background()

	err := env.Parse(&cfg)
	if err != nil {
		ln.FatalErr(ctx, err, ln.Action("parse config from env"))
	}

	cm := map[string]string{}
	rcm := map[string]string{}

	for _, val := range cfg.ChannelMap {
		sp := strings.Split(val, ":")
		if len(sp) != 2 {
			ln.Fatal(ctx, ln.Action("bad channel map entry, wanted discord:irc"), ln.F{"val": sp})
		}

		ln.Log(ctx, ln.Action("associated channel"), ln.F{"discord": sp[0], "irc": "#" + sp[1]})

		cm[sp[0]] = sp[1]
		rcm[sp[1]] = sp[0]
	}

	dg, err := discordgo.New("Bot " + cfg.DiscordToken)
	if err != nil {
		ln.FatalErr(ctx, err, ln.Action("discordgo create"))
	}
	ln.Log(ctx, ln.Action("discordgo client created"))

	b := &bot{
		dischans:          map[string]*bundler.Bundler{},
		channelMap:        cm,
		reverseChannelMap: rcm,
		s:                 dg,
	}

	for did := range b.channelMap {
		bd := bundler.NewBundler("", b.bundlerHandler(did))
		bd.DelayThreshold = 250 * time.Millisecond
		bd.BundleByteLimit = 1750
		bd.BundleByteThreshold = 1700
		bd.BufferedByteLimit = 4096
		b.dischans[did] = bd
	}

	dg.AddHandler(b.handleDiscord)

	ln.Log(ctx, ln.Action("dialing IRC"), ln.F{"server": cfg.IRCServer, "tls": cfg.IRCTLS})
	var conn net.Conn

	if cfg.IRCTLS {
		conn, err = tls.Dial("tcp", cfg.IRCServer, &tls.Config{})
		if err != nil {
			ln.FatalErr(ctx, err, ln.F{"server": cfg.IRCServer, "tls": cfg.IRCTLS})
		}
	} else {
		conn, err = net.Dial("tcp", cfg.IRCServer)
		if err != nil {
			ln.FatalErr(ctx, err, ln.F{"server": cfg.IRCServer, "tls": cfg.IRCTLS})
		}
	}

	b.i = irc.NewClient(conn, irc.ClientConfig{
		Nick:    cfg.IRCNick,
		Pass:    cfg.IRCPass,
		User:    cfg.IRCUser,
		Name:    "disrelay: the relay bot that doesn't make you want to yank your eyes out",
		Handler: b,
	})
	ln.Log(ctx, ln.Action("irc client created"))

	err = dg.Open()
	if err != nil {
		ln.FatalErr(ctx, err)
	}

	b.i.Run()
}

type bot struct {
	dischans          map[string]*bundler.Bundler
	channelMap        map[string]string // discord -> irc
	reverseChannelMap map[string]string // irc -> discord

	s     *discordgo.Session
	i     *irc.Client
	ilock sync.Mutex
}

func (b *bot) handleDiscord(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	ctx := context.Background()
	msg, err := m.ContentWithMoreMentionsReplaced(s)
	if err != nil {
		ln.Error(ctx, err)
		return
	}

	ic, ok := b.channelMap[m.ChannelID]
	if !ok {
		return
	}

	err = b.messageIRC(ic, m.Author.Username, msg)
	if err != nil {
		ln.FatalErr(ctx, err, ln.F{"discord_channel": m.ChannelID, "discord_username": m.Author.Username, "msg": msg})
	}
}

func (b *bot) bundlerHandler(channel string) func(interface{}) {
	return func(msgsi interface{}) {
		msgs, ok := msgsi.([]string)
		if !ok {
			return
		}

		msg := strings.Join(msgs, "\n")
		_, err := b.s.ChannelMessageSend(channel, msg)
		if err != nil {
			ln.Error(context.Background(), err, ln.Action("send bundled message"))
		}
	}
}

func (b *bot) messageDiscord(channel, who, what string) error {
	msg := fmt.Sprintf("<%s> %s", who, what)

	bd, ok := b.dischans[channel]
	if !ok {
		return errors.New("unknown channel")
	}

	bd.Add(msg, len(msg)+4)

	return nil
}

func (b *bot) messageIRC(channel, who, what string) error {
	b.ilock.Lock()
	defer b.ilock.Unlock()
	return b.i.Writef("PRIVMSG #%s :<%s> %s", channel, who, what)
}

func (b *bot) Handle(c *irc.Client, m *irc.Message) {
	if m.Prefix.Name == c.CurrentNick() {
		return
	}

	ctx := context.Background()

	ln.Log(ctx, ln.F{"message": m.String()})

	switch m.Command {
	case "001":
		time.Sleep(time.Second)

		cl := []string{}

		for ch, _ := range b.reverseChannelMap {
			cl = append(cl, "#"+ch)
		}

		msg := fmt.Sprintf("JOIN %s", strings.Join(cl, ","))
		c.Write(msg)
		ln.Log(ctx, ln.F{"line": msg})

	case "PRIVMSG":
		ch := m.Params[0][1:]
		nick := m.Prefix.Name
		msg := m.Params[1]

		channel, ok := b.reverseChannelMap[ch]
		if !ok {
			return
		}

		err := b.messageDiscord(channel, nick, msg)
		if err != nil {
			ln.FatalErr(ctx, err)
		}
	}
}
