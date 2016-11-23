package main

import (
  	"errors"
	"fmt"
	"strconv"
  	"time"
	"encoding/json"
  	"./imports"
	
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type SimpleChaincode struct {
}

type Account struct {
	ID          string  `json:"ID"`
	CashBalance int     `json:"CashBalance"`
}

type Transaction struct {
	ID          string   `json:"ID"`
 	Timestamp   string   `json:"timestamp"`
	FromUser    string   `json:"fromUser"`
	ToUser      string   `json:"toUser"`
	Quantity    int      `json:"quantity"`
}


// ============================================================================================================================
// Init - reset all the things
// ============================================================================================================================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var Aval int
	var err error

	// Initialize the chaincode
	Aval, err = strconv.Atoi(args[0])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}

	// Write the state to the ledger
	err = stub.PutState("abc", []byte(strconv.Itoa(Aval)))				//making a test var "abc", I find it handy to read/write to it right away to test the network
	if err != nil {
		return nil, err
	}

	// var empty []string
	// jsonAsBytes, _ := json.Marshal(empty)								//marshal an emtpy array of strings to clear the index
	// err = stub.PutState(marbleIndexStr, jsonAsBytes)
	// if err != nil {
	// 	return nil, err
	// }

	return nil, nil
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


// // ============================================================================================================================
// // Run - Our entry point for Invocations - [LEGACY] obc-peer 4/25/2016
// // ============================================================================================================================
// func (t *SimpleChaincode) Run(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
// 	fmt.Println("run is running " + function)
// 	return t.Invoke(stub, function, args)
// }

// ============================================================================================================================
// Invoke - Our entry point for  TO REMOVE SHITTTTTT
// ============================================================================================================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

  if function == "init" {													//initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	}  else if function == "write" {											//writes a value to the chaincode state
		return t.Write(stub, args)
	} else if function == "createAccount" {
    return t.CreateAccount(stub, args)
  } else if function == "set_user" {										//change owner of a marble
		res, err := t.set_user(stub, args)											//lets make sure all open trades are still valid
		return res, err
	}
	fmt.Println("invoke did not find func: " + function)					//error

	return nil, errors.New("Received unknown function invocation")
}

// ============================================================================================================================
// Query - Our entry point for Queries
// ============================================================================================================================
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" {													//read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)						//error

	return nil, errors.New("Received unknown function query")
}

// ============================================================================================================================
// Read - read a variable from chaincode state
// ============================================================================================================================
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var name, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the var to query")
	}

	name = args[0]
	valAsbytes, err := stub.GetState(name)									//get the var from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil													//send it onward
}

func (t *SimpleChaincode) Write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var name, value string // Entities
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the variable and value to set")
	}

	name = args[0]															//rename for funsies
	value = args[1]
	err = stub.PutState(name, []byte(value))								//write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (t *SimpleChaincode) CreateAccount(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
  // Obtain the username to associate with the account
  var username string
  var err error
 	fmt.Println("running write()")

  if len(args) != 1 {
     fmt.Println("Error obtaining username")
     return nil, errors.New("createAccount accepts a single username argument")
  }
  username = args[0]

  var account = Account{ID: username, CashBalance: 500}
  accountBytes, err := json.Marshal(&account)

  err = stub.PutState(username, accountBytes)
  if err != nil {
     return nil, err
  }
  return nil, nil
}
// ============================================================================================================================
// Set Trade - create an open trade for a marble you want with marbles you have
// ============================================================================================================================
func (t *SimpleChaincode) set_user(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
  var toRes Account
	//     0         1        2
	// "fromUser", "500", "toUser",
	if len(args) < 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	fromAccountAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return nil, errors.New("Failed to get thing")
	}
  toAccountAsBytes, err := stub.GetState(args[2])
	if err != nil {
		return nil, errors.New("Failed to get thing")
	}


	fromRes := Account{}
	json.Unmarshal(fromAccountAsBytes, &fromRes)										//un stringify it aka JSON.parse()

  toRes = Account{}
	json.Unmarshal(toAccountAsBytes, &toRes)



	accountBalance := fromRes.CashBalance


  transferAmount, err := strconv.Atoi(args[1])
   if err != nil {
      // handle error
   }
  if(accountBalance < transferAmount) {
    fmt.Println("- Insufficient funds")
    return nil, nil
  }

  toRes.CashBalance = toRes.CashBalance + transferAmount
  fromRes.CashBalance = fromRes.CashBalance - transferAmount

transID := uuid.NewV4().String()
  timestamp := time.Now().Format(time.RFC3339)
  trans := Transaction{ID: transID, Timestamp: timestamp, FromUser: fromRes.ID, ToUser: toRes.ID, Quantity: transferAmount}
  transBytes, err := json.Marshal(&trans)

  err = stub.PutState(transIDStr, transBytes)
  if err != nil {
     return nil, err
  }
	
	toJsonAsBytes, _ := json.Marshal(toRes)
	err = stub.PutState(args[2], toJsonAsBytes)								//rewrite the marble with id as key
	if err != nil {
		return nil, err
	}

  fromJsonAsBytes, _ := json.Marshal(fromRes)
	err = stub.PutState(args[0], fromJsonAsBytes)								//rewrite the marble with id as key
	if err != nil {
		return nil, err
	}

	fmt.Println("- end set trade")
	return nil, nil
}
