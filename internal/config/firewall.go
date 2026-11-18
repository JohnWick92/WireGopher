package config

type Firewall struct {
	FirewallModel string `yaml:"firewallModel"`
	OPNsense *OPNsense `yaml:"opnsense,omitempty"`
}

type OPNsense struct {
	URL       string `yaml:"url"`
	APIkey    string `yaml:"apiKey"`
	APIsecret string `yaml:"apiSecret"`
}
