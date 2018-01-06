package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var _MAIN_LOGGER = shim.NewLogger("SmartContractMain")

// Init initializes chaincode.
func (sc *SmartContract) Init(stub shim.ChaincodeStubInterface) pb.Response {
	_MAIN_LOGGER.Infof("Inside the init method ")
	response := sc.init(stub)
	return response
}

//Invoke is the entry point for any transaction
func (sc *SmartContract) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, _ := stub.GetFunctionAndParameters()
	_MAIN_LOGGER.Infof("Inside the invoke method with %s", function)
	if function == "insertDemographyRecordBulk" {
		return sc.insertDemographyRecord(stub)
	} else if function == "insertMedicalRecordsBulk" {
		return sc.insertMedicalRecords(stub)
	} else if function == "saveMedicalRecord" {
		return sc.saveMedicalRecord(stub)
	} else if function == "retrieveMedicalRecords" {
		return sc.retrieveMedicalRecords(stub)
	} else if function == "modifyMedicalRecord" {
		return sc.modifyMedicalRecord(stub)
	}
	return sc.handleFunctions(stub)
}
func (sc *SmartContract) modifyMedicalRecord(stub shim.ChaincodeStubInterface) pb.Response {
	recordMap := make(map[string]interface{})
	_, args := stub.GetFunctionAndParameters()
	json.Unmarshal([]byte(args[0]), &recordMap)
	return sc.ModifyRecord(stub, recordMap, "recordId")
}
func (sc *SmartContract) retrieveMedicalRecords(stub shim.ChaincodeStubInterface) pb.Response {
	criteriaMap := make(map[string]string)
	selectionCriteria := ""
	_, args := stub.GetFunctionAndParameters()
	json.Unmarshal([]byte(args[0]), &criteriaMap)
	//The critetia map looks like the following
	//{"type":"", //W_MRECID|W_AADHAAR|W_DOCTORID|W_HCAID
	//	"arg1":"value1","arg2":"value2","arg3":"value3",}
	switch criteriaMap["type"] {
	case "W_AADHAAR":
		selectionCriteria = fmt.Sprintf("{\"patientAadharNo\":\"%s\",\"objType\":\"com.hc.mrec\"}", criteriaMap["patientAadharNo"])
	case "W_MRECID":
		selectionCriteria = fmt.Sprintf("{\"recordId\":\"%s\",\"objType\":\"com.hc.mrec\"}", criteriaMap["recordId"])
	case "W_HCAID":
		selectionCriteria = fmt.Sprintf("{\"recordCreator\":\"%s\",\"objType\":\"com.hc.mrec\"}", criteriaMap["recordCreator"])
	case "W_NODOCTOR":
		selectionCriteria = fmt.Sprintf("{ \"objType\":\"com.hc.mrec\",\"doctorResponded\":{ \"$exists\": false} } ")
	case "W_DOCTORID":
		selectionCriteria = fmt.Sprintf("{\"doctorResponded\":\"%s\",\"objType\":\"com.hc.mrec\"}", criteriaMap["doctorResponded"])
	default:
		_MAIN_LOGGER.Info("Error in type of search")
	}
	if selectionCriteria != "" {
		_MAIN_LOGGER.Infof("Search criteria %s", selectionCriteria)
		records := sc.RetriveRecords(stub, selectionCriteria)
		_MAIN_LOGGER.Infof("Number of records found %d", len(records))
		outputJSON, _ := json.Marshal(records)
		return shim.Success(outputJSON)
	}
	return shim.Error("Invalid search type")
}

//Inserts  medical records
func (sc *SmartContract) insertMedicalRecords(stub shim.ChaincodeStubInterface) pb.Response {
	var errMsgBuf bytes.Buffer
	_, args := stub.GetFunctionAndParameters()
	isValid, validationMessage, records := sc.ValidateObjectIntegrity(args[0])
	if isValid {
		isAllGood := true
		for _, record := range records {
			isSaveDone, msg := sc.ValidateAndInsertObject(stub, record, "recordId")
			isAllGood = isAllGood && isSaveDone
			if !isSaveDone {
				errMsgBuf.WriteString(msg + ",")
			}

		}
		if isAllGood {
			return shim.Success([]byte("All records saved"))
		}
		return shim.Error(errMsgBuf.String())
	}
	return shim.Error(validationMessage)

}

