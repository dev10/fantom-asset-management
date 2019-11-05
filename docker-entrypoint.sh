#!/bin/bash
set -e

echo "docker entrypoint running..."

if [ "$1" = 'full_init' ]; then
  MONIKER=first_node
  echo "setting up famd...: " "${MONIKER}"
  # Initialize configuration files and genesis file
  # moniker is the name of your node
  famd init moniker --chain-id fantomchain

  # Copy the `Address` output here and save it for later use
  # [optional] add "--ledger" at the end to use a Ledger Nano S
  echo jack12345 | famcli keys add jack

  # Copy the `Address` output here and save it for later use
  echo alice12345 | famcli keys add alice

  # Add both accounts, with coins to the genesis file
  famd add-genesis-account "$(famcli keys show jack -a)" 1000fantomtoken,100000000stake
  famd add-genesis-account "$(famcli keys show alice -a)" 1000fantomtoken,100000000stake

  # Configure your CLI to eliminate need for chain-id flag
  famcli config chain-id fantomchain
  famcli config output json
  famcli config indent true
  famcli config trust-node true

  echo jack12345 | famd gentx --name jack

  # Input gentx into genesis file and validate
  famd collect-gentxs
  famd validate-genesis

  echo "genesis.json:"
  cat /root/.famd/config/genesis.json
elif [ "$1" = 'init' ]; then
  # get minified genesis: jq -c . < ~/.famd/config/genesis.json > genesis_min.json
  # examples:
  # docker run --rm fam init mon2 33c93e57636bbabad14c9803cffe1ce452ade3d5@127.0.0.1:26656 $(cat genesis_min.json)
  # docker-compose run --no-deps app-2 init mon2 d0afa724f130f4d55164e798ef0b6c87591f1ea0@172.18.0.2:26656 "$(jq -c . < ./app1/.famd/config/genesis.json)"
  # docker-compose run --no-deps app-3 init mon3 d0afa724f130f4d55164e798ef0b6c87591f1ea0@172.18.0.2:26656 "$(cat ./app1/.famd/config/genesis.json)"
  MONIKER=node-"$2"
  FIRST_NODE_ID="$3"
  GENESIS="$4"

  echo "setting up simple famd...: " "${MONIKER}"
  famd init "${MONIKER}" --chain-id fantomchain

  # in ~/.famd/config/config.toml update persistent_peers = "first_id@first_node_ip:26656""
  sed -i "s/persistent_peers =.*/persistent_peers = \"${FIRST_NODE_ID}\"/g" /root/.famd/config/config.toml
  sed -i "s/addr_book_strict =.*/addr_book_strict = \"false\"/g" /root/.famd/config/config.toml

  # overwrite with first nodes ~/.famd/config/genesis.json
  echo "${GENESIS}" > /root/.famd/config/genesis.json

  echo "genesis.json:"
  cat /root/.famd/config/genesis.json

  echo "validating genesis..."
  famd validate-genesis
else
  echo "starting cmd...: " "$@"
  exec "$@"
fi
