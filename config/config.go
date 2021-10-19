// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type FileResource struct {
	Path string `config:"path"`
}

type ProcResource struct {
	Name string `config:"name"`
}

type Resource struct {
	File      FileResource `config:"file"`
	Proc      ProcResource `config:"process"`
	Condition string       `config:"condition"`
}

type Rule struct {
	Id          string     `config:"id"`
	Name        string     `config:"name"`
	Description string     `config:"description"`
	Resources   []Resource `config:"resources"`
}

type Config struct {
	Period time.Duration `config:"period"`
	Rules  []Rule        `config:"rules"`
}

var DefaultConfig = Config{
	Period: 1 * time.Second,
	Rules:  make([]Rule, 0),
}
