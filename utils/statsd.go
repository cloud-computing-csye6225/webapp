package utils

import (
	"fmt"
	"github.com/smira/go-statsd"
	"sync"
	"webapp/config"
)

var once sync.Once
var statsdClient *statsd.Client

func InitStatsdClient(configs config.Config) *statsd.Client {
	once.Do(func() {
		statsdClient = statsd.NewClient(
			fmt.Sprintf("%s:%s", configs.StatsDConfig.Host, configs.StatsDConfig.Port),
			statsd.MaxPacketSize(1400),
			statsd.MetricPrefix("api."))
	})
	return statsdClient
}

func StatIncrement(stat string, incrementBy int64) {
	statsdClient.Incr(stat, incrementBy)
}
