package config

import (
	"encoding/json"
	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
)

func init()  {
	RegisterCodec(&YamlCodec{})
	RegisterCodec(&JSONCodec{})
	RegisterCodec(&TomlCodec{})
}

// YamlCodec 解码Yaml
type YamlCodec struct{}

// Name yaml codec
func (*YamlCodec) Name() string {
	return "yaml"
}

// Unmarshal yaml decode
func (c *YamlCodec) Unmarshal(in []byte, out interface{}) error {
	return yaml.Unmarshal(in, out)
}

// JSONCodec JSON codec
type JSONCodec struct{}

// Name JSON codec
func (*JSONCodec) Name() string {
	return "json"
}

// Unmarshal JSON decode
func (c *JSONCodec) Unmarshal(in []byte, out interface{}) error {
	return json.Unmarshal(in, out)
}

// TomlCodec toml codec
type TomlCodec struct{}

// Name toml codec
func (*TomlCodec) Name() string {
	return "toml"
}

// Unmarshal toml decode
func (c *TomlCodec) Unmarshal(in []byte, out interface{}) error {
	return toml.Unmarshal(in, out)
}
