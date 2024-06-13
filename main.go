package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	i18n "github.com/peanut996/CloudflareWarpSpeedTest/locale"
	"github.com/peanut996/CloudflareWarpSpeedTest/task"
	"github.com/peanut996/CloudflareWarpSpeedTest/utils"
)

var (
	Version string
)

func init() {
	var printVersion bool
	var minDelay, maxDelay int
	var maxLossRate float64
	localizer := i18n.Init_i18n("main")
	flag.IntVar(&task.Routines, "n", 200, i18n.Query_i18n(localizer, "n"))
	flag.IntVar(&task.PingTimes, "t", 10, i18n.Query_i18n(localizer, "t"))
	flag.IntVar(&task.MaxScanCount, "c", 5000, i18n.Query_i18n(localizer, "c"))

	flag.IntVar(&maxDelay, "tl", 300, i18n.Query_i18n(localizer, "tl"))
	flag.IntVar(&minDelay, "tll", 0, i18n.Query_i18n(localizer, "tll"))
	flag.Float64Var(&maxLossRate, "tlr", 1, i18n.Query_i18n(localizer, "tlr"))

	flag.BoolVar(&task.AllMode, "all", false, i18n.Query_i18n(localizer, "all"))
	flag.BoolVar(&task.IPv6Mode, "ipv6", false, i18n.Query_i18n(localizer, "ipv6"))
	flag.IntVar(&utils.PrintNum, "p", 10, i18n.Query_i18n(localizer, "p"))
	flag.StringVar(&task.IPFile, "f", "", i18n.Query_i18n(localizer, "f"))
	flag.StringVar(&task.IPText, "ip", "", i18n.Query_i18n(localizer, "ip"))
	flag.StringVar(&utils.Output, "o", "result.csv", i18n.Query_i18n(localizer, "o"))
	flag.StringVar(&task.PrivateKey, "pri", "", i18n.Query_i18n(localizer, "pri"))
	flag.StringVar(&task.PrivateKey, "pub", "", i18n.Query_i18n(localizer, "pub"))
	flag.StringVar(&task.ReservedString, "reserved", "", i18n.Query_i18n(localizer, "reserved"))
	flag.BoolVar(&printVersion, "v", false, i18n.Query_i18n(localizer, "v"))

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `CloudflareWarpSpeedTest `+`
`+
			Version+i18n.Query_i18n(localizer, "h"))
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
