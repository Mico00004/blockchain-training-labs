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
	 BilledTo  		string `json:"billedTo"`
	 InvoiceDate 	string `json:"invoiceDate"`
	 InvoiceAmount  	float64 `json:"invoiceAmount"`
	 ItemDescription string `json:"itemDescription"`
	 GR 				string `json:"gr"`
	 IsPaid 			string `json:"isPaid"`
	 PaidAmount 		float64 `json:"paidAmount"`
	 Repaid 			string `json:"repaid"`
	 RepaymentAmount float64 `json:"repaymentAmount"`
 }
 
 func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	 return shim.Success(nil)
 }
 
 
 func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
 
	 
	 function, args := APIstub.GetFunctionAndParameters()
	 
	 if function == "initLedger" {
		 return s.initLedger(APIstub)
	 } else if function == "raiseInvoice" {
		 return s.raiseInvoice(APIstub, args)
	 } else if function == "createInvoiceWithJsonInput" {
		 return s.createInvoiceWithJsonInput(APIstub, args)
	 } else if function == "queryAllInvoices" {
		 return s.queryAllInvoices(APIstub)
	 } else if function == "goodsReceive" {
		 return s.goodsReceive(APIstub, args)	
	 } else if function == "bankPaymentToSupplier" {
		 return s.bankPaymentToSupplier(APIstub, args)
	 } else if function == "getUsers" {
		 return s.getUsers(APIstub, args)
	 } else if function == "oemRepaysToBank" { 
		 return s.oemRepaysToBank(APIstub, args)	
	 }
 
	 return shim.Error("Invalid Smart Contract function name.")
 }
 
 func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	 invoices := []Invoice{
		 Invoice{
			 InvoiceNumber: "INVOICE0", 
			 BilledTo: "Asus Co.", 
			 InvoiceDate: "02/07/2019", 
			 InvoiceAmount: 1000000, 
			 ItemDescription: "Laptops", 
			 GR: "No", IsPaid: "No", 
			 PaidAmount: 0, 
			 Repaid: "No", 
			 RepaymentAmount: 0},
	 }
 
	 i := 0
	 for i < len(invoices) {
		 fmt.Println("i is ", i)
		 invoiceAsBytes, _ := json.Marshal(invoices[i])
		 APIstub.PutState("INVOICE"+strconv.Itoa(i), invoiceAsBytes)
		 fmt.Println("Added", invoices[i])
		 i = i + 1
	 }
 
	 return shim.Success(nil)
 }
 func (s *SmartContract) raiseInvoice(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	 if len(args) != 10 {
		 return shim.Error("Incorrect number of arguments. Expecting 5")
	 
	 }

	 InvoiceAmountFloat, _ := strconv.ParseFloat(args[3],64)
	 PaidAmountFloat,_:= strconv.ParseFloat(args[7],64)
	 RepaymentAmountFloat, _ := strconv.ParseFloat(args[9],64)

	 var invoice = Invoice{
		 InvoiceNumber: args[0], 
		 BilledTo: args[1], 
		 InvoiceDate: args[2], 
		 InvoiceAmount: InvoiceAmountFloat, 
		 ItemDescription: args[4], 
		 GR: args[5], 
		 IsPaid: args[6], 
		 PaidAmount: PaidAmountFloat, 
		 Repaid: args[8], 
		 RepaymentAmount: RepaymentAmountFloat}

	 invoiceAsBytes, _ := json.Marshal(invoice)
	 APIstub.PutState(args[0], invoiceAsBytes)
 
	 return shim.Success(nil)
 }
 
 func (s *SmartContract) createInvoiceWithJsonInput(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	 if len(args) != 1 {
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
 
	 
	 var buffer bytes.Buffer
	 buffer.WriteString("[")
 
	 bArrayMemberAlreadyWritten := false
	 for resultsIterator.HasNext() {
		 queryResponse, err := resultsIterator.Next()
		 if err != nil {
			 return shim.Error(err.Error())
		 }
		 
		 if bArrayMemberAlreadyWritten == true {
			 buffer.WriteString(",")
		 }
	
		 // Record is a JSON object, so we write as-is
		 buffer.WriteString(string(queryResponse.Value))
		 buffer.WriteString("}\n	")
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
	 APIstub.PutState(args[0], invoiceAsBytes)
 
	 return shim.Success(nil)
 }
 
 func (s *SmartContract) bankPaymentToSupplier(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	 if len(args) != 3 {
		 return shim.Error("Incorrect number of arguments. Expecting 2")
	 }
 
	 invoiceAsBytes, _ := APIstub.GetState(args[0])
	 invoice := Invoice{}


	 paidAmountFloat, _ := strconv.ParseFloat(args[2],64)
	 json.Unmarshal(invoiceAsBytes, &invoice)
	 if invoice.InvoiceAmount > paidAmountFloat {
		invoice.PaidAmount = paidAmountFloat
	 	invoice.IsPaid = args[1]
	 }
	 
	 invoiceAsBytes, _ = json.Marshal(invoice)
	 APIstub.PutState(args[0], invoiceAsBytes)
 
	 return shim.Success(nil)
 }
 
 
 func (s *SmartContract) getUsers(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
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
 
	 RepaymentAmountFloat, _ := strconv.ParseFloat(args[2],64)
	 json.Unmarshal(invoiceAsBytes, &invoice)
	 if RepaymentAmountFloat > invoice.PaidAmount{
	 invoice.Repaid = args[1]
	 invoice.RepaymentAmount = RepaymentAmountFloat
	 }
 
	 invoiceAsBytes, _ = json.Marshal(invoice)
	 APIstub.PutState(args[0], invoiceAsBytes)
 
	 return shim.Success(nil)
 }
 
 
 // The main function is only relevant in unit test mode. Only included here for completeness.
 func main() {
 
	 // Create a new Smart Contract
	 err := shim.Start(new(SmartContract))
	 if err != nil {
		 fmt.Printf("Error creating new Smart Contract: %s", err)
	 }
 }
 
 