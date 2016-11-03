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
  // "io/ioutil"
  // "net/http"
	// "encoding/base64"
  //  "bytes"

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
}

type ProductDetails struct {
	ProductionDate 		int `json:"productionDate"`
	Manufacturer 			string `json:"manufacturer"`
	PlantCode 				string `json:"plantCode"`
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
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
func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args)
	} else if function == "add_product" {
		return t.add_product(stub, args)
	}

	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	}

	if function == "read_number" { //read a variable
		return t.read_number(stub, args)
	}

	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

// write - invoke function to add a production to the blockchain
func (t *SimpleChaincode) add_product(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	////// args
	// [0]		[1]
	// pid		jsonString of product
	var pid string
	var err error

	fmt.Println("running add_product()")

	var product Product
	json.Unmarshal([]byte(args[0]), &product)

	pid = product.PID
	details := ProductDetails {
		ProductionDate: product.ProductionDate,
		Manufacturer: product.Manufacturer,
		PlantCode: product.PlantCode,
	}

	detailsAsJasonBytes, _:= json.Marshal (details)

	err = stub.PutState(pid, []byte(detailsAsJasonBytes)) //write the product into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil

}

// write - invoke function to write key/value pair
func (t *SimpleChaincode) write(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[0] //rename for funsies
	value = args[1]
	err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}


// read - query function to read key/value pair
func (t *SimpleChaincode) read(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error
  if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		// // WIoTP REST API --> event für Device "BCFakeDetector" eventtype "fake-alert" JSON {"PID":"<replace-me>","fake":"true"}
		// url := "http://20wql7.messaging.internetofthings.ibmcloud.com:1883/api/v0002/application/types/FakeDetector/devices/BCFakeDetector/events/fake-alert"
    // //https://orgId.messaging.internetofthings.ibmcloud.com:8883/api/v0002/application/types/typeId/devices/deviceId/events/eventId
    // //fmt.Println("URL:>", url)
    // var jsonStr = []byte("{ \"PID\":\"<replace-me>\",\"fake\":\"true\"}")
    // req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    // //req.Header.Set("X-Custom-Header", "myvalue")
    // req.Header.Set("Content-Type", "application/json")
		// var user string = "a-20wql7-b28fat8pmw"
		// var password string = "T)DwTzn+plN*9tL38N"
		// req.Header.Add("Authorization","Basic "+basicAuth(user, password))
		//
    // client := &http.Client{}
    // resp, err := client.Do(req)
    // if err != nil {
    //     panic(err)
    // }
    // defer resp.Body.Close()
    // //fmt.Println("response Status:", resp.Status)
    // //fmt.Println("response Headers:", resp.Header)
    // body, _ := ioutil.ReadAll(resp.Body)
    // fmt.Println("response Body:", string(body))


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
func (t *SimpleChaincode) read_number(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
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
