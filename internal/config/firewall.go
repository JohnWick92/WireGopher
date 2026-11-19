package config

type Firewall struct {
	FirewallModel *string    `yaml:"firewallModel" validate:"required"`
	OPNsense      *OPNsense `yaml:"opnsense,omitempty"`
}

type OPNsense struct {
	URL       *string `yaml:"url" validate:"required"`
	APIkey    *string `yaml:"apiKey" validate:"required"`
	APIsecret *string `yaml:"apiSecret" validate:"required"`
}
