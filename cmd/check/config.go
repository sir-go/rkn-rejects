package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

const WorkersMax = 75

type (
	Cfg struct {
		Workers     int           `json:"workers,omitempty"`
		Sleeps      time.Duration `json:"sleeps,omitempty"`
		Max         int           `json:"limit,omitempty"`
		Timeout     time.Duration `json:"timeout,omitempty"`
		LogInterval time.Duration `json:"log_interval,omitempty"`
		Redis       struct {
			Host        string        `json:"host,omitempty"`
			Port        int           `json:"port,omitempty"`
			Password    string        `json:"-,omitempty"`
			Db          int           `json:"db,omitempty"`
			SetKey      string        `json:"set_key,omitempty"`
			TimeoutConn time.Duration `json:"timeout_conn,omitempty"`
			TimeoutRead time.Duration `json:"timeout_read,omitempty"`
		} `json:"redis"`
		Out            string `json:"out,omitempty"`
		VerdictsOutDir string `json:"verdicts_out_dir,omitempty"`
		LogLevel       string `json:"log_level,omitempty"`
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
		"radis connection timeout")

	flag.DurationVar(&cfg.Redis.TimeoutRead, "rtr", time.Second*15,
		"radis read timeout")

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

	flag.StringVar(&cfg.Out, "o", "stdout",
		"buffered log file path")

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
