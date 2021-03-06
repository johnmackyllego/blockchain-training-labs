package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	// "github.com/hyperledger/fabric/core/chaincode/lib/cid"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the car structure, with 4 properties.  Structure tags are used by encoding/json library
type Invoice struct {
	InvoiceNumber string `json:"invoiceNumber"`
	BilledTo string `json:"billedTo"`
	InvoiceDate string `json:"invoiceDate"`
	InvoiceAmount string `json:"invoiceAmount"`
	ItemDescription string `json:"itemDescription"`
	GR string `json:"gR"`
	IsPaid string `json:"isPaid"`
	PaidAmount string `json:"paidAmount"`
	Repaid string `json:"repaid"`
	RepaymentAmount string `json:"repaymentAmount"`
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "raiseInvoice" {
		return s.raiseInvoice(APIstub, args)
	} else if function == "goodsReceived" {
		return s.raiseInvoice(APIstub, args)
	} 
	// else if function == "bankPaymentToSupplier" {
	// 	return s.raiseInvoice(APIstub, args)
	// } else if function == "oemRepaysToBank" {
	// 	return s.raiseInvoice(APIstub, args)
	// }
	else if function == "displayAllInvoices" {
		return s.displayAllInvoices(APIstub)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	invoices := []Invoice{
		Invoice{InvoiceNumber: "INVOICE1", BilledTo: "Prius", InvoiceDate: "blue", InvoiceAmount: "Tomoko", ItemDescription: "blue", GR: "true", IsPaid: "true", PaidAmount: "1", Repaid: "true", RepaymentAmount: "1"},
	}

	i := 0
	for i < len(invoices) {
		fmt.Println("i is ", i)
		invoiceAsBytes, _ := json.Marshal(invoices[i])
		APIstub.PutState("CAR"+strconv.Itoa(i), invoiceAsBytes)
		fmt.Println("Added", invoices[i])
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) raiseInvoice(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 11 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	var invoice = Invoice{InvoiceNumber: args[1], BilledTo: args[2], InvoiceDate: args[3], InvoiceAmount: args[4], ItemDescription: args[5], GR: args[6], IsPaid: args[7], PaidAmount: args[8], Repaid: args[9], RepaymentAmount: args[10]}


	invoiceAsBytes, _ := json.Marshal(invoice)
	APIstub.PutState(args[0], invoiceAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) goodsReceived(APIstub shim.ChaincodeStubInterface) sc.Response {

	invoiceAsBytes, _ := APIstub.GetState(args[0])
	invoice := Invoice{}

	json.Unmarshal(invoiceAsBytes, &invoice)
	invoice.Owner = args[1]

	invoiceAsBytes, _ = json.Marshal(invoice)
	APIstub.PutState(args[0], invoiceAsBytes)

	return shim.Success(nil)
}

// func (s *SmartContract) bankPaymentToSupplier(APIstub shim.ChaincodeStubInterface) sc.Response {

// 	invoiceAsBytes, _ := APIstub.GetState(args[0])
// 	invoice := Invoice{}

// 	json.Unmarshal(invoiceAsBytes, &invoice)
// 	invoice.Owner = args[1]

// 	invoiceAsBytes, _ = json.Marshal(invoice)
// 	APIstub.PutState(args[0], invoiceAsBytes)

// 	return shim.Success(nil)
// }

// func (s *SmartContract) oemRepaysToBank(APIstub shim.ChaincodeStubInterface) sc.Response {

// 	invoiceAsBytes, _ := APIstub.GetState(args[0])
// 	invoice := Invoice{}

// 	json.Unmarshal(invoiceAsBytes, &invoice)
// 	invoice.Owner = args[1]

// 	invoiceAsBytes, _ = json.Marshal(invoice)
// 	APIstub.PutState(args[0], invoiceAsBytes)

// 	return shim.Success(nil)
// }

func (s *SmartContract) displayAllInvoices(APIstub shim.ChaincodeStubInterface) sc.Response {

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

	fmt.Printf("- displayAllInvoices:\n%s\n", buffer.String())

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