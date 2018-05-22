/*
Copyright IBM Corp 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"errors"
	"fmt"
	"encoding/json"
	"strings"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {

}

// Structure of Patient Details : Custom block

type Patient struct {
	Username    							string  `json:"Username"`
	Name         							string  `json:"Name"`
	DescriptionOfCurrentAilment 					string  `json:"DescriptionOfCurrentAilment"`
	DateOfBirth 							string  `json:"DateOfBirth"`
	Gender 								string  `json:"Gender"`
	ReportType 							string  `json:"ReportType"`
	PreLunch							string  `json:"PreLunch"`
	PostLunch							string  `json:"PostLunch"`
	MinSize								string  `json:"MinSize"`
	MaxSize 							string  `json:"MaxSize"`
	Disease 							string  `json:"Disease"`
	OnGoingMedication 						string  `json:"OnGoingMedication"`
	Duration 							string  `json:"Duration"`
	Titanus 							string  `json:"Titanus"`
	HepatitisA 							string  `json:"HepatitisA"`
	HepatitisB 							string  `json:"HepatitisB"`
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	err := stub.PutState("hello_blockchain", []byte(args[0]))
	if err != nil {
		return nil, err
	}

	return nil, nil
}


// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	//fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args)
	} else if function== "assign" {
		return t.assign(stub, args)
		//cer, key
	}

	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	}
	if function == "readAssign" { //read a variable
		return t.readAssign(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

//asign pair of (cer, key)
func (t *SimpleChaincode) assign(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	//PutState(cer, key + uid)
	var merg string
	merg = args[1]+":"+ args[2]
	err = stub.PutState(args[0], []byte(merg)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// write - invoke function to write key/value pair
func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	var rtype string
	fmt.Println("running write()")

	if len(args) != 14 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	m_patient 				:= &Patient{}
	m_patient.Username 			= args[0]
	m_patient.Name 				= args[1]
	m_patient.DescriptionOfCurrentAilment	= args[2]
	m_patient.DateOfBirth			= args[3]
	m_patient.Gender			= args[4]
	rtype					= args[5]

	if(strings.ToLower(rtype)=="diabetes"){
		m_patient.ReportType	= args[6]
		m_patient.PreLunch	= args[7]
		m_patient.PostLunch	= args[8]
	}else if(strings.ToLower(rtype)=="kidney"){
		m_patient.ReportType	= args[5]
		m_patient.MinSize	= args[6]
		m_patient.MaxSize	= args[7]
	}

	m_patient.Disease			= args[8]
	m_patient.OnGoingMedication		= args[9]
	m_patient.Duration			= args[10]
	m_patient.Titanus			= args[11]
	m_patient.HepatitisA			= args[12]
	m_patient.HepatitisB			= args[13]


	var key = args[0]

	value, err := json.Marshal(&m_patient)

	if err != nil {
		return nil, err
	}

	err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// read - query function to read key/value pair
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}
	return valAsbytes, nil
}

// read - query function to read tripplet (cer, key+':'+uid)
func (t *SimpleChaincode) readAssign(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	/*if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}*/

	var cer, key, /*uid,*/ jsonResp string
	var err error

	cer=args[0]
	key=args[1]
	//uid=args[2]

	// cer = key + uid
	valAsbytes, err := stub.GetState(cer)

	s := strings.Split( string( valAsbytes), ":")

	if (s[0] == key){
		valAsbytes, err := stub.GetState(s[1])
		if err != nil {
			jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
			return nil, errors.New(jsonResp)
		}
		return valAsbytes, nil
	}
	return nil, err
}
