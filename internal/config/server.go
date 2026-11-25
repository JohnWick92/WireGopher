package config

type HTTP struct {
	Port string `yaml:"port" validate:"required"`
	URL  string `yaml:"url" validate:"required"`
}

type Server struct {
	HTTP *HTTP `yaml:"http" validate:"required"`
}
