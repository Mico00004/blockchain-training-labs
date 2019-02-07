/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/*
 * The sample smart contract for documentation topic:
 * Writing Your First Blockchain Application
 */

package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	//"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the   structure, with 4 properties.  Structure tags are used by encoding/json library
type Invoice struct {
	InvoiceNumber   string `json:"invoiceNumber"`
	BilledTo  string `json:"billedTo"`
	InvoiceDate string `json:"invoiceDate"`
	InvoiceAmount  string `json:"invoiceAmount"`
	ItemDescription string `json:"itemDescription"`
	GR string `json:"gr"`
	IsPaid string `json:"isPaid"`
	PaidAmount string `json:"paidAmount"`
	Repaid string `json:"repaid"`
	RepaymentAmount string `json:"RepaymentAmount"`
}

/*
 * The Init method is called when the Smart Contract "fabcar" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "fabInvoice"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryInvoice" {
		return s.queryInvoice(APIstub, args)
	} else if function == "raiseInvoice" {
		return s.raiseInvoice(APIstub, args)
	} else if function == "queryAllInvoices" {
		return s.queryAllInvoices(APIstub)
	} else if function == "goodsReceive" {
		return s.goodsReceive(APIstub, args)	
	} else if function == "bankPaymentToSupplier" {
		return s.bankPaymentToSupplier(APIstub, args)
	} else if function == "oemRepaysToBank" { 
		return s.oemRepaysToBank(APIstub, args)	
	} else if function == "getUser" {
		return s.getUser(APIstub, args)
	} else if function == "createInvoiceWithJsonInput" {
		return s.createInvoiceWithJsonInput(APIstub, args)
	} 

	return shim.Error("Invalid Smart Contract function name.")
}


func getQueryResultForQueryString(APIstub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {


	resultsIterator, err := APIstub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
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
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")	

	return buffer.Bytes(), nil
}


func (s *SmartContract) queryInvoice(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	invoiceAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(invoiceAsBytes)
}


func (s *SmartContract) raiseInvoice(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	var invoice = Invoice{InvoiceNumber: args[1], BilledTo: args[2], InvoiceDate: args[3], InvoiceAmount: args[4], ItemDescription: args[5], GR: args[6], IsPaid: args[6], PaidAmount: args[7], Repaid: args[8], RepaymentAmount: args[9]}

	invoiceAsBytes, _ := json.Marshal(invoice)
	APIstub.PutState(args[0], invoiceAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) createInvoiceWithJsonInput(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}
	fmt.Println("args[1] > ", args[1])
	invoiceAsBytes := []byte(args[1])
	invoice := Invoice{}
	err := json.Unmarshal(invoiceAsBytes, &invoice)

	if err != nil {
		return shim.Error("Error During Invoice Unmarshall")
	}
	APIstub.PutState(args[0], invoiceAsBytes)
	return shim.Success(nil)
}

func (s *SmartContract) queryAllInvoices(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "INVOICE0"
	endKey := "INVOICE999"

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

	fmt.Printf("- queryAllInvoices:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) goodsReceive(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	invoiceAsBytes, _ := APIstub.GetState(args[0])
	invoice := Invoice{}

	json.Unmarshal(invoiceAsBytes, &invoice)
	invoice.GR = args[1]

	invoiceAsBytes, _ = json.Marshal(invoice)
	APIstub.PutState(args[0], carAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) bankPaymentToSupplier(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	invoiceAsBytes, _ := APIstub.GetState(args[0])
	invoice := Invoice{}

	json.Unmarshal(invoiceAsBytes, &invoice)
	invoice.isPaid = args[1]
	invoice.paidAmount = args[2]

	invoiceAsBytes, _ = json.Marshal(invoice)
	APIstub.PutState(args[0], invoiceAsBytes)

	return shim.Success(nil)

	return shim.Success(buffer.Bytes())
}
}

func (s *SmartContract) getUser(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

/*	attr := args[0]
	attrValue, _, _ := cid.GetAttributeValue(APIstub,attr)

	msp, _ := cid.GetMSPID(APIstub)

	var buffer bytes.Buffer
		buffer.WriteString("{\"User\":")
		buffer.WriteString("\"")
		buffer.WriteString(attrValue)
		buffer.WriteString("\"")

		buffer.WriteString(", \"MSP\":")
		buffer.WriteString("\"")

		buffer.WriteString(msp+"_DUMMY_change")
		buffer.WriteString("\"")

		buffer.WriteString("}")
	

	return shim.Success(buffer.Bytes())
*/

return shim.Success(nil)

}

func (s *SmartContract) oemRepaysToBank(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	invoiceAsBytes, _ := APIstub.GetState(args[0])
	invoice := Invoice{}

	json.Unmarshal(invoiceAsBytes, &invoice)
	invoice.Repaid = args[1]
	invoice.RepaymentAmount = args[2]

	invoiceAsBytes, _ = json.Marshal(invoice)
	APIstub.PutState(args[0], invoiceAsBytes)

	return shim.Success(nil)

	return shim.Success(buffer.Bytes())
}


// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
