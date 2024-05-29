package config

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
)

type Config struct {
	PublisherService ServiceConfig    `yaml:"publisher_service"`
	OrderService     ServiceConfig    `yaml:"order_service"`
	HTTPServer       HTTPServerConfig `yaml:"http_server"`
	Database         DatabaseConfig   `yaml:"database"`
	CacheConfig      CacheConfig      `yaml:"cache_config"`
	FrontendPath     string           `yaml:"frontendPath"`
}

type ServiceConfig struct {
	Addr       string           `yaml:"addr"`
	Subject    string           `yaml:"subject"`
	Publisher  PublisherConfig  `yaml:"publisher"`
	Subscriber SubscriberConfig `yaml:"subscriber"`
}

type PublisherConfig struct {
	ClusterID string `yaml:"clusterID"`
	ClientID  string `yaml:"clientID"`
}

type SubscriberConfig struct {
	ClusterID string           `yaml:"clusterID"`
	ClientID  string           `yaml:"clientID"`
	Options   SubscribeOptions `yaml:"options"`
}

type SubscribeOptions struct {
	StartOpt         StartOptConfig `yaml:"startOpt"`
	Subject          string         `yaml:"subject"`
	Qgroup           string         `yaml:"qgroup"`
	DurableName      string         `yaml:"durable_name"`
	SetManualAckMode bool           `yaml:"set_manual_ack_mode"`
	AckWait          string         `yaml:"ack_wait"`
	StartSeq         uint64         `yaml:"start_seq"`
	DeliverAll       bool           `yaml:"deliver_all"`
	DeliverLast      bool           `yaml:"deliver_last"`
	StartDelta       string         `yaml:"start_delta"`
	NewOnly          bool           `yaml:"new_only"`
}

type StartOptConfig struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type DatabaseConfig struct {
	URL string `yaml:"URL"`
}

type HTTPServerConfig struct {
	Addr         string        `yaml:"address" env-default:"0.0.0.0:8888"`
	Timeout      time.Duration `yaml:"timeout" env-default:"5s"`
	Idle_timeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type CacheConfig struct {
	Expiration   time.Duration `yaml:"expiration"`
	IntervalGC   time.Duration `yaml:"interval_gc"`
	MaxItems     int           `yaml:"max_items"`
	MaxItemSize  int           `yaml:"max_item_size"`
	MaxKeySize   int           `yaml:"max_key_size"`
	MaxCacheSize int           `yaml:"max_cache_size"`
}

func Read() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v\n", err)
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable is not set")
	}
	file, err := os.Open(configPath)
	if err != nil {
		return Config{}, err
	}

	out, err := io.ReadAll(file)
	if err != nil {
		log.Fatal("err:", err)
	}
	var config Config
	if err := yaml.Unmarshal(out, &config); err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	config.HTTPServer.PrintHTTPServerConfig()
	config.CacheConfig.PrintCacheConfig()

	return config, nil
}

func (c *CacheConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var raw struct {
		Expiration   string `yaml:"expiration"`
		IntervalGC   string `yaml:"interval_gc"`
		MaxItems     int    `yaml:"max_items"`
		MaxItemSize  int    `yaml:"max_item_size"`
		MaxKeySize   int    `yaml:"max_key_size"`
		MaxCacheSize int    `yaml:"max_cache_size"`
	}
	if err := unmarshal(&raw); err != nil {
		return err
	}
	expiration, err := parseDuration(raw.Expiration)
	if err != nil {
		return fmt.Errorf("failed to parse expiration duration: %w", err)
	}
	intervalGC, err := parseDuration(raw.IntervalGC)
	if err != nil {
		return fmt.Errorf("failed to parse interval_gc duration: %w", err)
	}
	c.Expiration = expiration
	c.IntervalGC = intervalGC
	c.MaxItems = raw.MaxItems
	c.MaxItemSize = raw.MaxItemSize
	c.MaxKeySize = raw.MaxKeySize
	c.MaxCacheSize = raw.MaxCacheSize
	return nil
}

func parseDuration(durationStr string) (time.Duration, error) {
	if durationStr == "" {
		return 0, nil
	}

	durationStr = strings.TrimSpace(durationStr)

	units := map[string]time.Duration{
		"ns": time.Nanosecond,
		"us": time.Microsecond,
		"µs": time.Microsecond,
		"ms": time.Millisecond,
		"s":  time.Second,
		"m":  time.Minute,
		"h":  time.Hour,
		"d":  time.Hour * 24,
	}

	for unit, duration := range units {
		if strings.HasSuffix(durationStr, unit) {
			valStr := strings.TrimSuffix(durationStr, unit)
			val, err := strconv.ParseFloat(valStr, 64)
			if err != nil {
				return 0, err
			}
			return time.Duration(val) * duration, nil
		}
	}

	return 0, fmt.Errorf("unknown duration format: %s", durationStr)
}

func (cc *CacheConfig) PrintCacheConfig() {
	fmt.Println("Cache Configuration:")
	fmt.Printf("  Expiration:     %s\n", cc.Expiration)
	fmt.Printf("  Interval GC:    %s\n", cc.IntervalGC)
	fmt.Printf("  Max Items:      %d\n", cc.MaxItems)
	fmt.Printf("  Max Item Size:  %d\n", cc.MaxItemSize)
	fmt.Printf("  Max Key Size:   %d\n", cc.MaxKeySize)
	fmt.Printf("  Max Cache Size: %d\n", cc.MaxCacheSize)
}

func (h *HTTPServerConfig) PrintHTTPServerConfig() {
	fmt.Println("HTTP Server Configuration:")
	fmt.Printf("  Addr:           %s\n", h.Addr)
	fmt.Printf("  Timeout:        %s\n", h.Timeout)
	fmt.Printf("  Idle_Timeout:   %s\n", h.Idle_timeout)
}
