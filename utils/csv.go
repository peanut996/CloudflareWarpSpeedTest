package utils

import (
	"encoding/csv"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/peanut996/CloudflareWarpSpeedTest/i18n"
)

const (
	defaultOutput         = "warp.csv"
	maxDelay              = 9999 * time.Millisecond
	minDelay              = 0 * time.Millisecond
	maxLossRate   float32 = 1.0
)

var (
	InputMaxDelay    = maxDelay
	InputMinDelay    = minDelay
	InputMaxLossRate = maxLossRate
	Output           = defaultOutput
	PrintNum         = 10
)

func NoPrintResult() bool {
	return PrintNum == 0
}

func noOutput() bool {
	return Output == "" || Output == " "
}

type PingData struct {
	IP       *net.UDPAddr
	Sent     int
	Received int
	Delay    time.Duration
}

type CloudflareIPData struct {
	*PingData
	lossRate float32
}

func (cf *CloudflareIPData) getLossRate() float32 {
	if cf.lossRate == 0 {
		pingLost := cf.Sent - cf.Received
		cf.lossRate = float32(pingLost) / float32(cf.Sent)
	}
	return cf.lossRate
}

func (cf *CloudflareIPData) toString() []string {
	result := make([]string, 3)
	result[0] = cf.IP.String()
	result[1] = strconv.FormatFloat(float64(cf.getLossRate())*100, 'f', 0, 32) + "%"
	result[2] = strconv.FormatFloat(cf.Delay.Seconds()*1000, 'f', 2, 32)
	return result
}

func ExportCsv(data []CloudflareIPData) {
	if noOutput() || len(data) == 0 {
		return
	}
	fp, err := os.Create(Output)
	if err != nil {
		log.Fatalf(i18n.QueryTemplateI18n(i18n.CreateFileFailed, map[string]interface{}{"Output": Output, "err": err}))
		return
	}
	defer fp.Close()
	w := csv.NewWriter(fp)
	_ = w.Write([]string{"IP:Port", "Loss", "Latency"})
	_ = w.WriteAll(convertToString(data))
	w.Flush()
}

func convertToString(data []CloudflareIPData) [][]string {
	result := make([][]string, 0)
	for _, v := range data {
		result = append(result, v.toString())
	}
	return result
}

type PingDelaySet []CloudflareIPData

func (s PingDelaySet) FilterDelay() (data PingDelaySet) {
	if InputMaxDelay > maxDelay || InputMinDelay < minDelay {
		return s
	}
	if InputMaxDelay == maxDelay && InputMinDelay == minDelay {
		return s
	}
	for _, v := range s {
		if v.Delay > InputMaxDelay {
			break
		}
		if v.Delay < InputMinDelay {
			continue
		}
		data = append(data, v)
	}
	return
}

func (s PingDelaySet) FilterLossRate() (data PingDelaySet) {
	if InputMaxLossRate >= maxLossRate {
		return s
	}
	for _, v := range s {
		if v.getLossRate() > InputMaxLossRate {
			break
		}
		data = append(data, v)
	}
	return
}

func (s PingDelaySet) Len() int {
	return len(s)
}
func (s PingDelaySet) Less(i, j int) bool {
	iRate, jRate := s[i].getLossRate(), s[j].getLossRate()
	if iRate != jRate {
		return iRate < jRate
	}
	return s[i].Delay < s[j].Delay
}
func (s PingDelaySet) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s PingDelaySet) Print() {
	if NoPrintResult() {
		return
	}
	if len(s) <= 0 {
		fmt.Println(i18n.QueryI18n(i18n.TotalResultZeroSkipOutput))
		return
	}
	dataString := convertToString(s)
	if len(dataString) < PrintNum {
		PrintNum = len(dataString)
	}
	headFormat := "\n%-24s%-9s%-10s\n"
	dataFormat := "%-25s%-8s%-10s\n"
	for i := 0; i < PrintNum; i++ {
		if len(dataString[i][0]) > 15 {
			headFormat = "\n%-44s%-9s%-10s\n"
			dataFormat = "%-45s%-8s%-10s\n"
		}
	}
	fmt.Printf(headFormat, "IP:Port", i18n.QueryI18n(i18n.PacketLossRate), i18n.QueryI18n(i18n.Latency))
	for i := 0; i < PrintNum; i++ {
		fmt.Printf(dataFormat, dataString[i][0], dataString[i][1], dataString[i][2])
	}
	if !noOutput() {
		fmt.Println(i18n.QueryTemplateI18n(i18n.WriteResultToFileDone, map[string]interface{}{"Output": Output}))
	}
}
