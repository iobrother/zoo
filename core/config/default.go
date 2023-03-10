package config

var defaultConfig = New()

func Default() Config {
	return defaultConfig
}

func ResetDefault(c Config) {
	defaultConfig = c

	Unmarshal = defaultConfig.Unmarshal
	Scan = defaultConfig.Scan
	Get = defaultConfig.Get
	GetString = defaultConfig.GetString
	GetBool = defaultConfig.GetBool
	GetInt = defaultConfig.GetInt
	GetFloat64 = defaultConfig.GetFloat64
	GetDuration = defaultConfig.GetDuration
	GetIntSlice = defaultConfig.GetIntSlice
	GetStringSlice = defaultConfig.GetStringSlice
	GetStringMap = defaultConfig.GetStringMap
}

var (
	Unmarshal      = defaultConfig.Unmarshal
	Scan           = defaultConfig.Scan
	Get            = defaultConfig.Get
	GetString      = defaultConfig.GetString
	GetBool        = defaultConfig.GetBool
	GetInt         = defaultConfig.GetInt
	GetFloat64     = defaultConfig.GetFloat64
	GetDuration    = defaultConfig.GetDuration
	GetIntSlice    = defaultConfig.GetIntSlice
	GetStringSlice = defaultConfig.GetStringSlice
	GetStringMap   = defaultConfig.GetStringMap
	OnChange       = defaultConfig.OnChange
)
