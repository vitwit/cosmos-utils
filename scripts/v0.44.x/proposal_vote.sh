#/bin/sh

display_usage() {
    printf "** Please check the exported value of:: **\n Daemon : $DAEMON\n Chain ID : $CHAINID\n"
    exit 1
}

if [ -z $DAEMON ] || [ -z $CHAINID ]
then 
    display_usage
fi

# read no.of validators to be vote on the proposal
NODES=$1
if [ -z $NODES ]
then
    NODES=2
fi

echo "--------- No.of validators who have to vote on the proposal : $NODES ------------"

IP="$(dig +short myip.opendns.com @resolver1.opendns.com)"
echo "Public IP address: ${IP}"

if [ -z $IP ]
then
    IP=127.0.0.1
fi

echo "--------Get voting period proposals--------------"
vp=$("${DAEMON}" q gov proposals --status voting_period --output json)
len=$(echo "${vp}" | jq -r '.proposals | length' )
echo "** Length of voting period proposals : $len **"
echo

for row in $(echo "${vp}" | jq -r '.proposals | .[] | @base64'); do

  PID=$(echo "${row}" | base64 --decode | jq -r '.proposal_id')
  echo
  echo
  echo "** Checking votes for proposal id : $PID **"

  for (( a=1; a<=$NODES; a++ ))
  do
    DIFF=`expr $a - 1`
    INC=`expr $DIFF \* 2`
    PORT=`expr 16657 + $INC` #get ports
    RPC="http://${IP}:${PORT}"
    echo " **** NODE :: $RPC  ****"

    # get validator address
    validator=$("${DAEMON}" keys show validator${a} --bech val --keyring-backend test --home $DAEMON_HOME-${a} --output json)
    VALADDRESS=$(echo "${validator}" | jq -r '.address')
    FROMKEY="validator${a}"
    VOTER=$VALADDRESS

    echo "** voter address :: $VALADDRESS and from key :: $FROMKEY **"

    # Check vote status
    getVote=$( ("${DAEMON}" q gov vote "${PID}" "${VOTER}" --output json) 2>&1)
   
    if [ "$?" -eq 0 ];  #0 indiactes no reponse with gov vote query so we can proceed further
    then
      voted=$(echo "${getVote}" | jq -r '.option')
      #echo "*** Proposal Id : $PID and VOTER : $VOTER and VOTE OPTION : $voted ***"
      #cast vote
      castVote=$( ("${DAEMON}" tx gov vote "${PID}" yes --from "${FROMKEY}" --fees 1000"${DENOM}" --chain-id "${CHAINID}" --node "${RPC}" --home $DAEMON_HOME-${a} --keyring-backend test --output json -y) 2>&1) 
      
      sleep 6s

      txHash=$(echo "${castVote}"| jq -r '.txhash')
      echo "** TX HASH :: $txHash **"
      # query the txhash and check the code
      txResult=$("${DAEMON}" q tx "${txHash}" --node $RPC --output json)
      checkVote=$(echo "${txResult}"| jq -r '.code')

      if [[ "$checkVote" != "" ]];
      then
        if [ "$checkVote" -eq 0 ];
        then
          echo "**** $FROMKEY successfully voted on the proposal :: (proposal ID : $PID and address $VOTER ) !!  txHash is : $txHash ****"
        else 
          echo "**** $FROMKEY getting error while casting vote for ( Proposl ID : $PID and address $VOTER )!!!!  txHash is : $txhash and REASON : $(echo "${castVote}" | jq '.raw_log') ****"
        fi
      fi
    else
      echo "Error while getting votes of proposal ID : $PID of $FROMKEY address : $VOTER"
    fi
  done
done
