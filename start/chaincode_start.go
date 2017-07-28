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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type DeliveryInfo struct {
	PackageId   string `json:"packageId"`
	From        string `json:"from"`
	Destination string `json:"destination"`
	MinTemp     string `json:"minTemp"`
	MaxTemp     string `json:"maxTemp"`
	Carrier     string `json:"carrier"`
	Status      string `json:"status"`
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

	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	if function == "delivery" {
		return t.Delivery(stub, args)
	} else if function == "init_delivery" {
		return t.InitDelivery(stub, args)
	}

	return nil, nil
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
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

func (t *SimpleChaincode) Delivery(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var bufferIdHistoryDeliveryPackage bytes.Buffer
	var err error
	var listHistoryDeliveryPackage []DeliveryInfo

	if len(args) != 7 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	var newDeliveryInfo DeliveryInfo
	newDeliveryInfo.PackageId = args[0]
	newDeliveryInfo.From = args[1]
	newDeliveryInfo.Destination = args[2]
	newDeliveryInfo.MinTemp = args[3]
	newDeliveryInfo.MaxTemp = args[4]
	newDeliveryInfo.Carrier = args[5]
	newDeliveryInfo.Status = args[6]

	b, err := json.Marshal(newDeliveryInfo)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Errors while creating json string for newDeliveryInfo")
	}

	err = stub.PutState(args[0], b)
	if err != nil {
		return nil, err
	}

	//add to history delvery
	bufferIdHistoryDeliveryPackage.WriteString("history_delivery_")
	bufferIdHistoryDeliveryPackage.WriteString(args[0])

	historyDeliveryPackageArray, err := stub.GetState(bufferIdHistoryDeliveryPackage.String())
	if err != nil || len(historyDeliveryPackageArray) != 0 {

		listHistoryDeliveryPackageTemp, err := stub.GetState(bufferIdHistoryDeliveryPackage.String())
		if err != nil {
			return nil, err
		}

		var newHistoryDeliveryPackageArray []DeliveryInfo
		err = json.Unmarshal(listHistoryDeliveryPackageTemp, &newHistoryDeliveryPackageArray)
		if err != nil {
			return nil, err
		}

		newHistoryDeliveryPackageArray = append(newHistoryDeliveryPackageArray, newDeliveryInfo)

		d, err := json.Marshal(newHistoryDeliveryPackageArray)
		if err != nil {
			fmt.Println(err)
			return nil, errors.New("Errors while creating json string for participanttwo")
		}

		err = stub.PutState(bufferIdHistoryDeliveryPackage.String(), d)
		if err != nil {
			return nil, err
		}

	} else {

		listHistoryDeliveryPackage = []DeliveryInfo{
			{
				PackageId:   args[0],
				From:        args[1],
				Destination: args[2],
				MinTemp:     args[3],
				MaxTemp:     args[4],
				Carrier:     args[5],
				Status:      args[6],
			},
		}

		c, err := json.Marshal(listHistoryDeliveryPackage)
		if err != nil {
			fmt.Println(err)
			return nil, errors.New("Errors while creating json string for participanttwo")
		}

		err = stub.PutState(bufferIdHistoryDeliveryPackage.String(), c)
		if err != nil {
			return nil, err
		}

		stateJSON1 := []byte(args[0])
		err = stub.PutState("lastId", stateJSON1)
		if err != nil {
			return nil, err
		}

	}

	return nil, nil
}

func (t *SimpleChaincode) InitDelivery(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	var buffer bytes.Buffer

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	stateJSON1 := []byte(args[0])
	err = stub.PutState("lastId", stateJSON1)
	if err != nil {
		return nil, err
	}

	buffer.WriteString("package_")
	buffer.WriteString(args[0])
	stateJSON2 := []byte(buffer.String())
	err = stub.PutState("packageId", stateJSON2)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
