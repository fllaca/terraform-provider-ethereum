package main

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
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
				Optional: true,
			},
			"transaction": &schema.Schema{
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
		},
	}
}

func resourceSmartContractCreate(d *schema.ResourceData, m interface{}) error {
	cfg := m.(ethereumConfig)
	ethClient := cfg.client
	contractAbi := d.Get("abi").(string)
	contractBin := d.Get("bin").(string)
	parameters := d.Get("parameters").([]interface{})
	auth := buildAuth(ethClient, cfg.account_key)
	contractBackend := interface{}(ethClient).(bind.ContractBackend)
	address, _, tx, err := deployContract(auth, contractBackend, contractAbi, contractBin, parameters...)
	if err != nil {
		fmt.Print(err)
		return err
	}
	d.SetId(address.Hex())
	d.Set("address", address.Hex())
	fmt.Print(tx)
	//	d.Set("transaction", tx.Hash().Hex())

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

// DeployStore deploys a new Ethereum contract, binding an instance of Store to it.
func deployContract(auth *bind.TransactOpts, backend bind.ContractBackend, contractAbi string, contractBin string, params ...interface{}) (common.Address, *bind.BoundContract, *types.Transaction, error) {

	parsed, err := abi.JSON(strings.NewReader(contractAbi))
	if err != nil {
		return common.Address{}, nil, nil, errors.New("Error Parsing: " + err.Error())
	}

	inputs := make([]interface{}, 0)
	for i, param := range params {
		switch parsed.Constructor.Inputs[i].Type.String() {
		case "uint256":
			v, _ := strconv.Atoi(param.(string))
			bigInt := big.NewInt(int64(v))
			inputs = append(inputs, &bigInt)
		case "uint8":
			v, _ := strconv.Atoi(param.(string))
			inputs = append(inputs, uint8(v))
		// strings fall here
		default:
			inputs = append(inputs, param)
		}
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(contractBin), backend, inputs...)

	if err != nil {
		//return common.Address{}, contract, nil, err
		return common.Address{}, contract, nil, errors.New("Error Deploying: " + err.Error())
	}
	return address, contract, tx, nil
}

func buildAuth(ethClient *ethclient.Client, key string) *bind.TransactOpts {
	privateKey, err := crypto.HexToECDSA(key)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := ethClient.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := ethClient.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(1000000) // in units
	auth.GasPrice = gasPrice
	return auth
}
