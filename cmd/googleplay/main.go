package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/89z/format"
	gp "github.com/JoKr4/googleplay/pkg/googleplay"
)

func main() {
	// a
	var app string
	flag.StringVar(&app, "a", "", "app")
	// d
	dir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	dir = filepath.Join(dir, "googleplay")
	flag.StringVar(&dir, "d", dir, "user dir")
	// date
	var parse bool
	flag.BoolVar(&parse, "date", false, "parse date")
	// device
	var device bool
	flag.BoolVar(&device, "device", false, "create device")
	// dpi
	var screenDensity int
	flag.IntVar(&screenDensity, "dpi", 320, "screen density of device")
	// email
	var email string
	flag.StringVar(&email, "email", "", "your email")
	// log
	var level int
	flag.IntVar(&level, "log", 0, "log level")
	// p
	var platformID int64
	flag.Int64Var(&platformID, "p", 0, gp.Platforms.String())
	// password
	var password string
	flag.StringVar(&password, "password", "", "your password")
	// purchase
	var (
		buf      strings.Builder
		purchase bool
	)
	buf.WriteString("Purchase app. ")
	buf.WriteString("Only needs to be done once per Google account.")
	flag.BoolVar(&purchase, "purchase", false, buf.String())
	// s
	var single bool
	flag.BoolVar(&single, "s", false, "single APK")
	// v
	var version uint64
	flag.Uint64Var(&version, "v", 0, "app version")
	flag.Parse()
	gp.LogLevel = format.LogLevel(level)
	if email != "" {
		err := gp.DoToken(dir, email, password)
		if err != nil {
			panic(err)
		}
	} else {
		platform := gp.Platforms[platformID]
		if device {
			err := gp.DoDevice(dir, platform, screenDensity)
			if err != nil {
				panic(err)
			}
		} else if app != "" {
			head, err := gp.GetHeader(dir, platform, single)
			if err != nil {
				panic(err)
			}
			if purchase {
				err := head.Purchase(app)
				if err != nil {
					panic(err)
				}
			} else if version >= 1 {
				err := gp.DoDelivery(head, app, version, dir)
				if err != nil {
					panic(err)
				}
			} else {
				details, err := gp.GetDetails(head, app, parse)
				if err != nil {
					panic(err)
				}
				fmt.Print(details)
			}
		} else {
			flag.Usage()
		}
	}
}
