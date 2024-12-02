package config

type Config struct {
	DB            DBConfig     `json:"db"  yaml:"db"`
	Logger        LoggerConfig `json:"logger"  yaml:"logger"`
	Server        ServerConfig `json:"server"  yaml:"server"`
	Elasticsearch EsConfig     `json:"elasticsearch"  yaml:"elasticsearch"`
}

type DBConfig struct {
	Host      string `json:"host"  yaml:"host"`
	Port      string `json:"port"  yaml:"port"`
	QDatabase string `json:"q_database"  yaml:"q_database"`
	SDatabase string `json:"s_database"  yaml:"s_database"`
	Schema    string `json:"schema"  yaml:"schema"`
	Username  string `json:"username"  yaml:"username"`
	Password  string `json:"password"  yaml:"password"`
}

type LoggerConfig struct {
	Level  string `json:"level"  yaml:"level"`
	Output string `json:"output"  yaml:"output"`
	Path   string `json:"path"  yaml:"path"`
}

type ServerConfig struct {
	HttpPort            uint   `json:"http_port"  yaml:"http_port"`
	Secret              string `json:"secret"  yaml:"secret"`
	AuthExpMinute       uint   `json:"auth_exp_minute"  yaml:"auth_exp_minute"`
	AuthRefreshMinute   uint   `json:"auth_refresh_minute"  yaml:"auth_refresh_minute"`
	RateLimitMaxAttempt int    `json:"rate_limit_max_attempt"  yaml:"rate_limit_max_attempt"`
	RatelimitTimePeriod int    `json:"ratelimit_time_period"  yaml:"ratelimit_time_period"`
}

type EsConfig struct {
	Host     string `json:"host"  yaml:"host"`
	Port     string `json:"port"  yaml:"port"`
	Username string `json:"username"  yaml:"username"`
	Password string `json:"password"  yaml:"password"`
	Index    string `json:"index"  yaml:"index"`
}
