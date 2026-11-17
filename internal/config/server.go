package config

type HTTP struct {
	Port string `yaml:"port"`
}

type Server struct {
	HTTP HTTP `yaml:"http"`
}
