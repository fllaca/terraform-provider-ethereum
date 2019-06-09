provider "ethereum" {
  client_address = "http://localhost:8545"
  account_key    = "a7ba4cf92420bb9694130ab9bd215a7acdb96cb837ed70a5bf38e31c938d6c29"
}

resource "ethereum_smart_contract" "store" {
  abi = "${file("Store.abi")}"
  bin = "${file("Store.bin")}"

  parameters = ["v1.0"]
}

resource "ethereum_smart_contract_transact" "setItem" {
  abi = "${ethereum_smart_contract.store.abi}"
  to = "${ethereum_smart_contract.store.address}"

	method = "setItem"

  parameters = ["key", "value"]
}
