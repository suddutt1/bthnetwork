
#!/bin/bash -e




echo "Building channel for wbhealthchannel" 

. setpeer.sh WBHealthDept peer0
export CHANNEL_NAME="wbhealthchannel"
peer channel create -o orderer.wb.gov.in:7050 -c $CHANNEL_NAME -f ./wbhealthchannel.tx --tls true --cafile $ORDERER_CA -t 10000


. setpeer.sh WBHealthDept peer0
export CHANNEL_NAME="wbhealthchannel"
peer channel join -b $CHANNEL_NAME.block

