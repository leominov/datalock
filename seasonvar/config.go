package seasonvar

import "os"

const (
	DefaultHTTPAddr    = "127.0.0.1:7000"
	DefaultPublicDir   = "./public"
	DefaultMetricsPath = "/metrics"
	DefaultHealthzPath = "/healthz"
	DefaultDatabaseDir = "./database"
)

type Config struct {
	ListenAddr  string
	MetricsPath string
	HealthzPath string
	PublicDir   string
	DatabaseDir string
}

func NewConfig() *Config {
	return &Config{}
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
	return nil
}
