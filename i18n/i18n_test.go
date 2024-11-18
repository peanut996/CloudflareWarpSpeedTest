package i18n

import (
	"os"
	"testing"
)

func TestQueryI18n(t *testing.T) {
	tests := []struct {
		name     string
		lang     string
		messageID string
		want     string
	}{
		{
			name:     "english message",
			lang:     "en",
			messageID: TestThreadCount,
			want:     "Latency test threads; the more threads, the faster the latency test, but do not set it too high on low-performance devices (such as routers); [maximum 1000]",
		},
		{
			name:     "chinese message",
			lang:     "zh",
			messageID: TestThreadCount,
			want:     "指定延迟测试线程的数量。增加此值可以加快延迟测试过程，但不适合性能较低的设备，如路由器 [默认值为 200，最大为 1000]",
		},
		{
			name:     "default to english for unknown language",
			lang:     "fr",
			messageID: TestThreadCount,
			want:     "Latency test threads; the more threads, the faster the latency test, but do not set it too high on low-performance devices (such as routers); [maximum 1000]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save current LANG environment variable
			oldLang := os.Getenv("LANG")
			defer os.Setenv("LANG", oldLang)

			// Set test language
			os.Setenv("LANG", tt.lang)
			
			// Reset localizer for new language
			localizer = nil
			
			got := QueryI18n(tt.messageID)
			if got != tt.want {
				t.Errorf("QueryI18n() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueryTemplateI18n(t *testing.T) {
	tests := []struct {
		name      string
		lang      string
		messageID string
		tpData    map[string]interface{}
		want      string
	}{
		{
			name:      "english template",
			lang:      "en",
			messageID: OutputResultFile,
			tpData: map[string]interface{}{
				"file": "test.csv",
			},
			want: "Write result to file; add quotes if the path contains spaces; empty value means not writing to a file [-o \"\"]; ",
		},
		{
			name:      "chinese template",
			lang:      "zh",
			messageID: OutputResultFile,
			tpData: map[string]interface{}{
				"file": "test.csv",
			},
			want: "设置输出结果文件 [默认文件为 \"result.csv\"]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save current LANG environment variable
			oldLang := os.Getenv("LANG")
			defer os.Setenv("LANG", oldLang)

			// Set test language
			os.Setenv("LANG", tt.lang)
			
			// Reset localizer for new language
			localizer = nil
			
			got := QueryTemplateI18n(tt.messageID, tt.tpData)
			if got != tt.want {
				t.Errorf("QueryTemplateI18n() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInit(t *testing.T) {
	tests := []struct {
		name string
		lang string
	}{
		{
			name: "initialize with english",
			lang: "en",
		},
		{
			name: "initialize with chinese",
			lang: "zh",
		},
		{
			name: "initialize with unknown language",
			lang: "fr",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save current LANG environment variable
			oldLang := os.Getenv("LANG")
			defer os.Setenv("LANG", oldLang)

			// Set test language
			os.Setenv("LANG", tt.lang)
			
			// Reset and reinitialize
			localizer = nil
			initLocalizer()
			
			// Verify localizer is initialized
			if localizer == nil {
				t.Error("initLocalizer() failed to initialize localizer")
			}
		})
	}
}
