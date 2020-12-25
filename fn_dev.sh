#!/bin/bash

#
# Copyright Oraichain All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#
# export so other script can access

# colors
BROWN='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m'
BOLD=$(tput bold)
NORMAL=$(tput sgr0)

# environment
BASE_DIR=$PWD
SCRIPT_NAME=`basename "$0"`

# verify the result of the end-to-end test
verifyResult() {  
  if [ $1 -ne 0 ] ; then
    echo "!!!!!!!!!!!!!!! "$2" !!!!!!!!!!!!!!!!"
    echo "========= ERROR !!! FAILED to execute End-2-End Scenario ==========="
    echo
      exit 1
  fi
}

printCommand(){
  echo -e ""
  printBoldColor $BROWN "Command:"
  printBoldColor $BLUE "\t$1"  
}

printBoldColor(){
  echo -e "$1${BOLD}$2${NC}${NORMAL}"
}

# Print the usage message
printHelp () {

  echo $BOLD "Usage: "  
  echo "  $SCRIPT_NAME -h|--help (Show help)"  
  echo $NORMAL

  if [[ ! -z $2 ]]; then
    res=$(printHelp 0 | grep -A2 "\- '$2' \-")
    echo "$res"    
  else      
    printBoldColor $BROWN "      - 'start' - Run the full node"
    printBoldColor $BLUE  "          fn start"           
    echo
    printBoldColor $BROWN "      - 'broadcast' - broadcast transaction"
    printBoldColor $BLUE  "          fn broadcast --key value"           
    echo
    printBoldColor $BROWN "      - 'init' - Init the orai node"
    printBoldColor $BLUE  "          fn init"           
    echo
    printBoldColor $BROWN "      - 'sign' - sign transaction"
    printBoldColor $BLUE  "          fn sign --key value"           
    echo
    printBoldColor $BROWN "      - 'initScript' - init AI request script"
    printBoldColor $BLUE  "          fn initScript --key value"           
    echo
    printBoldColor $BROWN "      - 'clear' - Clear all existing data"
    printBoldColor $BLUE  "          fn clear"           
    echo
  fi

  echo
  echo "  $SCRIPT_NAME method --argument=value"
  
  # default exit as 0
  exit ${1:-0}
}


# Get a value:
getArgument() {     
  local key="args_${1/-/_}"  
  echo ${!key:-$2}  
}


# check first param is method
if [[ $1 =~ ^[a-z] ]]; then 
  METHOD=$1
  shift
fi

