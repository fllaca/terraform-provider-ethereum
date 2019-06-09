package main

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hashicorp/terraform/helper/schema"
)

type ethereumConfig struct {
	client      *ethclient.Client
	account_key string
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"client_address": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ETHEREUM_CLIENT_ADDRESS", ""),
			},
			"account_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ETHEREUM_ACCOUNT_KEY", ""),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"ethereum_smart_contract": resourceSmartContract(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	client_address := d.Get("client_address").(string)
	account_key := d.Get("account_key").(string)
	ethClient, err := ethclient.Dial(client_address)

	cfg := ethereumConfig{
		client:      ethClient,
		account_key: account_key,
	}

	return cfg, err
}
