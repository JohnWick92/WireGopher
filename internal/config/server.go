package config

type HTTP struct {
	Port string `yaml:"port" validate:"required"`
}

type Server struct {
	HTTP *HTTP `yaml:"http" validate:"required"`
}
