package opnsense

import (
	"github.com/JohnWick92/WireGopher/internal/config"
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/opnsense"
)

var client opnsense.Client

func getClient(cfg *config.OPNsense) {
	apiClient := api.NewClient(api.Options{
		Uri:           *cfg.URL,
		APIKey:        *cfg.APIkey,
		APISecret:     *cfg.APIsecret,
		AllowInsecure: true,
	})
	client = opnsense.NewClient(apiClient)
}

func Run(cfg *config.OPNsense) {
	getClient(cfg)
}
