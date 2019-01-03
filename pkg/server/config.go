package server

import (
	"os"
	"strings"

	"github.com/leominov/datalock/pkg/backends"
	"github.com/leominov/datalock/pkg/util/text"
)

const (
	DefaultHTTPAddr       = "127.0.0.1:7000"
	DefaultPublicDir      = "./public"
	DefaultMetricsPath    = "/metrics"
	DefaultHealthzPath    = "/healthz"
	DefaultHdHostname     = "ZGF0YS1oZC5kYXRhbG9jay5ydQ=="
	DefaultHostname       = "c2Vhc29udmFyLnJ1"
	DefaultPublicHostname = "MXNlYXNvbnZhci5ydQ=="
	DefaultTemplatesDir   = "./templates"
)

type Config struct {
	ListenAddr          string
	MetricsPath         string
	HealthzPath         string
	PublicDir           string
	HdHostname          string
	Hostname            string
	PublicHostname      string
	TemplatesDir        string
	NodeList            []string
	StorageClientConfig backends.Config
}

func NewConfig() *Config {
	return &Config{
		Hostname:            text.Base64Decode(DefaultHostname),
		StorageClientConfig: backends.Config{},
	}
}

func (c *Config) Load() error {
	err := c.LoadFromEnv()
	if err != nil {
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
	c.HdHostname = os.Getenv("DATALOCK_HD_HOSTNAME")
	if c.HdHostname == "" {
		c.HdHostname = text.Base64Decode(DefaultHdHostname)
	}
	c.PublicHostname = os.Getenv("DATALOCK_PUBLIC_HOSTNAME")
	if c.PublicHostname == "" {
		c.PublicHostname = text.Base64Decode(DefaultPublicHostname)
	}
	c.TemplatesDir = os.Getenv("DATALOCK_TEMPLATES_DIR")
	if c.TemplatesDir == "" {
		c.TemplatesDir = DefaultTemplatesDir
	}
	nodeListStr := os.Getenv("DATALOCK_NODE_LIST")
	if len(nodeListStr) != 0 {
		nodeListArr := strings.Split(nodeListStr, ",")
		for _, node := range nodeListArr {
			nodeAddr := strings.TrimSpace(node)
			c.NodeList = append(c.NodeList, nodeAddr)
		}
	}
	c.StorageClientConfig.Backend = os.Getenv("DATALOCK_STORAGE_BACKEND")
	c.StorageClientConfig.Directory = os.Getenv("DATALOCK_STORAGE_DIRECTORY")
	c.StorageClientConfig.Bucket = os.Getenv("DATALOCK_STORAGE_BUCKET")
	if databaseDir := os.Getenv("DATALOCK_DATABASE_DIR"); databaseDir != "" {
		c.StorageClientConfig.Directory = databaseDir
	}
	return nil
}
