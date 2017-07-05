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
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Balance  int    `json:"balance"`
}

type TransactionInfo struct {
	Id        string `json:"id"`
	Amount    int    `json:"amount"`
	UserInfoA User   `json:"userInfoA"`
	UserInfoB User   `json:"userInfoB"`
	Status    string `json:"status"`
}

var tanggalFormat string

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}

	t := time.Now()
	tanggalFormat = t.Format("2006_01_02_150405")
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	if function == "transaction" {
		return t.Transaction(stub, args)
	} else if function == "create_user" {
		return t.CreateUser(stub, args)
	}

	return nil, nil
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	} else if function == "list_users" {
		return t.listUsers(stub, args)
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

func (t *SimpleChaincode) listUsers(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp string
	var err error

	valAsbytes, err := stub.GetState("users")
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for users}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}

func (t *SimpleChaincode) CreateUser(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3. name,password,balance to create user")
	}

	usersArray, err := stub.GetState("users")
	if err != nil {
		return nil, err
	}

	var users []string

	err = json.Unmarshal(usersArray, &users)

	if err != nil {
		return nil, err
	}

	users = append(users, args[0])

	b, err := json.Marshal(users)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Errors while creating json string for usertwo")
	}

	err = stub.PutState("users", b)
	if err != nil {
		return nil, err
	}

	var userone User
	userone.Name = args[0]
	userone.Password = args[1]
	balance, err := strconv.Atoi(args[2])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding at 3 place")
	}

	userone.Balance = balance

	b, err = json.Marshal(userone)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Errors while creating json string for userone")
	}

	err = stub.PutState(args[0], b)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (t *SimpleChaincode) Transaction(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var X int // Transaction value
	var err error

	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	// Get the state from the ledger
	// TODO: will be nice to have a GetAllState call to ledger
	Avalbytes, err := stub.GetState(args[0])
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	var userA User
	err = json.Unmarshal(Avalbytes, &userA)
	if err != nil {
		return nil, errors.New("Failed to marshal string to struct of userA")
	}

	Bvalbytes, err := stub.GetState(args[1])
	if err != nil {
		return nil, errors.New("Failed to get state")
	}

	var userB User
	err = json.Unmarshal(Bvalbytes, &userB)
	if err != nil {
		return nil, errors.New("Failed to marshal string to struct of userB")
	}

	// Perform the execution
	X, err = strconv.Atoi(args[2])
	if err != nil {
		return nil, errors.New("Third argument must be integer")
	}

	userA.Balance = userA.Balance - X
	userB.Balance = userB.Balance + X
	fmt.Printf("Aval = %d, Bval = %d\n", userA.Balance, userB.Balance)

	b, err := json.Marshal(userA)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Errors while creating json string for usera")
	}

	// Write the state back to the ledger
	err = stub.PutState(userA.Name, b)
	if err != nil {
		return nil, err
	}

	b, err = json.Marshal(userB)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Errors while creating json string for userb")
	}

	err = stub.PutState(userB.Name, b)
	if err != nil {
		return nil, err
	}

	var newTransactionInfo TransactionInfo

	Cvalbytes, err := stub.GetState(args[3])
	if err != nil {
		return nil, errors.New("Failed to get state")
	}

	err = json.Unmarshal(Cvalbytes, &newTransactionInfo)
	if err != nil {
		return nil, errors.New("Failed to marshal string to struct of newTransactionInfo")
	}

	newTransactionInfo.Id = args[3]
	newTransactionInfo.Amount = X
	newTransactionInfo.UserInfoA = userA
	newTransactionInfo.UserInfoB = userB
	newTransactionInfo.Status = "success"

	b, err = json.Marshal(newTransactionInfo)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Errors while creating json string for newTransactionInfo")
	}

	err = stub.PutState(args[3], b)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
