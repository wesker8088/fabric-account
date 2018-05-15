package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type SmartContract struct {
}

type Account struct {
	Name   string `json:"name"`
	Gender string `json:"gender"`
	Age    string `json:"age"`
	Mail   string `json:"mail"`
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}


func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	function, args := APIstub.GetFunctionAndParameters()
	if function == "query" {
		return s.query(APIstub, args)
	} else if function == "init" {
		return s.initAccont(APIstub)
	} else if function == "create" {
		return s.create(APIstub, args)
	} else if function == "list" {
		return s.list(APIstub)
	} else if function == "update" {
		return s.update(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) query(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	carAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(carAsBytes)
}

func (s *SmartContract) initAccont(APIstub shim.ChaincodeStubInterface) sc.Response {
	Accounts := []Account{
		Account{Name: "wesker", Gender: "male", Age: "26", Mail: "wesker@gmail.com"},
		Account{Name: "jill", Gender: "female", Age: "21", Mail: "jill@gmail.com"},
		Account{Name: "leon", Gender: "male", Age: "22", Mail: "leon@gmail.com"},
		Account{Name: "chris", Gender: "male", Age: "25", Mail: "chris@gmail.com"},
	}

	i := 0
	for i < len(Accounts) {
		fmt.Println("i is ", i)
		accountAsBytes, _ := json.Marshal(Accounts[i])
		APIstub.PutState("ACCOUNT"+strconv.Itoa(i), accountAsBytes)
		fmt.Println("Added", Accounts[i])
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) create(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5{
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	var account = Account{Name: args[1], Gender: args[2], Age: args[3], Mail: args[4]}
        fmt.Println("New Added:", account)
	accountAsBytes, _ := json.Marshal(account)
        fmt.Println("New args[0]:", args[0])
        fmt.Println("New accountAsBytes:", accountAsBytes)
	APIstub.PutState(args[0], accountAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) list(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "ACCOUNT0"
	endKey := "ACCOUNT999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- listAllAcount:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) update(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
        fmt.Println("Account update start")
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	accountAsBytes, _ := APIstub.GetState(args[0])
	account := Account{}

	json.Unmarshal(accountAsBytes, &account)
	account.Name = args[1]
        fmt.Println("Account update name:",args[1])
	accountAsBytes, _ = json.Marshal(account)
	APIstub.PutState(args[0], accountAsBytes)

        fmt.Println("Account update end")
	return shim.Success(nil)
}

func main() {

	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}

