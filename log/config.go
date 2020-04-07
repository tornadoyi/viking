package log



type Config struct {
	Level 							Level 							`json:"level" yaml:"level"`
	DisableCaller 					bool 							`json:"disableCaller" yaml:"disableCaller"`
	DisableStacktrace 				bool 							`json:"disableStacktrace" yaml:"disableStacktrace"`
	Sampling 						*SamplingConfig 				`json:"sampling" yaml:"sampling"`
	Encoding 						string 							`json:"encoding" yaml:"encoding"`
	EncoderConfig 					EncoderConfig 					`json:"encoderConfig" yaml:"encoderConfig"`
	InitialFields 					map[string]interface{} 			`json:"initialFields" yaml:"initialFields"`
	File							*RotateLogger					`json:"file" yaml:"file"`
	Stdout							bool							`json:"stdout" yaml:"stdout"`
	Stderr							bool							`json:"stderr" yaml:"stderr"`
	Default							bool							`json:"default" yaml:"default"`
}


