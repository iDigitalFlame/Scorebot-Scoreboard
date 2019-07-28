package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/iDigitalFlame/scorebot-scoreboard/scoreboard"
	"github.com/iDigitalFlame/scorebot-scoreboard/scoreboard/web"
)

const (
	version = "v1.3"
)

func main() {
	ConfigFile := flag.String("c", "", "Scorebot Config File Path.")
	ConfigDefault := flag.Bool("d", false, "Print Default Config and Exit.")

	LogFile := flag.String("log", "", "Scoreboard Log File Path.")
	LogLevel := flag.Int("log-level", int(scoreboard.DefaultLogLevel), "Scoreboard Log Level.")

	Tick := flag.Int("tick", int(scoreboard.DefaultTick), "Scoreboard Poll Rate. (in seconds)")

	Listen := flag.String("bind", scoreboard.DefaultListen, "Address and Port to Listen on.")

	Timeout := flag.Int("timeout", int(scoreboard.DefaultTimeout), "Scoreboard Request Timeout. (in seconds)")

	Scorebot := flag.String("sbe", "", "Scorebot Core Address or URL.")

	Directory := flag.String("dir", "", "Scoreboard HTML Directory Path.")

	TwitterCosumerKey := flag.String("tw-ck", "", "Twitter Consumer API Key.")
	TwitterCosumerSecret := flag.String("tw-cs", "", "Twitter Consumer API Secret.")

	TwitterAccessKey := flag.String("tw-ak", "", "Twitter Access API Key.")
	TwitterAccessSecret := flag.String("tw-as", "", "Twitter Access API Secret.")

	TwitterKeywords := flag.String("tw-keywords", "", "Twitter Search Keywords. (comma seperated)")
	TwitterLanguage := flag.String("tw-lang", "", "Twitter Search Lanugage. (comma seperated)")

	TwitterExpire := flag.Int("tw-expire", int(scoreboard.DefaultExpire), "Tweet Display Time. (in seconds)")

	TwitterBlockWords := flag.String("tw-block-words", "", "Twitter Blocked Words. (comma seperated)")
	TwitterBlockUsers := flag.String("tw-block-user", "", "Twitter Blocked Usernames. (comma seperated)")

	TwitterOnlyUsers := flag.String("tw-only-users", "", "Twitter WHitelisted Usernames. (comma seperated)")

	flag.Usage = func() {
		fmt.Printf("Scorebot Scoreboard %s\n\nUsage:\n", version)
		flag.PrintDefaults()
	}
	flag.Parse()

	if *ConfigDefault {
		fmt.Printf("%s\n", scoreboard.Defaults())
		os.Exit(0)
	}

	var c *scoreboard.Config
	if len(*ConfigFile) > 0 {
		var err error
		c, err = scoreboard.Load(*ConfigFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(1)
		}
	} else {
		if len(*Scorebot) == 0 || len(*Listen) == 0 || *Tick <= 0 || *Timeout < 0 || *TwitterExpire <= 0 {
			flag.Usage()
			os.Exit(2)
		}
		c = &scoreboard.Config{
			Log: &scoreboard.Log{
				File:  *LogFile,
				Level: uint8(*LogLevel),
			},
			Tick:   uint16(*Tick),
			Listen: *Listen,
			Twitter: &scoreboard.Twitter{
				Filter: &web.Filter{
					Language:     scoreboard.SplitParm(*TwitterLanguage, scoreboard.ConfigSeperator),
					Keywords:     scoreboard.SplitParm(*TwitterKeywords, scoreboard.ConfigSeperator),
					OnlyUsers:    scoreboard.SplitParm(*TwitterOnlyUsers, scoreboard.ConfigSeperator),
					BlockedUsers: scoreboard.SplitParm(*TwitterBlockUsers, scoreboard.ConfigSeperator),
					BlockedWords: scoreboard.SplitParm(*TwitterBlockWords, scoreboard.ConfigSeperator),
				},
				Expire:  uint16(*TwitterExpire),
				Timeout: uint16(*Timeout),
				Credentials: &web.Credentials{
					AccessKey:      *TwitterAccessKey,
					ConsumerKey:    *TwitterCosumerKey,
					AccessSecret:   *TwitterAccessSecret,
					ConsumerSecret: *TwitterCosumerSecret,
				},
			},
			Timeout:   uint16(*Timeout),
			Scorebot:  *Scorebot,
			Directory: *Directory,
		}
	}

	scoreboard, err := scoreboard.NewScoreboard(c)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}

	if err := scoreboard.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}
