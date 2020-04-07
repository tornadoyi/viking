package config

import (
	"github.com/tornadoyi/viking/goplus/path/file"
	"github.com/tornadoyi/viking/goplus/path/filepath"
	"github.com/tornadoyi/viking/goplus/runtime"
	"github.com/tornadoyi/viking/log"
)


func AddFileConfig(name, filePath, defaultConfig string, parseFunc interface{}, args... interface{}) (*Config, error) {
	return AddParser(NewFileParser(name, filePath, defaultConfig, parseFunc, args...))
}


type FileParser struct {
	*BaseParser
	filePath				string
	parser					*runtime.JITFunc
}

func NewFileParser(name string, filePath, defaultConfig string, parseFunc interface{}, args... interface{}) *FileParser {
	p := &FileParser{
		filePath:		filePath,
		parser:			runtime.NewJITFunc(parseFunc, args...),
	}
	p.BaseParser = NewBaseParser(name, defaultConfig, p.parse)
	return p
}

func (h *FileParser) FilePath() string { return h.filePath}


func (h *FileParser) parse() interface{} {

	configContent := h.defaultConfig

	// check config file
	if filepath.IsFile(h.filePath) {
		if data, err := file.Load(h.filePath); err != nil {
			log.Warnw("Activate default configuration due to config file load failure",
				"name", h.name,  "file", h.filePath, "error", err)
		} else {
			configContent = string(data)
		}
	}

	data, err := h.parser.Call(configContent)
	if err != nil { panic(err) }
	return data
}