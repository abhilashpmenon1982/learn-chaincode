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
        "encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Contract struct {
	Customer     string `json:"customer"`
	Provider     string `json:"provider"`
	Period       string `json:"period"`
        Status       string `json:"status"`
}






// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Execting 4")
	}

	var contractArray []string

	var contractone Contract
	contractone.Customer = args[0]
	contractone.Provider = args[1]
	contractone.Period = args[2]
	contractone.Status = args[3]

	b, err := json.Marshal(contractone)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Errors while creating json string for contractone")
	}

	err = stub.PutState(args[0], b)
	if err != nil {
		return nil, err
	}


	contractArray = append(contractArray, args[0])

	b, err = json.Marshal(contractArray)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Errors while creating json string for contract")
	}

	err = stub.PutState("contracts", b)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	if function == "transaction" {
		return t.Transaction(stub, args)
	} else if function == "create_contract" {
		return t.CreateContract(stub, args)
	}

	return nil, nil
}


func (t *SimpleChaincode) CreateContract(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4. customer,provider,period,status to create contract")
	}

	contractArray, err := stub.GetState("contracts")
	if err != nil {
		return nil, err
	}

	var contracts []string

	err = json.Unmarshal(contractArray, &contracts)

	if err != nil {
		return nil, err
	}

	contracts = append(contracts, args[0])

	b, err := json.Marshal(contracts)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Errors while creating json string for contract")
	}

	err = stub.PutState("contracts", b)
	if err != nil {
		return nil, err
	}

	var contractone Contract
	contractone.Customer = args[0]
	contractone.Provider = args[1]
	contractone.Period = args[2]
	contractone.Status = args[3]

	b, err = json.Marshal(contractone)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Errors while creating json string for contract")
	}

	err = stub.PutState(args[0], b)
	if err != nil {
		return nil, err
	}

	return nil, nil
}



func (t *SimpleChaincode) Transaction(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {


	var err error

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}

	// Get the state from the ledger
	// TODO: will be nice to have a GetAllState call to ledger
	Avalbytes, err := stub.GetState(args[0])
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	var contractone Contract
	err = json.Unmarshal(Avalbytes, &contractone)
	if err != nil {
		return nil, errors.New("Failed to marshal string to struct of contractone")
	}


	contractone.Period = args[1]
	contractone.Status = args[2]
	fmt.Printf("Period = %s, Status = %s\n", contractone.Period, contractone.Status)

	b, err := json.Marshal(contractone)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Errors while creating json string for contractone")
	}

	// Write the state back to the ledger
	err = stub.PutState(contractone.Customer, b)
	if err != nil {
		return nil, err
	}


	return nil, nil
}


// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	} else if function == "list_contracts" {
		return t.listContracts(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
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

func (t *SimpleChaincode) listContracts(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp string
	var err error

	valAsbytes, err := stub.GetState("contracts")
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for contracts}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}
