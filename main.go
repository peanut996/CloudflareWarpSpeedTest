package main

import (
	"flag"
	"fmt"
	"github.com/peanut996/CloudflareWarpSpeedTest/i18n"
	"os"
	"time"

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
	flag.IntVar(&task.Routines, "n", 200, i18n.QueryI18n("n"))
	flag.IntVar(&task.PingTimes, "t", 10, i18n.QueryI18n("t"))
	flag.IntVar(&task.MaxScanCount, "c", 5000, i18n.QueryI18n("c"))

	flag.IntVar(&maxDelay, "tl", 300, i18n.QueryI18n("tl"))
	flag.IntVar(&minDelay, "tll", 0, i18n.QueryI18n("tll"))
	flag.Float64Var(&maxLossRate, "tlr", 1, i18n.QueryI18n("tlr"))

	flag.BoolVar(&task.AllMode, "all", false, i18n.QueryI18n("all"))
	flag.BoolVar(&task.IPv6Mode, "ipv6", false, i18n.QueryI18n("ipv6"))
	flag.IntVar(&utils.PrintNum, "p", 10, i18n.QueryI18n("p"))
	flag.StringVar(&task.IPFile, "f", "", i18n.QueryI18n("f"))
	flag.StringVar(&task.IPText, "ip", "", i18n.QueryI18n("ip"))
	flag.StringVar(&utils.Output, "o", "result.csv", i18n.QueryI18n("o"))
	flag.StringVar(&task.PrivateKey, "pri", "", i18n.QueryI18n("pri"))
	flag.StringVar(&task.PublicKey, "pub", "", i18n.QueryI18n("pub"))
	flag.StringVar(&task.ReservedString, "reserved", "", i18n.QueryI18n("reserved"))
	flag.BoolVar(&printVersion, "v", false, i18n.QueryI18n("v"))

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `CloudflareWarpSpeedTest `+`
`+Version+i18n.QueryI18n("h"))
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
	task.InitHandshakePacket()

	fmt.Printf("CloudflareWarpSpeedTest\n\n")

	pingData := task.NewWarping().Run().FilterDelay().FilterLossRate()
	utils.ExportCsv(pingData)
	pingData.Print()
}
