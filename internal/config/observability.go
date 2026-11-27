package config

type Observability struct {
	Alloy *Alloy `yaml:"alloy" validate:"required"`
}

type Alloy struct {
	Host *string `yaml:"host" validate:"required"`
	Port *string `yaml:"port" validate:"required"`
}
