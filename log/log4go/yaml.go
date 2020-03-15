package log4go

import (
	"gopkg.in/yaml.v2"
)

type YamlConsoleConfig struct {
	Enable_  bool   `yaml:"enable"`
	Level_   string `yaml:"level"`
	Pattern_ string `yaml:"pattern"`
}

func (h* YamlConsoleConfig) Enable() bool { return h.Enable_ }
func (h* YamlConsoleConfig) Level() string { return h.Level_ }
func (h* YamlConsoleConfig) Pattern() string { return h.Pattern_ }

type YamlFileConfig struct {
	Enable_   bool   `yaml:"enable"`
	Category_ string `yaml:"category"`
	Level_    string `yaml:"level"`
	Filename_ string `yaml:"filename"`

	Pattern_ string `yaml:"pattern"`

	Rotate_   bool   `yaml:"rotate"`
	Maxsize_  string `yaml:"maxsize"`  // \d+[KMG]? Suffixes are in terms of 2**10
	Maxlines_ string `yaml:"maxlines"` //\d+[KMG]? Suffixes are in terms of thousands
	Daily_    bool   `yaml:"daily"`    //Automatically rotates by day
	Sanitize_ bool   `yaml:"sanitize"` //Sanitize newlines to prevent log injection
}

func (h* YamlFileConfig) Enable() bool { return h.Enable_ }
func (h* YamlFileConfig) Category() string { return h.Category_ }
func (h* YamlFileConfig) Level() string { return h.Level_ }
func (h* YamlFileConfig) Filename() string { return h.Filename_ }
func (h* YamlFileConfig) Pattern() string { return h.Pattern_ }
func (h* YamlFileConfig) Rotate() bool { return h.Rotate_ }
func (h* YamlFileConfig) Maxsize() string { return h.Maxsize_ }
func (h* YamlFileConfig) Maxlines() string { return h.Maxlines_ }
func (h* YamlFileConfig) Daily() bool { return h.Daily_ }
func (h* YamlFileConfig) Sanitize() bool { return h.Sanitize_ }

type YamlSocketConfig struct {
	Enable_   bool   `yaml:"enable"`
	Category_ string `yaml:"category"`
	Level_    string `yaml:"level"`
	Pattern_  string `yaml:"pattern"`
	Addr_     string `yaml:"addr"`
	Protocol_ string `yaml:"protocol"`
}

func (h* YamlSocketConfig) Enable() bool { return h.Enable_ }
func (h* YamlSocketConfig) Category() string { return h.Category_ }
func (h* YamlSocketConfig) Level() string { return h.Level_ }
func (h* YamlSocketConfig) Pattern() string { return h.Pattern_ }
func (h* YamlSocketConfig) Addr() string { return h.Addr_ }
func (h* YamlSocketConfig) Protocol() string { return h.Protocol_ }


// LogConfig presents json log config struct
type YamlLogConfig struct {
	Console_ *YamlConsoleConfig  `yaml:"console"`
	Files_   []*YamlFileConfig   `yaml:"files"`
	Sockets_ []*YamlSocketConfig `yaml:"sockets"`
}

func (h* YamlLogConfig) Console() IConsoleConfig { return h.Console_ }
func (h* YamlLogConfig) Files() []IFileConfig {
	ret := make([]IFileConfig, len(h.Files_))
	for i, f := range h.Files_ { ret[i] = f }
	return ret
}
func (h* YamlLogConfig) Sockets() []ISocketConfig {
	ret := make([]ISocketConfig, len(h.Sockets_))
	for i, s := range h.Sockets_ { ret[i] = s }
	return ret
}


func ParseYamlConfig(config string) (*YamlLogConfig, error) {
	var lc YamlLogConfig
	err := yaml.Unmarshal([]byte(config), &lc)
	return &lc, err
}