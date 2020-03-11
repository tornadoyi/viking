package log4go

import (
	"gopkg.in/yaml.v2"
)

type _ConsoleConfig struct {
	Enable_  bool   `yaml:"enable"`
	Level_   string `yaml:"level"`
	Pattern_ string `yaml:"pattern"`
}

func (h* _ConsoleConfig) Enable() bool { return h.Enable_ }
func (h* _ConsoleConfig) Level() string { return h.Level_ }
func (h* _ConsoleConfig) Pattern() string { return h.Pattern_ }

type _FileConfig struct {
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

func (h* _FileConfig) Enable() bool { return h.Enable_ }
func (h* _FileConfig) Category() string { return h.Category_ }
func (h* _FileConfig) Level() string { return h.Level_ }
func (h* _FileConfig) Filename() string { return h.Filename_ }
func (h* _FileConfig) Pattern() string { return h.Pattern_ }
func (h* _FileConfig) Rotate() bool { return h.Rotate_ }
func (h* _FileConfig) Maxsize() string { return h.Maxsize_ }
func (h* _FileConfig) Maxlines() string { return h.Maxlines_ }
func (h* _FileConfig) Daily() bool { return h.Daily_ }
func (h* _FileConfig) Sanitize() bool { return h.Sanitize_ }

type _SocketConfig struct {
	Enable_   bool   `yaml:"enable"`
	Category_ string `yaml:"category"`
	Level_    string `yaml:"level"`
	Pattern_  string `yaml:"pattern"`
	Addr_     string `yaml:"addr"`
	Protocol_ string `yaml:"protocol"`
}

func (h* _SocketConfig) Enable() bool { return h.Enable_ }
func (h* _SocketConfig) Category() string { return h.Category_ }
func (h* _SocketConfig) Level() string { return h.Level_ }
func (h* _SocketConfig) Pattern() string { return h.Pattern_ }
func (h* _SocketConfig) Addr() string { return h.Addr_ }
func (h* _SocketConfig) Protocol() string { return h.Protocol_ }


// LogConfig presents json log config struct
type _LogConfig struct {
	Console_ *_ConsoleConfig  `yaml:"console"`
	Files_   []*_FileConfig   `yaml:"files"`
	Sockets_ []*_SocketConfig `yaml:"sockets"`
}

func (h* _LogConfig) Console() IConsoleConfig { return h.Console_ }
func (h* _LogConfig) Files() []IFileConfig {
	ret := make([]IFileConfig, len(h.Files_))
	for i, f := range h.Files_ { ret[i] = f }
	return ret
}
func (h* _LogConfig) Sockets() []ISocketConfig {
	ret := make([]ISocketConfig, len(h.Sockets_))
	for i, s := range h.Sockets_ { ret[i] = s }
	return ret
}


func (log Logger) InitWithYamlConfig(config string) {
	var lc _LogConfig
	if err := yaml.Unmarshal([]byte(config), lc); err != nil { panic(err) }
	log.InitWithConfig(&lc)
}
