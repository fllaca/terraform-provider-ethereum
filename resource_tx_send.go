package main

import (
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/fllaca/terraform-provider-ethereum/ethereum"
)

func resourceTransactionSend() *schema.Resource {
	return &schema.Resource{
		Create: resourceTransactionSendCreate,
		Read:   resourceTransactionSendRead,
		Update: resourceTransactionSendUpdate,
		Delete: resourceTransactionSendDelete,

		Schema: map[string]*schema.Schema{
			"abi": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"to": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"method": {
				Type:     schema.TypeString,
				Optional: true,
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
			"hash": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceTransactionSendCreate(d *schema.ResourceData, m interface{}) error {
	cfg := m.(ethereumConfig)
	ethClient := cfg.client
	contractAbi := d.Get("abi").(string)
	to := d.Get("to").(string)
	method := d.Get("method").(string)
	parameters := d.Get("parameters").([]interface{})

	auth := ethereum.NewAuth(ethClient, cfg.account_key)
	backend := interface{}(ethClient).(bind.ContractBackend)
	tx, err := ethereum.TransactContract(auth, backend, contractAbi, to, method, parameters...)
	if err != nil {
		return err
	}

	d.SetId(tx.Hash().Hex())
	d.Set("hash", tx.Hash().Hex())

	return resourceTransactionSendRead(d, m)
}

func resourceTransactionSendRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceTransactionSendUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceTransactionSendRead(d, m)
}

func resourceTransactionSendDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
