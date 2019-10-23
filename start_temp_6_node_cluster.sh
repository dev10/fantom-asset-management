#!/usr/bin/env bash

# this script requires jq to be installed
command -v jq >/dev/null 2>&1 || { echo >&2 "jq is required but it's not installed. Aborting."; exit 1; }

ROOT_DIR=/tmp/
APP_DIR=fam-app
NODE_COUNT=6
export APP_DATA_PATH=${ROOT_DIR}${APP_DIR}
echo 'setting APP_DATA_PATH'

printenv | grep APP_DATA_PATH
docker-compose config

function take_ownership_of_data() {
  echo 'need to access genesis configuration in data directories'
  sudo chown -R "$(whoami)": ${APP_DATA_PATH}-*
}

function initialise() {
  read -rp "Do you wish to remove any existing data first (y/n)?" answer
  case ${answer:0:1} in
    y|Y )
      docker-compose stop;
      docker-compose rm;
      take_ownership_of_data;
      for n in $(seq 1 $NODE_COUNT)
      do
        docker volume rm fantom-asset-management_famd_data_"$n";
        DELETE_PATH=$APP_DATA_PATH-$n
        read -rp "Do you wish to remove the directory: ${DELETE_PATH} (y/n)?" answer
        case ${answer:0:1} in
        y|Y )
          echo 'removing ' "${DELETE_PATH}" && rm -rfv "${DELETE_PATH}";;
        * )
          ;;
        esac
      done
      ;;
    * )
      ;;
  esac

  echo 'creating data directories'
  for n in {1..6}
  do
    mkdir -pv ${APP_DATA_PATH}-$n
  done

  echo 'initialising first node with genesis'
  docker-compose run app-1 full_init

  take_ownership_of_data;

  echo 'getting node address info from genesis.json'
  NODE_ADDRESS="$( jq -r '.app_state[] | select(.gentxs != null)? | { gentxs }[] | .[] | { value }[] | .memo' ${APP_DATA_PATH}-1/.famd/config/genesis.json )"

  echo 'setting up other nodes to connect to first node at: ' "${NODE_ADDRESS}"
  for n in {2..6}
  do
    # docker-compose run --no-deps app-${n} init mon${n} "${NODE_ADDRESS}" "$(jq -c . < /tmp/fam-app-1/.famd/config/genesis.json)"
    docker-compose run -e APP_DATA_PATH --no-deps app-${n} init mon${n} "${NODE_ADDRESS}" "$(cat ${APP_DATA_PATH}-1/.famd/config/genesis.json)"
  done

  take_ownership_of_data;
}

function start_nodes() {
  echo 'starting...'
  docker-compose up
}

read -rp "Do you wish to initialise first before starting (y/n)?" answer
case ${answer:0:1} in
  y|Y )
    initialise; start_nodes;;
  * )
    start_nodes;;
esac