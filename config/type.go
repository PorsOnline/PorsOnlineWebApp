package config

type Config struct {
	DB            DBConfig       `json:"db"  yaml:"db"`
	SecretDB      SecretDBConfig `json:"secretDB" yaml:"secretDB"`
	Logger        LoggerConfig   `json:"logger"  yaml:"logger"`
	Server        ServerConfig   `json:"server"  yaml:"server"`
	Elasticsearch EsConfig       `json:"elasticsearch"  yaml:"elasticsearch"`
}

type ServerConfig struct {
	HttpPort            uint   `json:"httpPort"  yaml:"httpPort"`
	Secret              string `json:"secret"  yaml:"secret"`
	AuthExpMinute       uint   `json:"authExpMin"  yaml:"authExpMin"`
	AuthRefreshMinute   uint   `json:"authExpRefreshMin"  yaml:"authExpRefreshMin"`
	RateLimitMaxAttempt int    `json:"rate_limit_max_attempt"  yaml:"rate_limit_max_attempt"`
	RatelimitTimePeriod int    `json:"ratelimit_time_period"  yaml:"ratelimit_time_period"`
}
type DBConfig struct {
	Host     string `json:"host"`
	Port     uint   `json:"port"`
	Database string `json:"database"`
	Schema   string `json:"schema"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type SecretDBConfig struct {
	Host     string `json:"host"`
	Port     uint   `json:"port"`
	Database string `json:"database"`
	Schema   string `json:"schema"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type LoggerConfig struct {
	Level  string `json:"level"  yaml:"level"`
	Output string `json:"output"  yaml:"output"`
	Path   string `json:"path"  yaml:"path"`
}

type EsConfig struct {
	Host     string `json:"host"  yaml:"host"`
	Port     string `json:"port"  yaml:"port"`
	Username string `json:"username"  yaml:"username"`
	Password string `json:"password"  yaml:"password"`
	Index    string `json:"index"  yaml:"index"`
}
