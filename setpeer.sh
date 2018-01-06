
#!/bin/bash
export ORDERER_CA=/opt/ws/crypto-config/ordererOrganizations/wb.gov.in/msp/tlscacerts/tlsca.wb.gov.in-cert.pem

if [ $# -lt 2 ];then
	echo "Usage : . setpeer.sh WBHealthDept| <peerid>"
fi
export peerId=$2

if [[ $1 = "WBHealthDept" ]];then
	echo "Setting to organization WBHealthDept peer "$peerId
	export CORE_PEER_ADDRESS=$peerId.wbhealth.gov.in:7051
	export CORE_PEER_LOCALMSPID=WBHealthDeptMSP
	export CORE_PEER_TLS_CERT_FILE=/opt/ws/crypto-config/peerOrganizations/wbhealth.gov.in/peers/$peerId.wbhealth.gov.in/tls/server.crt
	export CORE_PEER_TLS_KEY_FILE=/opt/ws/crypto-config/peerOrganizations/wbhealth.gov.in/peers/$peerId.wbhealth.gov.in/tls/server.key
	export CORE_PEER_TLS_ROOTCERT_FILE=/opt/ws/crypto-config/peerOrganizations/wbhealth.gov.in/peers/$peerId.wbhealth.gov.in/tls/ca.crt
	export CORE_PEER_MSPCONFIGPATH=/opt/ws/crypto-config/peerOrganizations/wbhealth.gov.in/users/Admin@wbhealth.gov.in/msp
fi

	