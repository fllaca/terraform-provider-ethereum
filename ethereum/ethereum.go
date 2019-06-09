package ethereum

import (
	"errors"
	"context"
	"crypto/ecdsa"
	"strconv"
	"math/big"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/core/types"
)


func ConvertArgumentsTypes(arguments abi.Arguments, params ...interface{}) []interface{} {
	inputs := make([]interface{}, 0)
	for i, param := range params {
		switch arguments[i].Type.String() {
		case "uint256":
			v, _ := strconv.Atoi(param.(string))
			bigInt := big.NewInt(int64(v))
			inputs = append(inputs, &bigInt)
		case "uint8":
			v, _ := strconv.Atoi(param.(string))
			inputs = append(inputs, uint8(v))
		case "bytes32":
			var arr [32]byte

			copy(arr[:], []byte(param.(string))[:32])
			inputs = append(inputs, arr)
		// strings fall here
		default:
			inputs = append(inputs, param)
		}
	}
	return inputs
}

func NewAuth(ethClient *ethclient.Client, key string) *bind.TransactOpts {
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

func DeployContract(auth *bind.TransactOpts, backend bind.ContractBackend, contractAbi string, contractBin string, params ...interface{}) (common.Address, *bind.BoundContract, *types.Transaction, error) {

	parsed, err := abi.JSON(strings.NewReader(contractAbi))
	if err != nil {
		return common.Address{}, nil, nil, errors.New("Error Parsing: " + err.Error())
	}

	inputs := ConvertArgumentsTypes(parsed.Constructor.Inputs, params...)

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(contractBin), backend, inputs...)

	if err != nil {
		//return common.Address{}, contract, nil, err
		return common.Address{}, contract, nil, errors.New("Error Deploying: " + err.Error())
	}
	return address, contract, tx, nil
}

func BindContract(address common.Address, backend bind.ContractBackend, parsed abi.ABI) (*bind.BoundContract, error) {
	return bind.NewBoundContract(address, parsed, backend, backend, backend), nil
}

func TransactContract(auth *bind.TransactOpts, backend bind.ContractBackend, contractAbi string, contractAddress string, method string, params ...interface{}) (*types.Transaction, error) {
	parsed, err := abi.JSON(strings.NewReader(contractAbi))
	if err != nil {
		return nil, errors.New("Error Parsing: " + err.Error())
	}

	to := common.HexToAddress(contractAddress)

	contract, err := BindContract(to, backend, parsed)
	if err != nil {
		return nil, errors.New("Error Binding Contract: " + err.Error())
	}

	inputs := ConvertArgumentsTypes(parsed.Methods[method].Inputs, params...)

	return contract.Transact(auth, method, inputs...)
}
