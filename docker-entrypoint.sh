#!/bin/bash
set -e

echo "docker entrypoint running..."

if [ "$1" = 'full_init' ]; then
  MONIKER=first_node
  echo "setting up famd...: " "${MONIKER}"
  # Initialize configuration files and genesis file
  # moniker is the name of your node
  famd init moniker --chain-id fantom-chain

  # Copy the `Address` output here and save it for later use
  # [optional] add "--ledger" at the end to use a Ledger Nano S
  echo jack12345 | famcli keys add jack

  # Copy the `Address` output here and save it for later use
  echo alice12345 | famcli keys add alice

  # Add both accounts, with coins to the genesis file
  famd add-genesis-account "$(famcli keys show jack -a)" 1000fantomtoken,100000000stake
  famd add-genesis-account "$(famcli keys show alice -a)" 1000fantomtoken,100000000stake

  # Configure your CLI to eliminate need for chain-id flag
  famcli config chain-id namechain
  famcli config output json
  famcli config indent true
  famcli config trust-node true

  echo jack12345 | famd gentx --name jack

  # Input gentx into genesis file and validate
  famd collect-gentxs
  famd validate-genesis

  echo "genesis.json:"
  cat /root/.famd/config/genesis.json

  famd start --trace
elif [ "$1" = 'init' ]; then
  MONIKER=node-"$2"
  FIRST_NODE_ID="$3"
  GENESIS="$4"

  echo "setting up simple famd...: " "${MONIKER}"
  famd init "${MONIKER}" --chain-id fantom-chain

  # in ~/.famd/config/config.toml update persistent_peers = "first_id@first_node_ip:26656""
  sed -i "s/persistent_peers =.*/persistent_peers = \"${FIRST_NODE_ID}\"/g" /root/.famd/config/config.toml

  # overwrite with first nodes ~/.famd/config/genesis.json
  echo "${GENESIS}" > /root/.famd/config/genesis.json

  echo "genesis.json:"
  cat /root/.famd/config/genesis.json
  famd validate-genesis

  famd start --trace
else
  echo "starting cmd...: " "$@"
  exec "$@"
fi
