package main

import (
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
)

type Address struct {
	Address_Line string `json:"address_line"`
	City string `json:"city"`
}

type Customer struct {
	Name string `json:"name"`
	Gender string `json:"gender"`
	DOB string `json:"dob"`
	Aadhar string `json:"aadhar_no"`
	Address Address `json:"address"`
	PAN string `json:"pan_no"`
	Cibil_Score int32 `json:"cibil_score"`
	Marital_Status string `json:"marital_status"`
	Education map[string]string `json:"education"`
	Employement map[string]string `json:"employement"`
	Health map[string]string `json:"health"`
	Possesions map[string]string `json:"possesions"`
}

type SimpleChainCode struct {
}

func main()  {
	err := shim.Start(new(SimpleChainCode))
	if err != nil {
		fmt.Printf("Error starting simple chaincode due to :%s", err)
	}
}

func (t *SimpleChainCode) Init(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Print("Initializing ChainCode.....")
	err := t.createTable(stub)
	if err != nil {
		return nil, errors.New("Error creating customer table. Cause : "+ err)
	}
	return []byte("INITIALIZATION O.K"), nil
}

func (t *SimpleChainCode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	switch function {
	case "insert":
		return t.insertDataIntoTable(stub, args)
	}

	return nil, errors.New("Unknown Function Invocation")
}

func (t *SimpleChainCode) createTable(stub shim.ChaincodeStubInterface) error {
	var columns []*shim.ColumnDefinition
	columns = append(columns, &shim.ColumnDefinition{Name:"name", Type: shim.ColumnDefinition_STRING, Key: true})
	columns = append(columns, &shim.ColumnDefinition{Name:"gender", Type: shim.ColumnDefinition_STRING, Key: false})
	columns = append(columns, &shim.ColumnDefinition{Name:"dob", Type: shim.ColumnDefinition_STRING, Key: false})
	columns = append(columns, &shim.ColumnDefinition{Name:"aadhar_no", Type: shim.ColumnDefinition_STRING, Key: true})
	columns = append(columns, &shim.ColumnDefinition{Name:"pan_no", Type: shim.ColumnDefinition_STRING, Key: true})
	columns = append(columns, &shim.ColumnDefinition{Name:"cibil_score", Type: shim.ColumnDefinition_INT32, Key: false})
	columns = append(columns, &shim.ColumnDefinition{Name:"address_line", Type: shim.ColumnDefinition_STRING, Key: false})
	columns = append(columns, &shim.ColumnDefinition{Name:"city", Type: shim.ColumnDefinition_STRING, Key: true})
	columns = append(columns, &shim.ColumnDefinition{Name:"marital_status", Type: shim.ColumnDefinition_STRING, Key: false})
	columns = append(columns, &shim.ColumnDefinition{Name:"education_map", Type: shim.ColumnDefinition_STRING, Key: false})
	columns = append(columns, &shim.ColumnDefinition{Name:"employement_map", Type: shim.ColumnDefinition_STRING, Key: false})
	columns = append(columns, &shim.ColumnDefinition{Name:"health_map", Type: shim.ColumnDefinition_STRING, Key: false})
	columns = append(columns, &shim.ColumnDefinition{Name:"possesion_map", Type: shim.ColumnDefinition_STRING, Key: false})

	return stub.CreateTable("customerTable", columns)
}

func (t *SimpleChainCode) insertDataIntoTable(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
	if len(args) != 1 {
		return nil, errors.New("Incorrect Number of Arguments. Required 1.")
	}
	var customer Customer
	err := json.Unmarshal(args[0], customer)
	if  err != nil {
		return nil, errors.New("Unable To Parse Argument String. Cause : " + err)
	}

	var columns []*shim.Column
	columns = append(columns, &shim.Column{Value: &shim.Column_String_{String_: customer.Name}})
	columns = append(columns, &shim.Column{Value: &shim.Column_String_{String_: customer.Gender}})
	columns = append(columns, &shim.Column{Value: &shim.Column_String_{String_: customer.DOB}})
	columns = append(columns, &shim.Column{Value: &shim.Column_String_{String_: customer.Aadhar}})
	columns = append(columns, &shim.Column{Value: &shim.Column_String_{String_: customer.PAN}})
	columns = append(columns, &shim.Column{Value: &shim.Column_Int32{Int32 : customer.Cibil_Score}})
	columns = append(columns, &shim.Column{Value: &shim.Column_String_{String_: customer.Address.Address_Line}})
	columns = append(columns, &shim.Column{Value: &shim.Column_String_{String_: customer.Address.City}})
	columns = append(columns, &shim.Column{Value: &shim.Column_String_{String_: customer.Marital_Status}})

	education_field, err := json.Marshal(customer.Education)
	if err != nil {
		columns = append(columns, &shim.Column{Value: &shim.Column_String_{String_: ""}})
	} else {
		columns = append(columns, &shim.Column{Value: &shim.Column_String_{String_: education_field}})
	}

	employement_field, err := json.Marshal(customer.Employement)
	if err != nil {
		columns = append(columns, &shim.Column{Value: &shim.Column_String_{String_: ""}})
	} else {
		columns = append(columns, &shim.Column{Value: &shim.Column_String_{String_: employement_field}})
	}
	health_field, err := json.Marshal(customer.Health)
	if err != nil {
		columns = append(columns, &shim.Column{Value: &shim.Column_String_{String_: ""}})
	} else {
		columns = append(columns, &shim.Column{Value: &shim.Column_String_{String_: health_field}})
	}
	possesions_field, err := json.Marshal(customer.Possesions)
	if err != nil {
		columns = append(columns, &shim.Column{Value: &shim.Column_String_{String_: ""}})
	} else {
		columns = append(columns, &shim.Column{Value: &shim.Column_String_{String_: possesions_field}})
	}
	row := shim.Row{Columns: columns}
	ok, err := stub.InsertRow("customerTable", row)
	if err != nil {
		return  nil, errors.New("Insert Operation Failed. Cause :" + err)
	}
	if !ok {
		return nil, errors.New("Insert Operation Failed. Row with the given key already exists.")
	}

	return []byte("Insert Success"), nil
}