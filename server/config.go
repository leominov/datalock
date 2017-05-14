package server

import (
	"os"

	"github.com/leominov/datalock/utils"
)

const (
	DefaultHTTPAddr    = "127.0.0.1:7000"
	DefaultPublicDir   = "./public"
	DefaultMetricsPath = "/metrics"
	DefaultHealthzPath = "/healthz"
	DefaultDatabaseDir = "./database"
	DefaultHdHostname  = "ZGF0YS1oZC5kYXRhbG9jay5ydQ=="
	DefaultHostname    = "c2Vhc29udmFyLnJ1"
)

type Config struct {
	ListenAddr  string
	MetricsPath string
	HealthzPath string
	PublicDir   string
	DatabaseDir string
	HdHostname  string
	Hostname    string
}

func NewConfig() *Config {
	return &Config{
		Hostname: utils.Base64Decode(DefaultHostname),
	}
}

func (c *Config) Load() error {
	if err := c.LoadFromEnv(); err != nil {
		return err
	}
	return nil
}

func (c *Config) LoadFromEnv() error {
	c.ListenAddr = os.Getenv("DATALOCK_LISTEN_ADDR")
	if c.ListenAddr == "" {
		c.ListenAddr = DefaultHTTPAddr
	}
	c.PublicDir = os.Getenv("DATALOCK_PUBLIC_DIR")
	if c.PublicDir == "" {
		c.PublicDir = DefaultPublicDir
	}
	c.MetricsPath = os.Getenv("DATALOCK_METRICS_PATH")
	if c.MetricsPath == "" {
		c.MetricsPath = DefaultMetricsPath
	}
	c.HealthzPath = os.Getenv("DATALOCK_HEALTHZ_PATH")
	if c.HealthzPath == "" {
		c.HealthzPath = DefaultHealthzPath
	}
	c.DatabaseDir = os.Getenv("DATALOCK_DATABASE_DIR")
	if c.DatabaseDir == "" {
		c.DatabaseDir = DefaultDatabaseDir
	}
	c.HdHostname = os.Getenv("DATALOCK_HD_HOSTNAME")
	if c.HdHostname == "" {
		c.HdHostname = utils.Base64Decode(DefaultHdHostname)
	}
	return nil
}
