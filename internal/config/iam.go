package config

type IAM struct {
	Provider *string   `yaml:"provider" validate:"required"`
	Keycloak *Keycloak `yaml:"keycloak" validate:"required,omitempty"`
}

type Keycloak struct {
	URL          *string `yaml:"url" validate:"required"`
	Issuer       *string `yaml:"issuer" validate:"required"`
	Realm        *string `yaml:"realm" validate:"required"`
	ClientID     *string `yaml:"clientID" validate:"required"`
	ClientSecret *string `yaml:"clientSecret" validate:"required"`
	RedirectURI  *string `yaml:"redirectURI" validate:"required"`
}
