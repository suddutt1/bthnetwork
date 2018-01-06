#!/bin/bash
. setpeer.sh WBHealthDept peer0 
export CHANNEL_NAME="wbhealthchannel"
peer chaincode install -n hcdm -v 1.0 -p github.com/hcdm
peer chaincode instantiate -o orderer.wb.gov.in:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C wbhealthchannel -n hcdm -v 1.0 -c '{"Args":["init",""]}' -P " OR( 'WBHealthDeptMSP.member' ) " 
