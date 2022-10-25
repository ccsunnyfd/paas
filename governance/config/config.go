package config

import (
	"github.com/asim/go-micro/plugins/config/source/consul/v4"
	"go-micro.dev/v4/config"
	"strconv"
)

func GetConfig(host string, port int, prefix string) (config.Config, error) {
	consulSource := consul.NewSource(
		// optionally specify consul address; default to localhost:8500
		consul.WithAddress(host+":"+strconv.Itoa(port)),
		// optionally specify prefix; defaults to /micro/config
		consul.WithPrefix(prefix),
		// optionally strip the provided prefix from the keys, defaults to false
		consul.StripPrefix(true),
	)

	// Create new config
	conf, err := config.NewConfig()
	if err != nil {
		return conf, err
	}

	// Load consul source
	return conf, conf.Load(consulSource)
}
