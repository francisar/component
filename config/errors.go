package config

import "errors"

var (
	ErrConfigNotSupport = errors.New("config: not support")
	ErrProviderNotExist = errors.New("config: provider not exist")
	ErrCodecNotExist = errors.New("config: codec not exist")
)
