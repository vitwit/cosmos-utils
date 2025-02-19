#/bin/sh

display_usage() {
    printf "** Please check the exported values:: **\n Deamon : $DEAMON\n Denom : $DENOM\n ChainID : $CHAINID\n Daemon home : $DAEMON_HOME\n"
    exit 1
}

if [ -z $DAEMON ] || [ -z $DENOM ] || [ -z $CHAINID ] || [ -z $DAEMON_HOME ]
then 
    display_usage
fi

echo

FROM=$1
if [ -z $FROM ]
then
    FROM=1
fi

TO=$2
if [ -z $TO ]
then
    TO=2
fi


IP="$(dig +short myip.opendns.com @resolver1.opendns.com)"
echo "Public IP address: ${IP}"
PORT=16657

if [ -z $IP ]
then
    IP=127.0.0.1
    PORT=16657
fi

RPC="http://${IP}:${PORT}"

acc1=$($DAEMON keys show account$FROM -a --home $DAEMON_HOME-1 --keyring-backend test)

acc2=$($DAEMON keys show account$TO -a --home $DAEMON_HOME-1 --keyring-backend test)


seq1=$("${DAEMON}" q account "${acc1}" --node $RPC --output json)

seq1no=$(echo "${seq1}" | jq -r '.sequence')

echo "Sequence of 1 is "

echo $seq1no

bound=`expr 10000 + $seq1no`

echo "bound is "

echo $bound

for (( a=$seq1no; a<$bound; a++ ))
do

	    sTx=$("${DAEMON}" tx bank send "${acc1}" "${acc2}" 1000000"${DENOM}" --chain-id "${CHAINID}" --keyring-backend test --home $DAEMON_HOME-1 --node $RPC --output json -y --sequence $a)
   
    sTxHash=$(echo "${sTx}" | jq -r '.txhash')
    echo "** TX HASH :: $sTxHash **"
    
done
