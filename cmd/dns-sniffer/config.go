package main

import (
	"encoding/json"
	"flag"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"rkn-rejects/internal/tools"
)

type (
	Cfg struct {
		// nf queue IDs
		Queues []uint16 `json:"queues,omitempty"`

		// maximum queue capacity
		QMaxLen uint32 `json:"q_max_len,omitempty"`

		// 'done' packet marker ID
		MarkDone int `json:"mark,omitempty"`

		// redis db connection parameters
		Redis struct {
			Host        string        `json:"host,omitempty"`
			Port        int           `json:"port,omitempty"`
			Password    string        `json:"-,omitempty"`
			Db          int           `json:"db,omitempty"`
			SetKey      string        `json:"set_key,omitempty"`
			TimeoutConn time.Duration `json:"timeout_conn,omitempty"`
			TimeoutRead time.Duration `json:"timeout_read,omitempty"`
		} `json:"redis"`

		// logging level
		LogLevel string `json:"log_level,omitempty"`

		// just configure and exit
		Dry bool `json:"dry,omitempty"`

		// netfilter configuration
		Nf struct {
			Table   string `json:"table,omitempty"`
			SetName string `json:"set_name,omitempty"`
		} `json:"nf,omitempty"`
	}
)

// Config stringer
func (c *Cfg) String() string {
	var (
		b   []byte
		err error
	)
	if c.Dry {
		b, err = json.MarshalIndent(c, "  ", "  ")
	} else {
		b, err = json.Marshal(c)
	}
	if err != nil {
		log.Warnln("config Marshal error:", err.Error())
		return ""
	}
	return string(b)
}

func initConfig() *Cfg {
	var (
		queuesStr string
		qMaxLen   uint64
	)
	cfg := new(Cfg)
	flag.StringVar(&queuesStr, "nfq", "100-103",
		"nf queues range")

	flag.IntVar(&cfg.MarkDone, "mdone", 1,
		"nf mark done")

	flag.Uint64Var(&qMaxLen, "nfql", 0xFF,
		"max nf queue length")

	flag.StringVar(&cfg.Redis.Host, "rh", "localhost",
		"redis host")

	flag.IntVar(&cfg.Redis.Port, "rp", 6379,
		"redis port")

	flag.StringVar(&cfg.Redis.Password, "ra", "",
		"redis password")

	flag.IntVar(&cfg.Redis.Db, "rd", 0,
		"redis db")

	flag.StringVar(&cfg.Redis.SetKey, "rk", "domains",
		"redis domains set key")

	flag.DurationVar(&cfg.Redis.TimeoutConn, "rtc", time.Second*15,
		"radius connection timeout")

	flag.DurationVar(&cfg.Redis.TimeoutRead, "rtr", time.Second*15,
		"radius read timeout")

	flag.StringVar(&cfg.LogLevel, "log", "info",
		"log level [panic < fatal < error < warn < info < debug < trace]")

	flag.StringVar(&cfg.Nf.Table, "nft", "rkn",
		"nf tables table name")

	flag.StringVar(&cfg.Nf.SetName, "nfs", "allow_sniffed",
		"nf tables set name")

	flag.BoolVar(&cfg.Dry, "dry", false,
		"just pretty print config")

	flag.Parse()

	if err := tools.ParseRange(queuesStr, &cfg.Queues); err != nil {
		panic(err)
	}
	cfg.QMaxLen = uint32(qMaxLen)
	log.Info(cfg)
	if cfg.Dry {
		os.Exit(0)
	}
	return cfg
}
