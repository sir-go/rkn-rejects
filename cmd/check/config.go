package main

// Running configuration. Parses running flags to the config struct.

import (
	"encoding/json"
	"flag"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

// WorkersMax - hardcoded maximum amount of checkers
const WorkersMax = 75

type (
	// Cfg stores the whole running configuration
	Cfg struct {
		// checking workers amount
		Workers int `json:"workers,omitempty"`

		// pause between checks for the worker
		Sleeps time.Duration `json:"sleeps,omitempty"`

		// maximum amount of targets to check (limits the target list)
		Max int `json:"limit,omitempty"`

		// timeout for the response for each target
		Timeout time.Duration `json:"timeout,omitempty"`

		// how often to store the buffered log
		LogInterval time.Duration `json:"log_interval,omitempty"`

		// redis b connection parameters
		Redis struct {
			Host        string        `json:"host,omitempty"`
			Port        int           `json:"port,omitempty"`
			Password    string        `json:"-,omitempty"`
			Db          int           `json:"db,omitempty"`
			SetKey      string        `json:"set_key,omitempty"`
			TimeoutConn time.Duration `json:"timeout_conn,omitempty"`
			TimeoutRead time.Duration `json:"timeout_read,omitempty"`
		} `json:"redis"`

		// path of the directory to store verdicts
		VerdictsOutDir string `json:"verdicts_out_dir,omitempty"`

		// logging level
		LogLevel string `json:"log_level,omitempty"`
	}
)

func (c *Cfg) String() string {
	var (
		b   []byte
		err error
	)
	b, err = json.Marshal(c)
	if err != nil {
		log.Warnln("config Marshal error:", err.Error())
		return ""
	}
	return string(b)
}

func initConfig() *Cfg {
	cfg := new(Cfg)
	flag.StringVar(&cfg.Redis.Host, "rh", "localhost",
		"redis host")

	flag.IntVar(&cfg.Redis.Port, "rp", 6379,
		"redis port")

	flag.StringVar(&cfg.Redis.Password, "ra", "",
		"redis password")

	flag.IntVar(&cfg.Redis.Db, "rd", 0,
		"redis db")

	flag.StringVar(&cfg.Redis.SetKey, "rk", "check",
		"redis set key for checks")

	flag.DurationVar(&cfg.Redis.TimeoutConn, "rtc", time.Second*15,
		"redis connection timeout")

	flag.DurationVar(&cfg.Redis.TimeoutRead, "rtr", time.Second*15,
		"redis read timeout")

	flag.StringVar(&cfg.LogLevel, "log", "info",
		"log level [panic < fatal < error < warn < info < debug < trace]")

	flag.IntVar(&cfg.Workers, "w", 10,
		fmt.Sprintf("workers amount [1..%d]", WorkersMax))

	flag.DurationVar(&cfg.Sleeps, "s", time.Millisecond*5,
		"sleep between checks")

	flag.IntVar(&cfg.Max, "m", -1,
		"check limit (-1 - infinite)")

	flag.DurationVar(&cfg.Timeout, "t", time.Second*3,
		"tcp timeout")

	flag.DurationVar(&cfg.LogInterval, "lt", time.Second*10,
		"log progress interval")

	flag.StringVar(&cfg.VerdictsOutDir, "d", "/tmp/",
		"directory for per-record verdict logs")

	flag.Parse()

	if cfg.Workers > WorkersMax {
		log.Warnf("workers amount (-w) is too mutch, reduced to %d",
			WorkersMax)
		cfg.Workers = WorkersMax
	}

	log.Info(cfg)
	return cfg
}
