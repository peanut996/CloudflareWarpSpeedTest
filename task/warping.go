package task

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"golang.org/x/crypto/blake2s"
	"golang.org/x/crypto/poly1305"
	"golang.zx2c4.com/wireguard/tai64n"
	"log"
	"math/rand"
	"net"
	"net/netip"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/peanut996/CloudflareWarpSpeedTest/utils"

	"golang.zx2c4.com/wireguard/conn"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/tun/netstack"
)

const (
	defaultRoutines             = 200
	defaultPingTimes            = 10
	udpConnectTimeout           = time.Millisecond * 1000
	wireguardHandshakeRespBytes = 92
	quickModeMaxIpNum           = 1000
	warpPublicKey               = "bmXOC+F1FxEMF9dyiK2H5/1SUtzH0JuVo51h2wPfgyo="
)

var (
	PrivateKey string

	PublicKey string

	QuickMode = false

	IPv6Mode = false

	ScanAllPort = false

	ReservedString = ""

	reserved = [3]byte{60, 189, 175}

	Routines = defaultRoutines

	PingTimes = defaultPingTimes

	ports = []int{
		500, 854, 859, 864, 878, 880, 890, 891, 894, 903,
		908, 928, 934, 939, 942, 943, 945, 946, 955, 968,
		987, 988, 1002, 1010, 1014, 1018, 1070, 1074, 1180, 1387,
		1701, 1843, 2371, 2408, 2506, 3138, 3476, 3581, 3854, 4177,
		4198, 4233, 4500, 5279, 5956, 7103, 7152, 7156, 7281, 7559, 8319, 8742, 8854, 8886,
	}

	commonIPv4CIDRs = []string{
		"162.159.192.0/24",
		"162.159.193.0/24",
		"162.159.195.0/24",
		"162.159.204.0/24",
		"188.114.96.0/24",
		"188.114.97.0/24",
		"188.114.98.0/24",
		"188.114.99.0/24",
	}

	commonIPv6CIDRs = []string{
		"2606:4700:d0::/48",
	}

	MaxWarpPortRange = 10000

	warpHandshakePacket, _ = hex.DecodeString("013cbdafb4135cac96a29484d7a0175ab152dd3e59be35049beadf758b8d48af14ca65f25a168934746fe8bc8867b1c17113d71c0fac5c141ef9f35783ffa5357c9871f4a006662b83ad71245a862495376a5fe3b4f2e1f06974d748416670e5f9b086297f652e6dfbf742fbfc63c3d8aeb175a3e9b7582fbc67c77577e4c0b32b05f92900000000000000000000000000000000")
)

type MessageInitiation struct {
	Type      uint8
	Reserved  [3]byte
	Sender    uint32
	Ephemeral device.NoisePublicKey
	Static    [device.NoisePublicKeySize + poly1305.TagSize]byte
	Timestamp [tai64n.TimestampSize + poly1305.TagSize]byte
	MAC1      [blake2s.Size128]byte
	MAC2      [blake2s.Size128]byte
}

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
		bar:     utils.NewBar(len(ips), "Available:", ""),
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

func (w *Warping) warpingHandler(ip *UDPAddr) {
	recv, totalDelay := w.warping(ip)
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
	addrs := generateIPAddrs(ips)
	if QuickMode && len(addrs) > quickModeMaxIpNum {
		return addrs[:quickModeMaxIpNum]
	}
	return addrs
}

