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
	"reflect"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Participant struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Balance  int    `json:"balance"`
}

type TransactionInfo struct {
	TransactionInfoId string      `json:"transactionInfoId"`
	TransactionId     string      `json:"transactionId"`
	Amount            int         `json:"amount"`
	ParticipantInfoA  Participant `json:"participantInfoA"`
	ParticipantInfoB  Participant `json:"participantInfoB"`
	Status            string      `json:"status"`
	Description       string      `json:"description"`
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
	//	if len(args) != 12 {
	//		return nil, errors.New("Incorrect number of arguments. Execting 6")
	//	}
	//
	//	var participantsArray []string
	//
	//	var participantone Participant
	//	participantone.Name = args[0]
	//	participantone.Password = args[1]
	//	balance, err := strconv.Atoi(args[2])
	//	if err != nil {
	//		return nil, errors.New("Expecting integer value for asset holding at 3 place")
	//	}
	//
	//	participantone.Balance = balance
	//
	//	b, err := json.Marshal(participantone)
	//	if err != nil {
	//		fmt.Println(err)
	//		return nil, errors.New("Errors while creating json string for participantone")
	//	}
	//
	//	err = stub.PutState(args[0], b)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	participantone.Name = args[3]
	//	participantone.Password = args[4]
	//	balance, err = strconv.Atoi(args[5])
	//	if err != nil {
	//		return nil, errors.New("Expecting integer value for asset holding at 3 place")
	//	}
	//
	//	participantone.Balance = balance
	//
	//	b, err = json.Marshal(participantone)
	//	if err != nil {
	//		fmt.Println(err)
	//		return nil, errors.New("Errors while creating json string for participantone")
	//	}
	//
	//	err = stub.PutState(args[3], b)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	participantone.Name = args[6]
	//	participantone.Password = args[7]
	//	balance, err = strconv.Atoi(args[8])
	//	if err != nil {
	//		return nil, errors.New("Expecting integer value for asset holding at 3 place")
	//	}
	//
	//	participantone.Balance = balance
	//
	//	b, err = json.Marshal(participantone)
	//	if err != nil {
	//		fmt.Println(err)
	//		return nil, errors.New("Errors while creating json string for participantone")
	//	}
	//
	//	err = stub.PutState(args[6], b)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	participantone.Name = args[9]
	//	participantone.Password = args[10]
	//	balance, err = strconv.Atoi(args[11])
	//	if err != nil {
	//		return nil, errors.New("Expecting integer value for asset holding at 3 place")
	//	}
	//
	//	participantone.Balance = balance
	//
	//	b, err = json.Marshal(participantone)
	//	if err != nil {
	//		fmt.Println(err)
	//		return nil, errors.New("Errors while creating json string for participantone")
	//	}
	//
	//	err = stub.PutState(args[9], b)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	participantsArray = append(participantsArray, args[0])
	//	participantsArray = append(participantsArray, args[3])
	//	participantsArray = append(participantsArray, args[6])
	//	participantsArray = append(participantsArray, args[9])
	//
	//	b, err = json.Marshal(participantsArray)
	//	if err != nil {
	//		fmt.Println(err)
	//		return nil, errors.New("Errors while creating json string for participanttwo")
	//	}
	//
	//	err = stub.PutState("participants", b)
	//	if err != nil {
	//		return nil, err
	//	}

	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	if function == "transaction" {
		return t.Transaction(stub, args)
	} else if function == "create_participant" {
		return t.CreateParticipant(stub, args)
	} else if function == "init_transaction" {
		return t.InitTransaction(stub, args)
	}

	return nil, nil
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	} else if function == "list_participants" {
		return t.listParticipants(stub, args)
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

func (t *SimpleChaincode) listParticipants(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp string
	var err error

	valAsbytes, err := stub.GetState("participants")
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for participants}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}

