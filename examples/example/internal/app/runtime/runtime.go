package runtime

import (
	"github.com/google/wire"
	"github.com/iobrother/zoo/core/config"
	zdb "github.com/iobrother/zoo/examples/example/pkg/database/db"
	zredis "github.com/iobrother/zoo/examples/example/pkg/database/redis"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var configProviderSet = wire.NewSet(
	config.Default,
	NewConfig,
	wire.FieldsOf(
		new(*Config),
		"Mysql",
		"Redis",
	),
)

type Config struct {
	Mysql zdb.Config    `yaml:"mysql"`
	Redis zredis.Config `yaml:"redis"`
}

func NewConfig(c config.Config) (*Config, error) {
	var conf Config
	err := c.Unmarshal(&conf)
	return &conf, err
}

var RuntimeProviderSet = wire.NewSet(
	configProviderSet,
	zdb.Open,
	zredis.NewClient,
	wire.NewSet(wire.Struct(new(Runtime), "*")),
)

type Runtime struct {
	Config *Config
	DB     *gorm.DB
	RC     *redis.Client
}
