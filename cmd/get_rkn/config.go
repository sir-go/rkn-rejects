package main

import (
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
)

const (
	ConfFileEnvVar  = "RKN_CONF"
	DefaultConfFile = "rkn.toml"
)

type (
	Net struct {
		*net.IPNet
	}

	Duration struct {
		time.Duration
	}

	Cfg struct {
		ActualVersions struct {
			Service string `toml:"service"`
			Dump    string `toml:"dump"`
			Doc     string `toml:"doc"`
		} `toml:"actual_versions"`
		Web struct {
			SoapUrl    string    `toml:"soap_url"`
			DocUrl     string    `toml:"doc_url"`
			TcpTimeout *Duration `toml:"tcp_timeout"`
			Attempts   int       `toml:"attempts"`
		} `toml:"web"`
		Req struct {
			File     string `toml:"file"`
			Operator struct {
				Name  string `toml:"name"`
				INN   string `toml:"inn"`
				OGRN  string `toml:"ogrn"`
				Email string `toml:"email"`
			} `toml:"operator"`
		} `toml:"req"`
		Sign struct {
			File   string `toml:"file"`
			Script string `toml:"script"`
		} `toml:"sign"`
		Res struct {
			DumpTo       *string   `toml:"dump_to"`
			Attempts     int       `toml:"attempts"`
			GetTimeout   *Duration `toml:"download_timeout"`
			RetryTimeout *Duration `toml:"retry_timeout"`
		} `toml:"res"`
		Parse struct {
			FromDump            *string   `toml:"from_dump"`
			ProgressPollTimeout *Duration `toml:"progress_poll_timeout"`
			BogusIp             struct {
				Subnets []Net `toml:"subnets"`
				MinMask int   `toml:"min_mask"`
			} `toml:"bogus_ip"`
			Redis struct {
				Host        string    `toml:"host"`
				Port        int       `toml:"port"`
				Password    string    `toml:"password"`
				Db          int       `toml:"db"`
				ChunkSize   int       `toml:"chunk_size"`
				Workers     int       `toml:"workers"`
				TimeoutConn *Duration `toml:"timeout_conn,omitempty"`
				TimeoutRead *Duration `toml:"timeout_read,omitempty"`
			} `toml:"redis"`
		} `toml:"parse"`
		Lists struct {
			BlackDomains string `toml:"black_domains"`
			WhiteDomains string `toml:"white_domains"`
		} `json:"lists"`
		Fw struct {
			IpDenyFile  string `toml:"ip_deny_file"`
			IpDenyTable string `toml:"ip_deny_table"`
			IpDenySet   string `toml:"ip_deny_set"`
		} `toml:"fw"`
		LogLevel     string `toml:"log_level"`
		WListDomains []string
		BListDomains []string
	}
)

func (d *Duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

func (n *Net) UnmarshalText(text []byte) error {
	return n.parse(string(text))
}

func (n *Net) parse(text string) error {
	var (
		err error
		np  *net.IPNet
	)
	if !strings.Contains(text, "/") {
		text += "/32"
	}
	_, np, err = net.ParseCIDR(text)
	n.IPNet = np
	return err
}

func initConfig() *Cfg {
	cfg := new(Cfg)
	cfgPath := os.Getenv(ConfFileEnvVar)
	if cfgPath == "" {
		cfgPath = DefaultConfFile
	}
	if _, err := toml.DecodeFile(cfgPath, cfg); err != nil {
		log.Panic(err)
	}
	absPath, err := filepath.Abs(cfgPath)
	if err != nil {
		log.Panic(err)
	}

	cfg.WListDomains = includesReadDomains(cfg.Lists.WhiteDomains)
	cfg.BListDomains = includesReadDomains(cfg.Lists.BlackDomains)

	log.Infoln("config:", absPath)

	return cfg
}
