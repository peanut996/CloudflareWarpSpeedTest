package mobile

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/peanut996/CloudflareWarpSpeedTest/task"
	"github.com/peanut996/CloudflareWarpSpeedTest/utils"
)

// SpeedTest represents the mobile interface for the speed test functionality
type SpeedTest struct {
	warping    *task.Warping
	isRunning  bool
	mutex      sync.Mutex
	resultChan chan string
}

// TestConfig represents the configuration for speed testing
type TestConfig struct {
	ThreadCount        int     `json:"threadCount"`
	PingTimes         int     `json:"pingTimes"`
	MaxScanCount      int     `json:"maxScanCount"`
	MaxDelay          int     `json:"maxDelay"`
	MinDelay          int     `json:"minDelay"`
	MaxLossRate       float64 `json:"maxLossRate"`
	TestAllCombos     bool    `json:"testAllCombos"`
	IPv6Mode          bool    `json:"ipv6Mode"`
	ResultDisplayCount int    `json:"resultDisplayCount"`
	CustomIPFile      string  `json:"customIpFile"`
	CustomIPText      string  `json:"customIpText"`
	PrivateKey        string  `json:"privateKey"`
	PublicKey         string  `json:"publicKey"`
}

// NewSpeedTest creates a new instance of SpeedTest
func NewSpeedTest() *SpeedTest {
	return &SpeedTest{
		resultChan: make(chan string, 100),
	}
}

// Configure sets up the speed test with the provided configuration
func (s *SpeedTest) Configure(configJSON string) error {
	var config TestConfig
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return err
	}

	task.Routines = config.ThreadCount
	task.PingTimes = config.PingTimes
	task.MaxScanCount = config.MaxScanCount
	task.AllMode = config.TestAllCombos
	task.IPv6Mode = config.IPv6Mode
	task.IPFile = config.CustomIPFile
	task.IPText = config.CustomIPText
	task.PrivateKey = config.PrivateKey
	task.PublicKey = config.PublicKey
	utils.PrintNum = config.ResultDisplayCount
	utils.InputMaxDelay = time.Duration(config.MaxDelay) * time.Millisecond
	utils.InputMinDelay = time.Duration(config.MinDelay) * time.Millisecond
	utils.InputMaxLossRate = float32(config.MaxLossRate)

	return nil
}

// Start begins the speed test
func (s *SpeedTest) Start() error {
	s.mutex.Lock()
	if s.isRunning {
		s.mutex.Unlock()
		return nil
	}
	s.isRunning = true
	s.mutex.Unlock()

	// Initialize the handshake packet
	task.InitHandshakePacket()

	// Create new warping instance
	s.warping = task.NewWarping()

	// Run in a goroutine
	go func() {
		// Run the speed test
		pingData := s.warping.Run().FilterDelay().FilterLossRate()
		
		// Convert results to JSON
		results, err := json.Marshal(pingData)
		if err != nil {
			s.resultChan <- `{"error": "` + err.Error() + `"}`
		} else {
			s.resultChan <- string(results)
		}

		s.mutex.Lock()
		s.isRunning = false
		s.mutex.Unlock()
	}()

	return nil
}

// Stop stops the current speed test
func (s *SpeedTest) Stop() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	if s.isRunning && s.warping != nil {
		// Implement stopping mechanism
		s.isRunning = false
	}
}

// GetResults returns the latest results as a JSON string
func (s *SpeedTest) GetResults() string {
	select {
	case result := <-s.resultChan:
		return result
	default:
		return `{"status": "running"}`
	}
}

// IsRunning returns whether a test is currently in progress
func (s *SpeedTest) IsRunning() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.isRunning
}
