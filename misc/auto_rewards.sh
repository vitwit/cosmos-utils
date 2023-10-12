#!/bin/bash

# Read configuration from JSON file
config_file="config.json"
networks=($(jq -c '.[]' "$config_file"))

# Loop through the networks
while true
do
    for network in "${networks[@]}"
    do
        # Extract parameters from JSON
        binary=$(echo "$network" | jq -r '.binary')
        granter=$(echo "$network" | jq -r '.granter')
        grantee=$(echo "$network" | jq -r '.grantee')
        chain_id=$(echo "$network" | jq -r '.chainId')
        node=$(echo "$network" | jq -r '.node')
        feepayer=$(echo "$network" | jq -r '.feepayer')

        echo "About to withdraw commission and reward for network: $chain_id"

        # Withdraw rewards and execute on the network
        $binary tx distribution withdraw-rewards $validator --from $granter --chain-id $chain_id -y --generate-only > rewards.json
        $binary tx authz exec rewards.json --chain-id $chain_id --node $node --fees 200uatom --fee-account $feepayer --from $grantee -y

        # Add sleep if needed before processing the next network
        # sleep 60
    done
done
