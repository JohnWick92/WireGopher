package config

type Firewall struct {
	OPNsense *OPNsense `yaml:"opnsense"`
}

type OPNsense struct {
	URL       string `yaml:"url"`
	APIkey    string `yaml:"apiKey"`
	APIsecret string `yaml:"apiSecret"`
}
