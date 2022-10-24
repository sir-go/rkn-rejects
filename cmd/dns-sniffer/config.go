package main

import (
	"encoding/json"
	"flag"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type (
	Cfg struct {
		Queues   []uint16 `json:"queues,omitempty"`
		QMaxLen  uint32   `json:"q_max_len,omitempty"`
		MarkDone int      `json:"mark,omitempty"`
		Redis    struct {
			Host        string        `json:"host,omitempty"`
			Port        int           `json:"port,omitempty"`
			Password    string        `json:"-,omitempty"`
			Db          int           `json:"db,omitempty"`
			SetKey      string        `json:"set_key,omitempty"`
			TimeoutConn time.Duration `json:"timeout_conn,omitempty"`
			TimeoutRead time.Duration `json:"timeout_read,omitempty"`
		} `json:"redis"`
		LogLevel string `json:"log_level,omitempty"`
		Dry      bool   `json:"dry,omitempty"`
		Nf       struct {
			Table   string `json:"table,omitempty"`
			SetName string `json:"set_name,omitempty"`
		} `json:"nf,omitempty"`
	}
)

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

// 100-103 -> [100 101 102 103]
func parseRange(s string, a *[]uint16) {
	var (
		v0, v1 uint64
		err    error
	)
	if !strings.ContainsRune(s, '-') {
		v0, err = strconv.ParseUint(s, 10, 16)
		if err != nil {
			log.Panicln("range parsing error:", err.Error())
		}
		*a = []uint16{uint16(v0)}
		return
	}

	p := strings.Split(s, "-")
	if v0, err = strconv.ParseUint(p[0], 10, 16); err != nil {
		log.Panicln("range begin parsing error:", err.Error())
	}
	if v1, err = strconv.ParseUint(p[1], 10, 16); err != nil {
		log.Panicln("range end parsing error:", err.Error())
	}
	for v0 <= v1 {
		*a = append(*a, uint16(v0))
		v0++
	}
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

	parseRange(queuesStr, &cfg.Queues)
	cfg.QMaxLen = uint32(qMaxLen)
	log.Info(cfg)
	if cfg.Dry {
		os.Exit(0)
	}
	return cfg
}
