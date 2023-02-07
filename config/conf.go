package config

import (
	"flag"
	"strings"
)

type Conf struct {
	Hosts     []string
	StoreFile string
	Debug     bool
}

func Load() *Conf {
	hostsFlag := flag.String("hosts", "ya.ru,google.com", "")
	storeFlag := flag.String("store-file", "store.db", "")
	debugFlag := flag.Bool("debug", false, "")
	flag.Parse()

	hosts := strings.Split(*hostsFlag, ",")

	c := &Conf{
		Hosts:     make([]string, len(hosts)),
		StoreFile: *storeFlag,
		Debug:     *debugFlag,
	}

	for i, host := range hosts {
		c.Hosts[i] = strings.Trim(host, " ")
	}

	return c
}