func (t *SimpleChaincode) CreateParticipant(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3. name,password,balance to create participant")
	}

	participantsArray, err := stub.GetState("participants")
	if err != nil {
		return nil, err
	}

	var participants []string

	err = json.Unmarshal(participantsArray, &participants)

	if err != nil {
		return nil, err
	}

	participants = append(participants, args[0])

	b, err := json.Marshal(participants)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Errors while creating json string for participanttwo")
	}

	err = stub.PutState("participants", b)
	if err != nil {
		return nil, err
	}

	var participantone Participant
	participantone.Name = args[0]
	participantone.Password = args[1]
	balance, err := strconv.Atoi(args[2])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding at 3 place")
	}

	participantone.Balance = balance

	b, err = json.Marshal(participantone)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Errors while creating json string for participantone")
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

	if len(args) != 7 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	// Get the state from the ledger
	// TODO: will be nice to have a GetAllState call to ledger
	Avalbytes, err := stub.GetState(args[2])
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	var participantA Participant
	err = json.Unmarshal(Avalbytes, &participantA)
	if err != nil {
		return nil, errors.New("Failed to marshal string to struct of participantA")
	}

	Bvalbytes, err := stub.GetState(args[3])
	if err != nil {
		return nil, errors.New("Failed to get state")
	}

	var participantB Participant
	err = json.Unmarshal(Bvalbytes, &participantB)
	if err != nil {
		return nil, errors.New("Failed to marshal string to struct of participantB")
	}

	// Perform the execution
	X, err = strconv.Atoi(args[2])
	if err != nil {
		return nil, errors.New("Third argument must be integer")
	}

	participantA.Balance = participantA.Balance - X
	participantB.Balance = participantB.Balance + X
	fmt.Printf("Aval = %d, Bval = %d\n", participantA.Balance, participantB.Balance)

	b, err := json.Marshal(participantA)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Errors while creating json string for participanta")
	}

	// Write the state back to the ledger
	err = stub.PutState(participantA.Name, b)
	if err != nil {
		return nil, err
	}

	b, err = json.Marshal(participantB)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Errors while creating json string for participantb")
	}

	err = stub.PutState(participantB.Name, b)
	if err != nil {
		return nil, err
	}

	var newTransactionInfo TransactionInfo

	//	Cvalbytes, err := stub.GetState(args[3])
	//	if err != nil {
	//		return nil, errors.New("Failed to get state")
	//	}

	/*s1 := `{ "Name": "davin" , "Password": "password" , "Balance": 100 }`
	s2 := `{ "Name": "ardian" , "Password": "password" , "Balance": 200 }`
	bytes1 := []byte(s1)
	bytes2 := []byte(s2)

	var participantA_unmarshal Participant
	err = json.Unmarshal(bytes1, &participantA_unmarshal)
	if err != nil {
		panic(err)
	}

	var participantB_unmarshal Participant
	err = json.Unmarshal(bytes2, &participantB_unmarshal)
	if err != nil {
		panic(err)
	}*/

	newTransactionInfo.TransactionInfoId = args[0]
	newTransactionInfo.TransactionId = args[1]
	newTransactionInfo.Amount = X
	newTransactionInfo.ParticipantInfoA = participantA
	newTransactionInfo.ParticipantInfoB = participantB
	newTransactionInfo.Status = args[4]
	newTransactionInfo.Description = args[5]

	//	err = json.Unmarshal(Cvalbytes, &newTransactionInfo)
	//	if err != nil {
	//		return nil, errors.New("Failed to marshal string to struct of newTransactionInfo")
	//	}

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

func (t *SimpleChaincode) InitTransaction(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var X int = 0 // Transaction value
	var err error

	if len(args) != 5 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	// Get the state from the ledger
	// TODO: will be nice to have a GetAllState call to ledger
	Avalbytes, err := stub.GetState(args[2])
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	var participantA Participant
	err = json.Unmarshal(Avalbytes, &participantA)
	if err != nil {
		return nil, errors.New("Failed to marshal string to struct of participantA")
	}

	participantA.Balance = X

	b, err := json.Marshal(participantA)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Errors while creating json string for participanta")
	}

	// Write the state back to the ledger
	err = stub.PutState(participantA.Name, b)
	if err != nil {
		return nil, err
	}

	var participantB Participant

	var oldTransactionInfo TransactionInfo
	var newTransactionInfo TransactionInfo
	newTransactionInfo.TransactionInfoId = args[0]
	newTransactionInfo.TransactionId = args[1]
	newTransactionInfo.Amount = X
	newTransactionInfo.ParticipantInfoA = participantA
	newTransactionInfo.ParticipantInfoB = participantB
	newTransactionInfo.Status = args[3]
	newTransactionInfo.Description = args[4]

	assetBytes, err := stub.GetState(args[0])
	if err != nil || len(assetBytes) != 0 {

		// This is an update scenario
		err = json.Unmarshal(assetBytes, &oldTransactionInfo)
		if err != nil {
			err = errors.New("Unable to unmarshal JSON data from stub")
			return nil, err
			// state is an empty instance of asset state
		}
		// Merge partial state updates
		oldTransactionInfo, err = t.mergePartialState(oldTransactionInfo, newTransactionInfo)
		if err != nil {
			err = errors.New("Unable to merge state")
			return nil, err
		}

	}

	b, err = json.Marshal(oldTransactionInfo)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Errors while creating json string for newTransactionInfo")
	}

	err = stub.PutState(args[0], b)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (t *SimpleChaincode) mergePartialState(oldState TransactionInfo, newState TransactionInfo) (TransactionInfo, error) {

	old := reflect.ValueOf(&oldState).Elem()
	new := reflect.ValueOf(&newState).Elem()
	for i := 0; i < old.NumField(); i++ {
		oldOne := old.Field(i)
		newOne := new.Field(i)
		if !reflect.ValueOf(newOne.Interface()).IsNil() {
			oldOne.Set(reflect.Value(newOne))
		}
	}
	return oldState, nil
}
