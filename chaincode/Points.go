package main

import (
  "errors"
	"fmt"
	"strconv"
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type SimpleChaincode struct {
}

type Account struct {
	ID                 string          `json:"ID"`
	Password           string          `json:"Password"`
	GiveBalance        int             `json:"GiveBalance"`
	PointsBalance      int             `json:"PointsBalance"`
}

type Product struct{
  ID                 string          `json:"ID"`
  Name               string          `json:"Name"`
  Cost               int             `json:"Cost"`
  Owner              *string         `json:"Owner"`
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
  } else if function == "createProduct" {
    return t.CreateProduct(stub, args)
  } else if function == "purchaseProduct" {
    return t.PurchaseProduct(stub, args)
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

  if len(args) != 2 {
     fmt.Println("Error obtaining username")
     return nil, errors.New("createAccount accepts a single username argument")
  }
  username = args[0]
  password := args[1]

  var account = Account{ID: username, Password: password, GiveBalance: 500, PointsBalance:50}
  accountBytes, err := json.Marshal(&account)

  err = stub.PutState(username, accountBytes)
  if err != nil {
     return nil, err
  }
  return nil, nil
}

func (t *SimpleChaincode) CreateProduct(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
  var err error
 	fmt.Println("running write()")

  if len(args) != 3 {
     return nil, errors.New("CreateProduct accepts a 3 argument")
  }
  ID := args[0]
  name := args[1]
  cost, err := strconv.Atoi(args[2])
   if err != nil {
      return nil, err
   }

  prod := Product{ID: ID, Name: name, Cost: cost, Owner:nil}
  prodBytes, err := json.Marshal(&prod)

  err = stub.PutState(ID, prodBytes)
  if err != nil {
     return nil, err
  }

  return nil, nil
}

func (t *SimpleChaincode) PurchaseProduct(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
  var err error
 	fmt.Println("running write()")

  if len(args) != 2 {
     return nil, errors.New("createAccount 2 argument")
  }
  ProdID := args[0]
  username := args[1]

  fromAccountAsBytes, err := stub.GetState(username)
	if err != nil {
		return nil, errors.New("Failed to get thing")
	}

  prodAsBytes, err := stub.GetState(ProdID)
	if err != nil {
		return nil, errors.New("Failed to get thing")
	}

  fromRes := Account{}
	json.Unmarshal(fromAccountAsBytes, &fromRes)										//un stringify it aka JSON.parse()

  prodRes := Product{}
	json.Unmarshal(prodAsBytes, &prodRes)

  if(prodRes.Owner != nil){
    return nil, errors.New("Already owned")
  }

  if(fromRes.PointsBalance < prodRes.Cost){
    fmt.Println("- Insufficient funds")
    return []byte("Insufficient funds"), errors.New("Insufficient funds")
  }

  prodRes.Owner = &fromRes.ID
  fromRes.PointsBalance -= prodRes.Cost

  fromJsonAsBytes, _ := json.Marshal(fromRes)
	err = stub.PutState(username, fromJsonAsBytes)								//rewrite the marble with id as key
	if err != nil {
		return nil, err
	}

  toJsonAsBytes, _ := json.Marshal(prodRes)
  err = stub.PutState(ProdID, toJsonAsBytes)								//rewrite the marble with id as key
  if err != nil {
    return nil, err
  }

	fmt.Println("- end set PurchaseProduct")
	return nil, nil
}

// ============================================================================================================================
// Set Trade - create an open trade for a marble you want with marbles you have
// ============================================================================================================================
func (t *SimpleChaincode) set_user(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
  var toRes Account
	//     0         1        2        3
	// "fromUser", "500", "toUser", "reason"
	if len(args) < 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
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



	accountBalance := fromRes.GiveBalance


  transferAmount, err := strconv.Atoi(args[1])
   if err != nil {
      // handle error
   }
  if(accountBalance < transferAmount) {
    fmt.Println("- Insufficient funds")
    return nil, nil
  }

  toRes.PointsBalance = toRes.PointsBalance + transferAmount
  fromRes.GiveBalance = fromRes.GiveBalance - transferAmount

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
