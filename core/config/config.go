package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type OnChangeFunc func(Config)

type Config interface {
	Unmarshal(val any) error
	Scan(key string, val any) error
	Get(key string) any
	GetString(key string) string
	GetBool(key string) bool
	GetInt(key string) int
	GetFloat64(key string) float64
	GetDuration(key string) time.Duration
	GetIntSlice(key string) []int
	GetStringSlice(key string) []string
	GetStringMap(key string) map[string]any
	OnChange(f OnChangeFunc)
}

type config struct {
	v        *viper.Viper
	data     []byte
	onChange OnChangeFunc
}

func New() Config {
	c := config{
		v: viper.New(),
	}

	c.load()

	return &c
}

func (c *config) load() {
	c.v.SetConfigName("config")

	c.v.AddConfigPath(".")
	c.v.AddConfigPath("./conf")
	if err := c.v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(err)
		}
	}

	env := c.v.GetString("env")
	if env == "" {
		env = "dev"
	}
	c.v.SetConfigName(fmt.Sprintf("config-%s", env))
	if err := c.v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(err)
		}
	}

	c.data, _ = json.Marshal(c.v.AllSettings())
	c.v.OnConfigChange(func(e fsnotify.Event) {
		data, _ := json.Marshal(c.v.AllSettings())
		if !bytes.Equal(data, c.data) {
			c.data = data
			c.onChange(c)
		}
	})

	c.v.WatchConfig()
}

func (c *config) Unmarshal(val any) error {
	return c.v.Unmarshal(val)
}

func (c *config) Scan(key string, val any) error {
	return c.v.UnmarshalKey(key, val)
}

func (c *config) Get(key string) any {
	return c.v.Get(key)
}

func (c *config) GetString(key string) string {
	return c.v.GetString(key)
}

func (c *config) GetBool(key string) bool {
	return c.v.GetBool(key)
}

func (c *config) GetInt(key string) int {
	return c.v.GetInt(key)
}

func (c *config) GetFloat64(key string) float64 {
	return c.v.GetFloat64(key)
}

func (c *config) GetDuration(key string) time.Duration {
	return c.v.GetDuration(key)
}

func (c *config) GetIntSlice(key string) []int {
	return c.v.GetIntSlice(key)
}

func (c *config) GetStringSlice(key string) []string {
	return c.v.GetStringSlice(key)
}

func (c *config) GetStringMap(key string) map[string]any {
	return c.v.GetStringMap(key)
}

func (c *config) OnChange(f OnChangeFunc) {
	c.onChange = f
}
