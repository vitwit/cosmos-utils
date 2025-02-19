#/bin/sh

NODES=$1
if [ -z $NODES ]
then
    NODES=1
fi

echo "**** Number of nodes to be paused: $NODES ****"

echo

echo "---------- Stopping systemd service files --------"

for (( a=1; a<=$NODES; a++ ))
do
    sudo -S systemctl stop $DAEMON-${a}.service

    echo "-- Executed sudo -S systemctl stop $DAEMON-${a}.service --"
done

echo