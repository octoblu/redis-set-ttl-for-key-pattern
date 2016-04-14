package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/codegangsta/cli"
	"github.com/coreos/go-semver/semver"
	"github.com/fatih/color"
	"github.com/octoblu/redis-set-ttl-for-key-pattern/cleaner"
	De "github.com/tj/go-debug"
)

var debug = De.Debug("redis-set-ttl-for-key-pattern:main")

func main() {
	app := cli.NewApp()
	app.Name = "redis-set-ttl-for-key-pattern"
	app.Version = version()
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "pattern, p",
			EnvVar: "REDIS_SET_TTL_FOR_KEY_PATTERN_PATTERN",
			Usage:  "Pattern for keys to expire. (ex: '*job*')",
		},
		cli.StringFlag{
			Name:   "redis-uri, r",
			EnvVar: "REDIS_SET_TTL_FOR_KEY_PATTERN_REDIS_URI",
			Usage:  "Fully qualified URI where redis is",
		},
	}
	app.Run(os.Args)
}

func run(context *cli.Context) {
	pattern, redisURI := getOpts(context)

	sigTerm := make(chan os.Signal)
	signal.Notify(sigTerm, syscall.SIGTERM)

	sigTermReceived := false

	go func() {
		<-sigTerm
		fmt.Println("SIGTERM received, waiting to exit")
		sigTermReceived = true
	}()

	cursor := 0
	client, err := cleaner.New(pattern, redisURI)
	panicIfError("Failed to construct cleaner: ", err)
	defer client.Close()

	for {
		if sigTermReceived {
			fmt.Println("I'll be back.")
			os.Exit(0)
		}

		fmt.Println("Cleaning: ", cursor)
		var err error
		cursor, err = client.Clean(cursor)
		panicIfError("Error on client.clean", err)
		if cursor == 0 {
			fmt.Println("All done.")
			os.Exit(0)
		}
	}
}

func getOpts(context *cli.Context) (string, string) {
	pattern := context.String("pattern")
	redisURI := context.String("redis-uri")

	if redisURI == "" {
		cli.ShowAppHelp(context)

		if pattern == "" {
			color.Red("  Missing required flag --pattern or REDIS_SET_TTL_FOR_KEY_PATTERN_PATTERN")
		}
		if redisURI == "" {
			color.Red("  Missing required flag --redis-uri or REDIS_SET_TTL_FOR_KEY_PATTERN_REDIS_URI")
		}
		os.Exit(1)
	}

	return pattern, redisURI
}

func panicIfError(msg string, err error) {
	if err == nil {
		return
	}

	log.Panicln(msg, err.Error())
}

func version() string {
	version, err := semver.NewVersion(VERSION)
	panicIfError(fmt.Sprintf("Error with version number: %v", VERSION), err)

	return version.String()
}
