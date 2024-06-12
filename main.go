package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/peanut996/CloudflareWarpSpeedTest/task"
	"github.com/peanut996/CloudflareWarpSpeedTest/utils"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var (
	Version string
)

func init() {
	var printVersion bool
	lang := os.Getenv("LANG")
	var bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.LoadMessageFile("locale/i18n.en.toml")
	bundle.LoadMessageFile("locale/i18n.zh.toml")
	localizer := i18n.NewLocalizer(bundle, strings.Split(lang, ".")[0])

	var minDelay, maxDelay int
	var maxLossRate float64
	flag.IntVar(&task.Routines, "n", 200, localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "n"}))
	flag.IntVar(&task.PingTimes, "t", 10, localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "t"}))
	flag.IntVar(&task.MaxScanCount, "c", 5000, localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "c"}))

	flag.IntVar(&maxDelay, "tl", 300, localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "tl"}))
	flag.IntVar(&minDelay, "tll", 0, localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "tll"}))
	flag.Float64Var(&maxLossRate, "tlr", 1, localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "tlr"}))

	flag.BoolVar(&task.AllMode, "all", false, localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "all"}))
	flag.BoolVar(&task.IPv6Mode, "ipv6", false, localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "ipv6"}))
	flag.IntVar(&utils.PrintNum, "p", 10, localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "p"}))
	flag.StringVar(&task.IPFile, "f", "", localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "f"}))
	flag.StringVar(&task.IPText, "ip", "", localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "ip"}))
	flag.StringVar(&utils.Output, "o", "result.csv", localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "o"}))
	flag.StringVar(&task.PrivateKey, "pri", "", localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "pri"}))
	flag.StringVar(&task.PrivateKey, "pub", "", localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "pub"}))
	flag.StringVar(&task.ReservedString, "reserved", "", localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "reserved"}))
	flag.BoolVar(&printVersion, "v", false, localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "v"}))

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `CloudflareWarpSpeedTest `+`
`+
			Version+localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "h"}))
		flag.PrintDefaults()
	}
	flag.Parse()

	utils.InputMaxDelay = time.Duration(maxDelay) * time.Millisecond
	utils.InputMinDelay = time.Duration(minDelay) * time.Millisecond
	utils.InputMaxLossRate = float32(maxLossRate)

	if printVersion {
		fmt.Println(Version)
		os.Exit(0)
	}
}

func main() {
	task.InitRandSeed()
	task.InitHandshakePacket()

	fmt.Printf("CloudflareWarpSpeedTest\n\n")

	pingData := task.NewWarping().Run().FilterDelay().FilterLossRate()
	utils.ExportCsv(pingData)
	pingData.Print()
}
