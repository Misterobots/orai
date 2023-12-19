#!/bin/bash
set -ux

# always returns true so set -e doesn't exit if it is not running.
killall oraid || true
rm -rf $HOME/.oraid/
killall screen

# make four orai directories
mkdir $HOME/.oraid
end_value=${end_value:-3}
for i in $(seq 1 "$end_value")
do
    mkdir $HOME/.oraid/validator$i

    # init all three validators
    oraid init --chain-id=testing validator$i --home=$HOME/.oraid/validator$i

    # create keys for all three validators
    oraid keys add validator$i --keyring-backend=test --home=$HOME/.oraid/validator$i

    update_genesis () {    
        cat $HOME/.oraid/validator1/config/genesis.json | jq "$1" > $HOME/.oraid/validator1/config/tmp_genesis.json && mv $HOME/.oraid/validator1/config/tmp_genesis.json $HOME/.oraid/validator1/config/genesis.json
    }

    # change staking denom to orai
    update_genesis '.app_state["staking"]["params"]["bond_denom"]="orai"'

    # create validator node 1
    if [ "$i" -eq "1" ]; then
        oraid add-genesis-account $(oraid keys show validator$i -a --keyring-backend=test --home=$HOME/.oraid/validator$i) 1000000000000orai,1000000000000stake --home=$HOME/.oraid/validator$i
        oraid gentx validator$i 500000000orai --keyring-backend=test --home=$HOME/.oraid/validator$i --chain-id=testing
        oraid collect-gentxs --home=$HOME/.oraid/validator$i
        oraid validate-genesis --home=$HOME/.oraid/validator$i
    fi

    # update staking genesis
    update_genesis '.app_state["staking"]["params"]["unbonding_time"]="240s"'
    # update crisis variable to orai
    update_genesis '.app_state["crisis"]["constant_fee"]["denom"]="orai"'
    # udpate gov genesis
    update_genesis '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="orai"'
    # update mint genesis
    update_genesis '.app_state["mint"]["params"]["mint_denom"]="orai"'
    update_genesis '.app_state["gov"]["voting_params"]["voting_period"]="30s"'
    # port key (validator1 uses default ports)
    # validator1 1317, 9090, 9091, 26658, 26657, 26656, 6060
    # validator2 1316, 9088, 9089, 26655, 26654, 26653, 6061
    # validator3 1315, 9086, 9087, 26652, 26651, 26650, 6062


    # change app.toml values
    VALIDATOR_APP_TOML=$HOME/.oraid/validator$i/config/app.toml

    # change config.toml values
    VALIDATOR_CONFIG=$HOME/.oraid/validator$i/config/config.toml

    sed -i -E 's|tcp://0.0.0.0:1317|tcp://0.0.0.0:131'$i'|g' $VALIDATOR_APP_TOML
    sed -i -E 's|0.0.0.0:9090|0.0.0.0:809'$i'|g' $VALIDATOR_APP_TOML
    sed -i -E 's|0.0.0.0:9091|0.0.0.0:709'$i'|g' $VALIDATOR_APP_TOML

    # Pruning - comment this configuration if you want to run upgrade script
    pruning="custom"
    pruning_keep_recent="5"
    pruning_keep_every="10"
    pruning_interval="10000"

    sed -i -e "s%^pruning *=.*%pruning = \"$pruning\"%; " $VALIDATOR_APP_TOML
    sed -i -e "s%^pruning-keep-recent *=.*%pruning-keep-recent = \"$pruning_keep_recent\"%; " $VALIDATOR_APP_TOML
    sed -i -e "s%^pruning-keep-every *=.*%pruning-keep-every = \"$pruning_keep_every\"%; " $VALIDATOR_APP_TOML
    sed -i -e "s%^pruning-interval *=.*%pruning-interval = \"$pruning_interval\"%; " $VALIDATOR_APP_TOML

    # state sync  - comment this configuration if you want to run upgrade script
    snapshot_interval="10"
    snapshot_keep_recent="2"

    sed -i -e "s%^snapshot-interval *=.*%snapshot-interval = \"$snapshot_interval\"%; " $VALIDATOR_APP_TOML
    sed -i -e "s%^snapshot-keep-recent *=.*%snapshot-keep-recent = \"$snapshot_keep_recent\"%; " $VALIDATOR_APP_TOML

    sed -i -E 's|allow_duplicate_ip = false|allow_duplicate_ip = true|g' $VALIDATOR_CONFIG

    # validator2
    sed -i -E 's|tcp://127.0.0.1:26658|tcp://0.0.0.0:2765'$i'|g' $VALIDATOR_CONFIG
    sed -i -E 's|tcp://127.0.0.1:26657|tcp://0.0.0.0:2665'$i'|g' $VALIDATOR_CONFIG
    sed -i -E 's|tcp://0.0.0.0:26656|tcp://0.0.0.0:2565'$i'|g' $VALIDATOR_CONFIG
    sed -i -E 's|allow_duplicate_ip = false|allow_duplicate_ip = true|g' $VALIDATOR_CONFIG

    if [ "$i" -ne "1" ]; then
        cp $HOME/.oraid/validator1/config/genesis.json $HOME/.oraid/validator$i/config/genesis.json
        # copy tendermint node id of validator1 to persistent peers of validator2-3
        sed -i -E "s|persistent_peers = \"\"|persistent_peers = \"$(oraid tendermint show-node-id --home=$HOME/.oraid/validator1)@localhost:25651\"|g" $VALIDATOR_CONFIG
        prev_i=$((i - 1))
        sed -i -E "s|persistent_peers = \"\"|persistent_peers = \"$(oraid tendermint show-node-id --home=$HOME/.oraid/validator1)@localhost:2565$prev_i\"|g" $VALIDATOR_CONFIG
    fi
done

for i in $(seq 1 "$end_value")
do
    # start all validators
    screen -S validator$i -d -m oraid start --home=$HOME/.oraid/validator$i --minimum-gas-prices=0.00001orai
done

# send orai from first validator to second validator
echo "Waiting 10 seconds to send funds to validators"
sleep 10

for i in $(seq 1 "$end_value")
do
    oraid tx send $(oraid keys show validator1 -a --keyring-backend=test --home=$HOME/.oraid/validator1) $(oraid keys show validator$i -a --keyring-backend=test --home=$HOME/.oraid/validator$i) 5000000000orai --keyring-backend=test --home=$HOME/.oraid/validator1 --chain-id=testing --broadcast-mode block --gas 200000 --fees 2orai --node http://localhost:26651 --yes

    if [ "$i" -ne "1" ]; then
        oraid tx staking create-validator --amount=500000000orai --from=validator$i --pubkey=$(oraid tendermint show-validator --home=$HOME/.oraid/validator$i) --moniker="validator$i" --chain-id="testing" --commission-rate="0.1" --commission-max-rate="0.2" --commission-max-change-rate="0.05" --min-self-delegation="500000000" --keyring-backend=test --home=$HOME/.oraid/validator$i --broadcast-mode block --gas 200000 --fees 2orai --node http://localhost:26651 --yes
    fi
done

total_validators=$(oraid query staking validators --home ~/.oraid/validator1 --node http://localhost:26651 --output json | jq '.validators | length')

if [ "$total_validators" -eq "$end_value" ]; then
    echo "All Validators are up and running!"
else
    echo "Error setting up multi-node local testnet with $end_value validators"
fi