package config

type Config struct {
	DB     DBConfig     `json:"db"  yaml:"db"`
	Logger LoggerConfig `json:"logger"  yaml:"logger"`
	Server ServerConfig `json:"server"  yaml:"server"`
}

type DBConfig struct {
	Host     string `json:"host"  yaml:"host"`
	Port     string `json:"port"  yaml:"port"`
	Database string `json:"database"  yaml:"database"`
	Schema   string `json:"schema"  yaml:"schema"`
	Username string `json:"username"  yaml:"username"`
	Password string `json:"password"  yaml:"password"`
}

type LoggerConfig struct {
	Path  string `json:"path"  yaml:"path"`
	Level string `json:"level"  yaml:"level"`
}

type ServerConfig struct {
	IPAddress			string `json:"ip_address" yaml:"ip_address"`
	HttpPort            uint   `json:"http_port"  yaml:"http_port"`
	Secret              string `json:"secret"  yaml:"secret"`
	AuthExpMinute       uint   `json:"auth_exp_minute"  yaml:"auth_exp_minute"`
	AuthRefreshMinute   uint   `json:"auth_refresh_minute"  yaml:"auth_refresh_minute"`
	RateLimitMaxAttempt int    `json:"rate_limit_max_attempt"  yaml:"rate_limit_max_attempt"`
	RatelimitTimePeriod int    `json:"ratelimit_time_period"  yaml:"ratelimit_time_period"`
}
