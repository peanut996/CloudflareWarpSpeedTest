package main

import (
	"flag"
	"fmt"
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
	var help = `
CloudflareWarpSpeedTest ` + Version + `
Test the latency and speed of all Cloudflare Warp IPs to obtain the lowest latency and port.

Parameters:
    -n 200
        Latency test threads; the more threads, the faster the latency test, but do not set it too high on low-performance devices (such as routers); (default 200, maximum 1000)
    -t 10
        Number of latency tests; the number of times to test latency for a single IP; (default 10 times)
    -q
        Quick mode; test results for randomly scanning 1000 addresses; on by default, [-q=false] turns off quick mode
    -ipv6
        IPv6 support. Only effect when not provide extra ip cidr.
    -tl 300
        Average latency upper limit; only output IPs with average latency lower than the specified limit, various upper and lower limit conditions can be used together; (default 300 ms)
    -tll 40
        Average latency lower limit; only output IPs with average latency higher than the specified limit; (default 0 ms)
    -tlr 0.2
        Packet loss rate upper limit; only output IPs with packet loss rate lower than or equal to the specified rate, range 0.00~1.00, 0 filters out any IPs with packet loss; (default 1.00)
    -sl 5
        Download speed lower limit; only output IPs with download speed higher than the specified speed, testing stops when reaching the specified quantity [-dn]; (default 0.00 MB/s)
    -p 10
        Number of results to display; directly display the specified number of results after testing, 0 means not displaying results and exiting directly; (default 10)
    -f ip.txt
        IP segment data file; add quotes if the path contains spaces;
    -ip 1.1.1.1,2.2.2.2/24,2606:4700::/32
        Specify IP segment data; directly specify the IP segment data to be tested through parameters, separated by commas; (default empty)
    -o result.csv
        Write result to file; add quotes if the path contains spaces; empty value means not writing to a file [-o ""]; (default result.csv)
    -pri PrivateKey
        Specify your WireGuard private key
    -pub PublicKey
        Specify your WireGuard public key, default is the Warp public key
    -reserved Reserved
        Add custom reserved field. format: [0, 0, 0]
    -full
        Test all ports; test all ports for each IP in the IP segment
    -h
        Print the help explanation
    -v 
        Print the version
`

	var minDelay, maxDelay int
	var maxLossRate float64
	flag.IntVar(&task.Routines, "n", 200, "Latency test threads")
	flag.IntVar(&task.PingTimes, "t", 10, "Number of latency test times")

	flag.IntVar(&maxDelay, "tl", 300, "Average latency upper limit")
	flag.IntVar(&minDelay, "tll", 0, "Average latency lower limit")
	flag.Float64Var(&maxLossRate, "tlr", 1, "Packet loss rate upper limit")

	flag.BoolVar(&task.ScanAllPort, "full", false, "Scan all ports")
	flag.BoolVar(&task.QuickMode, "q", true, "Quick mode, test results for randomly scanning 1000 IPs")
	flag.BoolVar(&task.IPv6Mode, "ipv6", false, "IPv6 support. Only effect when not provide extra ip cidr.")
	flag.IntVar(&utils.PrintNum, "p", 10, "Number of results to display")
	flag.StringVar(&task.IPFile, "f", "", "IP segment data file")
	flag.StringVar(&task.IPText, "ip", "", "Specify IP segment data")
	flag.StringVar(&utils.Output, "o", "result.csv", "Output result file")
	flag.StringVar(&task.PrivateKey, "pri", "", "Specify private key")
	flag.StringVar(&task.PrivateKey, "pub", "", "Specify public key")
	flag.StringVar(&task.ReservedString, "reserved", "", "Add custom reserved field")
	flag.BoolVar(&printVersion, "v", false, "Print program version")

	flag.Usage = func() { fmt.Print(help) }
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
