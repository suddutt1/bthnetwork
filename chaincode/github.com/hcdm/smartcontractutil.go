package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var _SC_LOGGER = shim.NewLogger("SmartContract")

type SmartContract struct {
}
type AbstractInput interface{}

func (sc *SmartContract) init(stub shim.ChaincodeStubInterface) pb.Response {
	_SC_LOGGER.Info("Inside init method")

	return shim.Success(nil)
}
func (sc *SmartContract) probe(stub shim.ChaincodeStubInterface) pb.Response {

	_SC_LOGGER.Info("Inside probe method")
	ts := time.Now().Format(time.UnixDate)
	output := "{\"status\":\"Success\",\"ts\" : \"" + ts + "\" }"
	_SC_LOGGER.Info("Retuning " + output)
	return shim.Success([]byte(output))
}

//ValidateObjectIntegrity Validates the input json array
func (sc *SmartContract) ValidateObjectIntegrity(jsonInput string) (bool, string, []interface{}) {

	var errMsgBuf bytes.Buffer
	var inputObject interface{}
	parsedObjects := make([]interface{}, 0)
	json.Unmarshal([]byte(jsonInput), &inputObject)
	switch inputObject.(type) {
	case []interface{}:
		_SC_LOGGER.Info("Array detected")
		allGood := true

		for index, item := range inputObject.([]interface{}) {
			isGood := sc.CheckObjects(item)
			allGood = allGood && isGood
			if !isGood {
				errMsgBuf.WriteString(fmt.Sprintf("\"Object type missing for %d\",", index))
			}
			parsedObjects = append(parsedObjects, item)
		}

		return allGood, errMsgBuf.String(), parsedObjects
	case interface{}:
		_SC_LOGGER.Info("Object detected")
		isGood := sc.CheckObjects(inputObject)
		if !isGood {
			errMsgBuf.WriteString("Object type missing")
		}
		parsedObjects = append(parsedObjects, inputObject)
		return isGood, errMsgBuf.String(), parsedObjects
	default:
		_SC_LOGGER.Info("Unkown data type")
	}

	return false, "Unkown data type", nil
}

//ValidateAndInsertObject validates the non existance of the object and inserts
func (sc *SmartContract) ValidateAndInsertObject(stub shim.ChaincodeStubInterface, input interface{}, idField string) (bool, string) {
	_SC_LOGGER.Info("ValidateAndInsertObject:Start")
	isSuccess := false
	errMsg := ""
	dataMap, mapOk := input.(map[string]interface{})
	if mapOk == true {
		id, idOk := dataMap[idField].(string)
		if idOk == true && id != "" {
			existingRecord, err := stub.GetState(id)
			_SC_LOGGER.Infof("Existing record %s", string(existingRecord))
			if len(existingRecord) == 0 && err == nil {
				json, _ := json.Marshal(dataMap)
				errSave := stub.PutState(id, json)
				if errSave == nil {
					isSuccess = true
					_SC_LOGGER.Info("Save success")
				} else {
					errMsg = "Not able to save the record in hyperledger"
					_SC_LOGGER.Info("Not able to save")
				}
			} else {
				errMsg = "Id already exists"
				_SC_LOGGER.Info("Id already exists")
			}
		} else {
			errMsg = "Id filed is invalid "
		}
	} else {
		errMsg = "Interface is not a map object"
	}

	return isSuccess, errMsg
}

//GetObjectByKey returns data from hyperledger using the key
func (sc *SmartContract) GetObjectByKey(stub shim.ChaincodeStubInterface, id string) interface{} {
	var outputObject interface{}
	recordBytes, err := stub.GetState(id)
	if len(recordBytes) > 0 && err == nil {
		json.Unmarshal(recordBytes, &outputObject)
		return outputObject
	}
	return nil
}

//RetriveRecords based on the selector criteria
func (sc *SmartContract) RetriveRecords(stub shim.ChaincodeStubInterface, criteria string) []map[string]interface{} {
	records := make([]map[string]interface{}, 0)
	selectorString := fmt.Sprintf("{\"selector\":%s }", criteria)
	_SC_LOGGER.Info("Query Selector :" + selectorString)
	resultsIterator, _ := stub.GetQueryResult(selectorString)
	for resultsIterator.HasNext() {
		record := make(map[string]interface{})
		recordBytes, _ := resultsIterator.Next()
		err := json.Unmarshal(recordBytes.Value, &record)
		if err != nil {
			_SC_LOGGER.Infof("Unable to unmarshal data retived:: %v", err)
		}
		records = append(records, record)
	}
	return records
}

//CheckObjects checks only the objType attribute. TODO more validations
func (sc *SmartContract) CheckObjects(input interface{}) bool {
	dataMap, ok := input.(map[string]interface{})
	if ok == true && dataMap["objType"] != nil {
		return true
	}
	return false
}

//func (sc *SmartContract) InsertObject()
func (sc *SmartContract) handleFunctions(stub shim.ChaincodeStubInterface) pb.Response {
	_SC_LOGGER.Info("InsidehandleFunctions")
	function, _ := stub.GetFunctionAndParameters()
	if function == "probe" {
		return sc.probe(stub)
	}
	return shim.Error("Invalid function provided")
}