# use [[ ]] we dont have to quote string
args=()
case "$METHOD" in
  bash)        
    while [[ ! -z $2 ]];do         
      if [[ ${1:0:2} == '--' ]]; then
        KEY=${1/--/}            
        if [[ $KEY =~ ^([a-zA-Z_-]+)=(.+) ]]; then                
          declare "args_${BASH_REMATCH[1]/-/_}=${BASH_REMATCH[2]}"
        else          
          declare "args_${KEY/-/_}=$2" 
          shift
        fi    
      else 
        args+=($1)
      fi
      shift
    done
    QUERY="$@"            
  ;;
  config)        
    while [[ $# -gt 0 ]] ; do            
      if [[ ${1:0:2} == '--' ]]; then
        KEY=${1/--/}            
        if [[ $KEY =~ ^([a-zA-Z_-]+)=(.+) ]]; then                
          declare "args_${BASH_REMATCH[1]/-/_}=${BASH_REMATCH[2]}"
        else          
          declare "args_${KEY/-/_}=$2" 
          shift
        fi    
      else 
        args+=($1)
      fi
      shift
    done     
  ;;
  *) 
    # normal processing
    while [[ $# -gt 0 ]] ; do                
      if [[ ${1:0:2} == '--' ]]; then
        KEY=${1/--/}                
        if [[ $KEY =~ ^([a-zA-Z_-]+)=(.+) ]]; then         
          declare "args_${BASH_REMATCH[1]/-/_}=${BASH_REMATCH[2]}"
        else
          declare "args_${KEY/-/_}=$2"        
          shift
        fi    
      else 
        case "$1" in
          -h|\?)            
            printHelp 0 $2
          ;;
          *)  
            args+=($1)            
          ;;  
        esac    
      fi 
      shift
    done 
  ;; 
esac


clear(){
    rm -rf .oraid/
    rm -rf .oraicli/
    rm -rf .oraifiles/    
}

oraidFn(){
    # oraid start
    orai start --chain-id $CHAIN_ID --laddr tcp://0.0.0.0:1317 --node tcp://0.0.0.0:26657 # --trust-node
}


initFn(){

    sleep 8

    ./init.sh $CHAIN_ID $USER
    # run at background without websocket
    oraid start --minimum-gas-prices 0.025orai &
    # 30 seconds timeout
    timeout 30 bash -c 'while [[ "$(curl -s -o /dev/null -w ''%{http_code}'' localhost:26657/health)" != "200" ]]; do sleep 1; done' || false
    local reporter="${USER}_reporter"
    ./websocket.sh $USER $reporter #$reporter
    sleep 10
    pkill oraid
}

initDevFn(){

    make all

    ./init.sh $CHAIN_ID $USER
    # run at background without websocket
    oraid start --minimum-gas-prices 0.025orai &
    # 30 seconds timeout
    timeout 30 bash -c 'while [[ "$(curl -s -o /dev/null -w ''%{http_code}'' localhost:26657/health)" != "200" ]]; do sleep 1; done' || false
    local reporter="${USER}_reporter"
    ./websocket.sh $USER $reporter #$reporter
    sleep 10
    pkill oraid
}


initScriptFn(){
  echo "y" | oraicli tx provider set-datasource coingecko_eth ./testfiles/coingecko_eth.py "A data source that fetches the ETH price from Coingecko API" "4000orai" --from $USER --fees 5000orai

  sleep 5

  echo "y" | oraicli tx provider set-datasource crypto_compare_eth ./testfiles/crypto_compare_eth.py "A data source that collects ETH price from crypto compare" "5000orai" --from $USER --fees 5000orai

  sleep 5

  echo "y" | oraicli tx provider set-testcase testcase_price ./testfiles/testcase_price.py "A sample test case that uses the expected output of users provided to verify the bitcoin price from the datasource" "5000orai" --from $USER --fees 5000orai

  sleep 5

  echo "y" | oraicli tx provider set-oscript oscript_eth ./testfiles/oscript_eth.py "An oracle script that fetches and aggregates ETH price from different sources" --ds coingecko_eth,crypto_compare_eth --tc testcase_price --from $USER --fees 5000orai
}

unsignedFn(){
  local id=$(curl -s "http://localhost:1317/auth/accounts/$(oraicli keys show $USER -a)" | jq ".result.value.address" -r)
  local unsigned=$(curl --location --request POST 'http://localhost:1317/airequest/aireq' \
--header 'Content-Type: application/json' \
--data-raw '{
    "base_req":{
        "from":"'$id'",
        "chain_id":"'$CHAIN_ID'"
    },
    "oracle_script_name":"oscript_eth",
    "input":"",
    "expected_output":{"price":"5000"},
    "fees":"60000orai",
    "validator_count": "1"
}' > tmp/unsignedTx.json)

    res=$?  
    verifyResult $res "Unsigned failed"
}

unsignedSetDsFn(){
  local id=$(curl -s "http://localhost:1317/auth/accounts/$(oraicli keys show $USER -a)" | jq ".result.value.address" -r)
  local unsigned=$(curl --location --request POST 'http://localhost:1317/provider/datasource' \
--header 'Content-Type: application/json' \
--data-raw '{
    "base_req":{
        "from":"'$id'",
        "chain_id":"Oraichain"
    },
    "name":"coingecko_eth",
    "code_path":"/workspace/testfiles/coingecko_eth.py",
    "description":"NTAwMA==",
    "fees":"60000orai",
    "test":["abc","efgh"]
}' > tmp/unsignedTx.json)
}

signFn(){     
    # $1 is account number
    local sequence=$(curl -s "http://localhost:1317/auth/accounts/$(oraicli keys show $USER -a)" | jq ".result.value.sequence" -r)
    local acc_num=$(curl -s "http://localhost:1317/auth/accounts/$(oraicli keys show $USER -a)" | jq ".result.value.account_number" -r)
    oraicli tx sign tmp/unsignedTx.json --from $USER --offline --chain-id $CHAIN_ID --sequence $sequence --account-number $acc_num > tmp/signedTx.json
    oraicli tx broadcast tmp/signedTx.json

    res=$?  
    verifyResult $res "Signed failed"
}

createValidatorFn() {
  local amount=$(getArgument "amount" 10000orai)
  local pubkey=$(getArgument "pubkey" oraivalconspub1addwnpepqvydmv22mkzc9rc92g43unew08cmj4q46dhk7vz0a9fj2xjsjn2lvqj0dfr)
  local moniker=$(getArgument "moniker" ducphamle)
  local commissionRate=$(getArgument "commission-rate" 0.10)
  local commissionMaxRate==$(getArgument "commission-max-rate" 0.20)
  local commissionMaxChangeRate==$(getArgument "commission-max-change-rate" 0.01)
  oraicli tx staking create-validator --amount $amount --pubkey $pubkey --moniker $moniker --chain-id $CHAIN_ID --commission-rate $commissionRate --commission-max-rate $commissionMaxRate --commission-max-change-rate $commissionMaxChangeRate --min-self-delegation 100 --gas auto --gas-adjustment 1.15 --gas-prices 0.025orai --from $USER
}

USER=$(getArgument "user" duc)
CHAIN_ID=$(getArgument "chain-id" Oraichain)

# processing
case "${METHOD}" in     
  hello)
    helloFn
  ;;
  init)
    initFn
  ;;
  initDev)
    initDevFn
  ;;
  start)
    oraidFn
  ;;  
  unsign)
    unsignedFn
  ;;
  unsignedSetDs)
    unsignedSetDsFn
  ;;
  initScript)
    initScriptFn
  ;;
  sign)
    signFn
  ;;
  broadcast)
    broadcastFn
  ;;  
  createValidator)
    createValidatorFn
  ;;
  *) 
    printHelp 1 ${args[0]}
  ;;
esac