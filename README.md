# Ethereum Terraform Provider

Deploys smart contracts and executes transactions on Ethereum


## Usage

```
provider "ethereum" {
  client_address = "http://localhost:8545"
  account_key = "<your-account-key>"
}

resource "ethereum_smart_contract" "store" {
  abi = "${file("Store.abi")}"
  bin = "${file("Store.bin")}"
	
  parameters = ["v1.0"]
}

```
