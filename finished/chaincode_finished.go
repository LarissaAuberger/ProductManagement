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
	"errors"
	"fmt"
	"encoding/json"
  "io/ioutil"
  "net/http"
	 "encoding/base64"
   "bytes"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)


// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Manufacturer struct {
	CompanyName string `json:"CompanyName"`
}

type Product struct {
	PID								string `json:"pid"`
	ProductionDate 		int `json:"productionDate"`
	Manufacturer 			string `json:"manufacturer"`
	PlantCode 				string `json:"plantCode"`
	Shipments         []Shipment `json:"shimpents"`
}


type Shipment struct {
	Id								string `json:"id"`
	Origin						string `json:"origin"`
	Destination				string `json:"destination"`
	Carrier						string `json:"carrier"`
	DepartureDate			int `json:"departureDate"`
	ArrivalDate				int `json:"arrivalDate"`
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	err := stub.PutState("changing state in init ", []byte(args[0]))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args)
	} else if function == "register_product" {
		return t.register_product(stub, args)
	} else if function == "add_shipment" {
		return t.add_shipment(stub, args)
	}

	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "queryAsConsumer" { //read a variable
		return t.queryAsConsumer(stub, args)
	}

	if function == "queryAsManufacturer" { //read a variable
		return t.queryAsManufacturer(stub, args)
	}

	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}


// write - invoke function to add a production to the blockchain
func (t *SimpleChaincode) register_product(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	////// args
	// [0]		[1]
	// pid		jsonString of product
	var pid string
	var err error

	fmt.Println("running register_product()")

	var product Product
	json.Unmarshal([]byte(args[0]), &product)

	pid = product.PID
	details := Product {
	  PID: product.PID,
		ProductionDate: product.ProductionDate,
		Manufacturer: product.Manufacturer,
		PlantCode: product.PlantCode,
	}

	detailsAsJsonBytes, _:= json.Marshal (details)

	err = stub.PutState(pid, []byte(detailsAsJsonBytes)) //write the product into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil

}


func (t *SimpleChaincode) add_shipment (stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	fmt.Println("running add_shipment()")

  var shipment Shipment
	json.Unmarshal([]byte(args[1]), &shipment)

	details := Shipment {
		Id: shipment.Id,
		Origin: shipment.Origin,
		Destination: shipment.Destination,
		Carrier: shipment.Carrier,
		DepartureDate: shipment.DepartureDate,
		ArrivalDate: shipment.ArrivalDate,
	}

	var p Product
	var product_id = args[0]
	pAsBytes, err := stub.GetState(product_id)
	json.Unmarshal(pAsBytes, &p)

  p.Shipments = append(p.Shipments, details)

	pAsBytes, err = json.Marshal(p)
	err = stub.PutState(product_id, pAsBytes)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// write - invoke function to write key/value pair
func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[0] 
	value = args[1]
	err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (t *SimpleChaincode) queryAsManufacturer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

var key, location, jsonResp string
	var err
	
if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}
	
key = args[0]
	valAsbytes, err := stub.GetState(key)
	return valAsbytes, nil
	
}


// read - query function to read key/value pair
func (t *SimpleChaincode) queryAsConsumer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, location, jsonResp string
	var err error
  if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

  location = args [1]
	key = args[0]
	valAsbytes, err := stub.GetState(key)
	var product Product
	json.Unmarshal(valAsbytes, &product)


	fmt.Printf("location: %s", location)

	if location != product.Shipments[0].Destination {
		// WIoTP REST API --> event für Device "BCFakeDetector" eventtype "fake-alert" JSON {"PID":"<replace-me>","fake":"true"}
		url := "http://20wql7.messaging.internetofthings.ibmcloud.com:1883/api/v0002/application/types/FakeDetector/devices/BCFakeDetector/events/fake-alert"
    //https://orgId.messaging.internetofthings.ibmcloud.com:8883/api/v0002/application/types/typeId/devices/deviceId/events/eventId
    //fmt.Println("URL:>", url)

    var jsonStr = []byte("{ \"PID\":\"" + product.PID + "\",\"fake\":\"true\"}")

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    //req.Header.Set("X-Custom-Header", "myvalue")
    req.Header.Set("Content-Type", "application/json")
		var user string = "a-20wql7-b28fat8pmw"
		var password string = "T)DwTzn+plN*9tL38N"
		req.Header.Add("Authorization","Basic "+basicAuth(user, password))

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    //fmt.Println("response Status:", resp.Status)
    //fmt.Println("response Headers:", resp.Header)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("response Body:", string(body))
	} else  {

		}

	//	shipmentAsJsonBytes, _ := json.Marshal (details)

	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}



func basicAuth(username, password string) string {
  auth := username + ":" + password
   return base64.StdEncoding.EncodeToString([]byte(auth))
}

// read - query function to read key/value pair
func (t *SimpleChaincode) read_number(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
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
