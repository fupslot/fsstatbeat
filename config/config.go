// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Resource struct {
	Id          string
	Description string
	Path        string
	Condition   string
}

type Config struct {
	Period   time.Duration `config:"period"`
	Resource Resource      `config:"resource"`
}

var DefaultConfig = Config{
	Period:   1 * time.Second,
	Resource: Resource{},
}
