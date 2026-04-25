package scanner

import "time"

type Option func(config *Config)

type Config struct {
	dialTimeout time.Duration
	workerCount int
}

func WithWorkerCount(workerCount int) func(config *Config) {
	return func(conf *Config) {
		conf.workerCount = workerCount
	}
}
func WithDialTimeout(timeout time.Duration) func(config *Config) {
	return func(conf *Config) {
		conf.dialTimeout = timeout
	}
}