//Inserts the demography records
func (sc *SmartContract) insertDemographyRecord(stub shim.ChaincodeStubInterface) pb.Response {
	var errMsgBuf bytes.Buffer
	_, args := stub.GetFunctionAndParameters()
	isValid, validationMessage, records := sc.ValidateObjectIntegrity(args[0])
	if isValid {
		isAllGood := true
		for index, record := range records {
			isSaveDone, msg := sc.ValidateAndInsertObject(stub, record, "patientAadharNo")
			isAllGood = isAllGood && isSaveDone
			if !isSaveDone {
				errMsgBuf.WriteString(fmt.Sprintf("\"Record no %d %s\",", index, msg))
			}

		}
		if isAllGood {
			return shim.Success([]byte("All records saved"))
		}
		return shim.Error(errMsgBuf.String())
	}
	return shim.Error(validationMessage)

}

//This method takes  a composite input
func (sc *SmartContract) saveMedicalRecord(stub shim.ChaincodeStubInterface) pb.Response {
	var errMsgBuf bytes.Buffer
	_, args := stub.GetFunctionAndParameters()
	rootRec := make(map[string]interface{})
	parseErr := json.Unmarshal([]byte(args[0]), &rootRec)
	if parseErr == nil {
		demographicRec := rootRec["demographicDetail"]
		medicalRecord := rootRec["medicalRecord"]
		if demographicRec != nil && medicalRecord != nil {
			//Check if the demograhic detail existing or not
			demographicDataJSON, _ := json.Marshal(demographicRec)
			medicalRecordJSON, _ := json.Marshal(medicalRecord)
			isDRValid, validationMessage, recordsDR := sc.ValidateObjectIntegrity(string(demographicDataJSON))
			if isDRValid == false {
				_MAIN_LOGGER.Info("Demographic data object is not valid")
				errMsgBuf.WriteString(fmt.Sprintf("\"%s\",", validationMessage))
			}
			isMRValid, validationMessage, recordsMR := sc.ValidateObjectIntegrity(string(medicalRecordJSON))
			if isMRValid == false {
				_MAIN_LOGGER.Info("Medical report data object is not valid")
				errMsgBuf.WriteString(fmt.Sprintf("\"%s\",", validationMessage))
			}
			if isMRValid == true && isDRValid == true {
				saveMRData := true
				_MAIN_LOGGER.Infof("Saving the demographic record")
				isSuccess, msg := sc.ValidateAndInsertObject(stub, recordsDR[0], "patientAadharNo")
				if isSuccess == true {
					_MAIN_LOGGER.Info(" Demographic data saved successfully")
				} else if isSuccess == false && msg == "Id already exists" {
					_MAIN_LOGGER.Infof(" Demographic data already exists for record %s", string(demographicDataJSON))
				} else {
					saveMRData = false
					_MAIN_LOGGER.Infof("Failed to save demographic record")
					errMsgBuf.WriteString(fmt.Sprintf("\"%s\",", msg))
				}
				if saveMRData == true {
					saveMRData, msg = sc.ValidateAndInsertObject(stub, recordsMR[0], "recordId")
				}
				if saveMRData == true {
					_MAIN_LOGGER.Info("Both the records are saved successfully")
					return shim.Success([]byte("Records save successfully "))
				}
				errMsgBuf.WriteString(fmt.Sprintf("\"%s\",", msg))
				_MAIN_LOGGER.Infof("Error returned %s", errMsgBuf.String())
				return shim.Error(errMsgBuf.String())
			}
			_MAIN_LOGGER.Info("Error retuning %s", errMsgBuf.String())
			return shim.Error(errMsgBuf.String())
		}
		return shim.Error("Both demographicDetail and medicalRecord is necessary")
	}
	return shim.Error("Invaild input provided")

}
func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		_MAIN_LOGGER.Criticalf("Error starting  chaincode: %v", err)
	}
}
