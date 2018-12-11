package config

import (
	"github.com/zhengcf/goutil/util/errors"
	"github.com/zhengcf/goutil/util/logutil"
	"github.com/BurntSushi/toml"
)

type Config struct {
	Host        string   `toml;"host" json:"host"`
	Port       	uint  	 `toml:"port" json:"port"`
	Psaddr  	string 	 `toml:"psaddr" json:"psaddr"`
	PrometheusPort uint  `toml:"prometheus_port" json:"prometheus_port"`
	Log         Log      `toml:"Log" json:"log"`
}


// Log is the log section of config.
type Log struct {
	// Log level.
	Level string `toml:"level" json:"level"`
	// Log format. one of json, text, or console.
	Format string `toml:"format" json:"format"`
	// Disable automatic timestamps in output.
	DisableTimestamp bool `toml:"disable-timestamp" json:"disable_timestamp"`
	// File log config.
	File logutil.FileLogConfig `toml:"file" json:"file"`

	SlowQueryFile      string `toml:"slow-query-file" json:"slow_query_file"`
	SlowThreshold      uint64 `toml:"slow-threshold" json:"slow_threshold"`
	ExpensiveThreshold uint   `toml:"expensive-threshold" json:"expensive_threshold"`
	QueryLogMaxLen     uint64 `toml:"query-log-max-len" json:"query_log_max_len"`
}

var defaultConf = Config{
	Host:      "0.0.0.0",
	Port:      8000,
	Psaddr:    "",
	PrometheusPort:  8080,
	Log: Log{
		Level:  "debug",
		Format: "text",
		File: logutil.FileLogConfig{
			LogRotate: true,
			MaxSize:   logutil.DefaultLogMaxSize,
		},
		SlowThreshold:      logutil.DefaultSlowThreshold,
		ExpensiveThreshold: 10000,
		QueryLogMaxLen:     logutil.DefaultQueryLogMaxLen,
	},
}

var globalConf = defaultConf

func NewConfig() *Config {
	conf := defaultConf
	return &conf
}

func GetGlobalConfig() *Config {
	return &globalConf
}

func(c *Config) Load(confFile string) error {
	_, err := toml.DecodeFile(confFile, c)
	return errors.Trace(err)
}

func (l *Log) ToLogConfig() *logutil.LogConfig {
	return &logutil.LogConfig{
		Level:            l.Level,
		Format:           l.Format,
		DisableTimestamp: l.DisableTimestamp,
		File:             l.File,
		SlowQueryFile:    l.SlowQueryFile,
	}
}
