package main

import (
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/fllaca/terraform-provider-ethereum/ethereum"
)

func resourceSmartContract() *schema.Resource {
	return &schema.Resource{
		Create: resourceSmartContractCreate,
		Read:   resourceSmartContractRead,
		Update: resourceSmartContractUpdate,
		Delete: resourceSmartContractDelete,

		Schema: map[string]*schema.Schema{
			"abi": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"bin": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"transaction": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"parameters": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"account_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceSmartContractCreate(d *schema.ResourceData, m interface{}) error {
	cfg := m.(ethereumConfig)
	ethClient := cfg.client
	contractAbi := d.Get("abi").(string)
	contractBin := d.Get("bin").(string)
	parameters := d.Get("parameters").([]interface{})

	auth := ethereum.NewAuth(ethClient, cfg.account_key)
	contractBackend := interface{}(ethClient).(bind.ContractBackend)
	address, _, tx, err := ethereum.DeployContract(auth, contractBackend, contractAbi, contractBin, parameters...)
	if err != nil {
		return err
	}

	d.SetId(address.Hex())
	d.Set("address", address.Hex())
	d.Set("transaction", tx.Hash().Hex())

	return resourceSmartContractRead(d, m)
}

func resourceSmartContractRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceSmartContractUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceSmartContractRead(d, m)
}

func resourceSmartContractDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
