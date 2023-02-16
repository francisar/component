package redis

import "time"

type Option struct {
	Address string `json:"address" yaml:"address"`
	IsCluster bool `json:"is_cluster" yaml:"is_cluster"`
	UserName string `json:"user_name" yaml:"user_name"`
	Password string `json:"password" yaml:"password"`
	Timeout  time.Duration `json:"timeout" yaml:"timeout"`
}