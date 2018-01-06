package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func checkInit(t *testing.T, stub *shim.MockStub) {
	args := make([][]byte, 0)
	args = append(args, []byte("init"))
	res := stub.MockInit("1", args)
	if res.Status != shim.OK {
		fmt.Println("Init failed", string(res.Message))
		t.FailNow()
	}
}
func checkProbe(t *testing.T, stub *shim.MockStub) {
	args := make([][]byte, 0)
	args = append(args, []byte("probe"))

	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		t.FailNow()
	} else {
		fmt.Printf("\n %s\n", string(res.Payload))
	}
}
func checkInsertDemographicRecordsValid(t *testing.T, stub *shim.MockStub, payload string) {
	args := make([][]byte, 0)
	args = append(args, []byte("insertDemographyRecord"))
	args = append(args, []byte(payload))

	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		t.FailNow()
	} else {
		fmt.Printf("%s\n", string(res.Payload))
	}
}
func checkInsertDemographicRecordsInValid(t *testing.T, stub *shim.MockStub, payload string) {
	args := make([][]byte, 0)
	args = append(args, []byte("insertDemographyRecord"))
	args = append(args, []byte(payload))

	res := stub.MockInvoke("1", args)
	if res.Status != shim.ERROR {
		t.FailNow()
	} else {
		fmt.Printf("%s\n", res.Message)
	}
}

//Test_Init tests the input
func Test_Init(t *testing.T) {
	scc := new(SmartContract)
	stub := shim.NewMockStub("hcdm", scc)

	// Init A=123 B=234
	checkInit(t, stub)

}
func Test_Probe(t *testing.T) {
	scc := new(SmartContract)
	stub := shim.NewMockStub("hcdm", scc)
	checkInit(t, stub)
	checkProbe(t, stub)
}
func Test_InsertRecords(t *testing.T) {
	scc := new(SmartContract)
	stub := shim.NewMockStub("hcdm", scc)
	checkInit(t, stub)
	testData := make(map[string]string)
	testData["objType"] = "com.hc.patientinfo"
	testData["patientAadharNo"] = "1238721321873"
	jsonBytes, _ := json.Marshal(testData)
	checkInsertDemographicRecordsValid(t, stub, string(jsonBytes))
	//Validating array of objects
	arrayData := make([]map[string]string, 0)
	for index := 0; index < 10; index++ {
		testData := make(map[string]string)
		testData["objType"] = "com.hc.patientinfo"
		testData["patientAadharNo"] = fmt.Sprintf("patientAadharNo%d", index)
		arrayData = append(arrayData, testData)
	}
	jsonBytes, _ = json.Marshal(arrayData)
	checkInsertDemographicRecordsValid(t, stub, string(jsonBytes))
	testDataInvalid := make(map[string]string)
	testDataInvalid["objType"] = "com.hc.patientinfo"
	testDataInvalid["patientAadharNo"] = "1238721321873"
	jsonBytes, _ = json.Marshal(testDataInvalid)
	checkInsertDemographicRecordsInValid(t, stub, string(jsonBytes))
	arrayData = make([]map[string]string, 0)
	for index := 0; index < 10; index++ {
		testData := make(map[string]string)
		testData["objType"] = "com.hc.patientinfo"
		testData["patientAadharNo"] = fmt.Sprintf("patientAadharNo%d", index)
		arrayData = append(arrayData, testData)
	}

	jsonBytes, _ = json.Marshal(arrayData)
	checkInsertDemographicRecordsInValid(t, stub, string(jsonBytes))
}
