package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/peanut996/CloudflareWarpSpeedTest/i18n"

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
	flag.IntVar(&task.Routines, "n", 200, i18n.QueryI18n(i18n.TestThreadCount))
	flag.IntVar(&task.PingTimes, "t", 10, i18n.QueryI18n(i18n.LatencyTestTimes))
	flag.IntVar(&task.MaxScanCount, "c", 5000, i18n.QueryI18n(i18n.ScanAddressCount))

	flag.IntVar(&maxDelay, "tl", 300, i18n.QueryI18n(i18n.LatencyUpperLimit))
	flag.IntVar(&minDelay, "tll", 0, i18n.QueryI18n(i18n.LatencyLowerLimit))
	flag.Float64Var(&maxLossRate, "tlr", 1, i18n.QueryI18n(i18n.PacketLossRateUpperLimit))

	flag.BoolVar(&task.AllMode, "all", false, i18n.QueryI18n(i18n.TestAllIpPortCombinations))
	flag.BoolVar(&task.IPv6Mode, "ipv6", false, i18n.QueryI18n(i18n.ScanIpv6Only))
	flag.IntVar(&utils.PrintNum, "p", 10, i18n.QueryI18n(i18n.ResultDisplayCount))
	flag.StringVar(&task.IPFile, "f", "", i18n.QueryI18n(i18n.IpDataFile))
	flag.StringVar(&task.IPText, "ip", "", i18n.QueryI18n(i18n.SpecifyIpData))
	flag.StringVar(&utils.Output, "o", "result.csv", i18n.QueryI18n(i18n.OutputResultFile))
	flag.StringVar(&task.PrivateKey, "pri", "", i18n.QueryI18n(i18n.CustomWireguardPrivateKey))
	flag.StringVar(&task.PublicKey, "pub", "", i18n.QueryI18n(i18n.CustomWireguardPublicKey))
	flag.StringVar(&task.ReservedString, "reserved", "", i18n.QueryI18n(i18n.CustomReservedField))
	flag.BoolVar(&printVersion, "v", false, i18n.QueryI18n(i18n.ProgramVersion))

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `CloudflareWarpSpeedTest `+"\n\n"+i18n.QueryI18n(i18n.HelpMessage))
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
