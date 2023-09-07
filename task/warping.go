package task

import (
	"CloudflareWarpSpeedTest/utils"
	"encoding/hex"
	"fmt"
	"net"
	"sort"
	"strconv"
	"sync"
	"time"
)

const (
	defaultRoutines   = 200
	defaultPingTimes  = 5
	udpConnectTimeout = time.Millisecond * 300

	warpValidatePacket = "cf000000628748824150e38f5c64b477"
)

var (
	ScanAllPort = false

	Routines = defaultRoutines

	PingTimes int = defaultPingTimes

	commonWarpPorts = []int{
		2408, 500, 1701, 4500, 854, 859, 864, 878, 880, 890, 891, 894, 903, 908, 928, 934, 939, 942, 943, 945, 946, 955, 968, 987, 988, 1002, 1010, 1014, 1018, 1070, 1074, 1180, 1387, 1843, 2371, 2506,
	}

	commonWarpCIDRs = []string{
		"162.159.192.0/24",
		"162.159.193.0/24",
		"162.159.195.0/24",
	}

	MaxWarpPortRange = 10000

	warpHandshakePacket, _ = hex.DecodeString("04e77a11628748824150e38f5c64b4776d82d118ed6ee00d8ede7ae82405df0c380000000000000000000000004154e7e7b6bbbb84ab8cd5e9b0f82a1c")
)

type UDPAddr struct {
	IP   *net.IPAddr
	Port int
}

type Warping struct {
	wg      *sync.WaitGroup
	m       *sync.Mutex
	ips     []*UDPAddr
	csv     utils.PingDelaySet
	control chan bool
	bar     *utils.Bar
}

func NewWarping() *Warping {
	checkPingDefault()
	ips := loadWarpIPRanges()
	return &Warping{
		wg:      &sync.WaitGroup{},
		m:       &sync.Mutex{},
		ips:     ips,
		csv:     make(utils.PingDelaySet, 0),
		control: make(chan bool, Routines),
		bar:     utils.NewBar(len(ips), "可用:", ""),
	}
}

func checkPingDefault() {
	if Routines <= 0 {
		Routines = defaultRoutines
	}
	if PingTimes <= 0 {
		PingTimes = defaultPingTimes
	}
}

func (w *Warping) Run() utils.PingDelaySet {
	if len(w.ips) == 0 {
		return w.csv
	}
	for _, ip := range w.ips {
		w.wg.Add(1)
		w.control <- false
		go w.start(ip)
	}
	w.wg.Wait()
	w.bar.Done()
	sort.Sort(w.csv)
	return w.csv
}

func (w *Warping) start(ip *UDPAddr) {
	defer w.wg.Done()
	w.warpingHandler(ip)
	<-w.control
}

func (w *Warping) checkConnection(ip *UDPAddr) (recv int, totalDelay time.Duration) {
	for i := 0; i < PingTimes; i++ {
		if ok, delay := w.warping(ip); ok {
			recv++
			totalDelay += delay
		}
	}
	return
}

func (w *Warping) warpingHandler(ip *UDPAddr) {
	recv, totalDelay := w.checkConnection(ip)
	nowAble := len(w.csv)
	if recv != 0 {
		nowAble++
	}
	w.bar.Grow(1, strconv.Itoa(nowAble))
	if recv == 0 {
		return
	}
	data := &utils.PingData{
		IP:       ip.ToUDPAddr(),
		Sended:   PingTimes,
		Received: recv,
		Delay:    totalDelay / time.Duration(recv),
	}
	w.appendIPData(data)
}

func (w *Warping) appendIPData(data *utils.PingData) {
	w.m.Lock()
	defer w.m.Unlock()
	w.csv = append(w.csv, utils.CloudflareIPData{
		PingData: data,
	})
}

func loadWarpIPRanges() (ipAddrs []*UDPAddr) {
	ips := loadIPRanges()
	return generateIPAddrs(ips)
}

func generateIPAddrs(ips []*net.IPAddr) (udpAddrs []*UDPAddr) {
	if !ScanAllPort {
		for _, port := range commonWarpPorts {
			udpAddrs = append(udpAddrs, generateSingleIPAddr(ips, port)...)
		}
		return
	}
	for port := 1; port <= MaxWarpPortRange; port++ {
		udpAddrs = append(udpAddrs, generateSingleIPAddr(ips, port)...)
	}
	return udpAddrs
}

func generateSingleIPAddr(ips []*net.IPAddr, port int) []*UDPAddr {
	udpAddrs := make([]*UDPAddr, 0)
	for _, ip := range ips {
		udpAddrs = append(udpAddrs, &UDPAddr{
			IP:   ip,
			Port: port,
		})
	}
	return udpAddrs
}

func (i *UDPAddr) FullAddress() string {
	if isIPv4(i.IP.String()) {
		return fmt.Sprintf("%s:%d", i.IP.String(), i.Port)
	}
	return fmt.Sprintf("[%s]:%d", i.IP.String(), i.Port)

}

func (i *UDPAddr) ToUDPAddr() (addr *net.UDPAddr) {
	addr, _ = net.ResolveUDPAddr("udp", i.FullAddress())
	return
}

func (w *Warping) warping(ip *UDPAddr) (bool, time.Duration) {

	fullAddress := ip.FullAddress()
	conn, err := net.DialTimeout("udp", fullAddress, udpConnectTimeout)
	if err != nil {
		return false, 0
	}
	defer conn.Close()
	startTime := time.Now()

	_, err = conn.Write(warpHandshakePacket)
	if err != nil {
		return false, 0
	}

	revBuff := make([]byte, 1024)

	err = conn.SetDeadline(time.Now().Add(udpConnectTimeout))
	if err != nil {
		return false, 0
	}
	n, err := conn.Read(revBuff)
	if err != nil {
		return false, 0
	}
	handshakeResponse := hex.EncodeToString(revBuff[:n])
	if handshakeResponse != warpValidatePacket {
		return false, 0
	}

	duration := time.Since(startTime)
	return true, duration
}