func generateIPAddrs(ips []*net.IPAddr) (udpAddrs []*UDPAddr) {
	if !ScanAllPort {
		for _, port := range ports {
			udpAddrs = append(udpAddrs, generateSingleIPAddr(ips, port)...)
		}
	} else {
		for port := 1; port <= MaxWarpPortRange; port++ {
			udpAddrs = append(udpAddrs, generateSingleIPAddr(ips, port)...)
		}
	}
	shuffleAddrs(&udpAddrs)
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

func (w *Warping) warping(ip *UDPAddr) (received int, totalDelay time.Duration) {
	fullAddress := ip.FullAddress()
	con, err := net.DialTimeout("udp", fullAddress, udpConnectTimeout)
	if err != nil {
		return 0, 0
	}
	defer con.Close()

	for i := 0; i < PingTimes; i++ {
		ok, rtt := handshake(con)
		if ok {
			received++
			totalDelay += rtt
		}
	}
	return

}

func handshake(conn net.Conn) (bool, time.Duration) {
	startTime := time.Now()
	_, err := conn.Write(warpHandshakePacket)
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
	if n != wireguardHandshakeRespBytes {
		return false, 0
	}

	duration := time.Since(startTime)
	return true, duration
}

func shuffleAddrs(udpAddrs *[]*UDPAddr) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	r.Shuffle(len(*udpAddrs), func(i, j int) {
		(*udpAddrs)[i], (*udpAddrs)[j] = (*udpAddrs)[j], (*udpAddrs)[i]
	})
}

func InitHandshakePacket() {
	if ReservedString != "" {
		if PrivateKey == "" {
			log.Fatalln("Reserved field must be used with private key")
		}
		r, err := utils.ParseReservedString(ReservedString)
		if err != nil {
			log.Fatalln("Failed to parse reserved, it must be 3 bytes slice: " + err.Error())
		}
		reserved = r
	}

	if PrivateKey == "" && PublicKey == "" {
		return
	}

	if PublicKey == "" {
		PublicKey = warpPublicKey
	}

	pri, err := getNoisePrivateKeyFromBase64(PrivateKey)
	if err != nil {
		log.Fatalln("Failed to parse private key: " + err.Error())
	}

	pub, err := getNoisePublicKeyFromBase64(PublicKey)
	if err != nil {
		log.Fatalln("Failed to parse public key: " + err.Error())
	}

	packet := buildHandshakePacket(pri, pub)

	warpHandshakePacket = packet[:]
}

func buildHandshakePacket(pri device.NoisePrivateKey, pub device.NoisePublicKey) []byte {
	d, _, err := netstack.CreateNetTUN([]netip.Addr{}, []netip.Addr{}, 1480)
	if err != nil {
		log.Fatalln("Failed to build handshake packet: " + err.Error())
	}
	dev := device.NewDevice(d, conn.NewDefaultBind(), device.NewLogger(0, ""))

	dev.SetPrivateKey(pri)

	peer, err := dev.NewPeer(pub)
	if err != nil {
		log.Fatalln("Failed to build handshake packet: " + err.Error())
	}
	msg, err := dev.CreateMessageInitiation(peer)
	if err != nil {
		log.Fatalln("Failed to build handshake packet: " + err.Error())
	}

	var buf [device.MessageInitiationSize]byte
	writer := bytes.NewBuffer(buf[:0])

	binary.Write(writer, binary.LittleEndian, msg)
	packet := writer.Bytes()

	generator := device.CookieGenerator{}
	generator.Init(pub)
	generator.AddMacs(packet)

	AddReserved(packet)
	return packet
}

func AddReserved(packet []byte) {
	packet[1], packet[2], packet[3] = reserved[0], reserved[1], reserved[2]
}

func getNoisePrivateKeyFromBase64(b string) (device.NoisePrivateKey, error) {
	pk := device.NoisePrivateKey{}
	h, err := encodeBase64ToHex(b)
	if err != nil {
		return pk, err
	}
	err = pk.FromHex(h)
	if err != nil {
		return pk, err
	}
	return pk, nil
}

func getNoisePublicKeyFromBase64(b string) (device.NoisePublicKey, error) {
	pk := device.NoisePublicKey{}
	h, err := encodeBase64ToHex(b)
	if err != nil {
		return pk, err
	}
	err = pk.FromHex(h)
	if err != nil {
		return pk, err
	}
	return pk, nil
}

func encodeBase64ToHex(key string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return "", errors.New("Invalid base64 string: " + key)
	}
	if len(decoded) != 32 {
		return "", errors.New("Noise key should be 32 bytes: " + key)
	}
	return hex.EncodeToString(decoded), nil
}
