#!/bin/bash
if [[ ! -z "$1" ]]; then  
	. setpeer.sh WBHealthDept peer0 
export CHANNEL_NAME="wbhealthchannel"
	peer chaincode install -n hcdm -v $1 -p github.com/hcdm
	peer chaincode upgrade -o orderer.wb.gov.in:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C wbhealthchannel -n hcdm -v $1 -c '{"Args":["init",""]}' -P " OR( 'WBHealthDeptMSP.member' ) " 
else
	echo ". hcdm_updchain.sh  <Version Number>" 
fi
