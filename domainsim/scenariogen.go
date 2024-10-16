package domainsim

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	//"log"
	"os"

	"gopkg.in/yaml.v3"
)

// Test struct that maps to each item in the YAML ProcessTest list
type Test struct {
	TestID string `yaml:"testid"`
	Key    string `yaml:"key"`
	Input  string `yaml:"input"`
	Expect string `yaml:"expect"`
	Error  string `yaml:"error"`
}

// ParsedTestResult holds the parsed JSON data
type ParsedTestResult struct {
	TestID string
	Key    string
	//Input  map[string]interface{}
	//Expect map[string]interface{}
	Input  any 
	Expect any 
	Error  string
}

// ScenarioGen struct
type ScenarioGen struct {
	Parsemap map[string]func(string,string)(any,any,error)
	Checkmap map[string]func(any,any,error,error)(bool)
}

// New function to initialize ScenarioGen
func NewScenarioGen() *ScenarioGen {
	sc := ScenarioGen{
		Parsemap: make(map[string]func(string,string)(any,any,error)),
		Checkmap: make(map[string]func(any,any,error,error)(bool)),
	}
	sc.Parsemap["testhandle1"] = sc.KeyHandle1
	sc.Checkmap["testhandle1"] = sc.CheckHandle1
	return &sc
}
func (sg *ScenarioGen) TestCheck(key string,resval any,expectval any,reserr error,expecterr error)(bool) {

	if f,ok :=sg.Checkmap[key]; ok {
		valid := f(resval,expectval,reserr,expecterr)
		return valid
	}
	return false
}

func (sg *ScenarioGen) CheckHandle1(resval any,expectval any,reserr error,expecterr error)(bool) {
	res := resval.(Entity3)
	exv := expectval.(Entity3)
	
	if res != exv {
		fmt.Println("resval:",res)
		fmt.Println("expect:",exv)
		fmt.Println("return value is invalid")
		return false
	}
	if reserr != expecterr {
		fmt.Println("error value is invalid")
		return false
	}
	return true
}

func (sg *ScenarioGen) KeyHandle1(ipath string,epath string)(any,any,error) {
	var idata Entity1
	var edata Entity3

	ib,err := sg.readJSONFileByte(ipath)
	if err != nil {
		return nil,nil,err	
	} 
	err = json.Unmarshal(ib, &idata)
	if err != nil {
		return nil,nil,err	
	} 
	eb,err := sg.readJSONFileByte(epath)
	if err != nil {
		return nil,nil,err	
	} 
	err = json.Unmarshal(eb, &edata)
	if err != nil {
		return nil,nil,err	
	} 
	return idata,edata,nil
}


// ParseYAMLAndJSON is a method of ScenarioGen
func (sg *ScenarioGen) ParseYAMLAndJSON(yamlPath string) ([]ParsedTestResult, error) {
	// Read YAML file
	yamlFile, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML file: %w", err)
	}

	// Unmarshal YAML into struct
	var processTest struct {
		Tests []Test `yaml:"ProcessTest"`
	}
	err = yaml.Unmarshal(yamlFile, &processTest)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	var results []ParsedTestResult

	// Loop through each test case
	for _, test := range processTest.Tests {
		// Read and unmarshal Input JSON 
		/*
		inputJSON, err := sg.readJSONFile(test.Input)
		if err != nil {
			return nil, fmt.Errorf("failed to read input JSON: %w", err)
		}

		// Read and unmarshal Expect JSON
		expectJSON, err := sg.readJSONFile(test.Expect)
		if err != nil {
			return nil, fmt.Errorf("failed to read expect JSON: %w", err)
		}
		*/

		if f,ok :=sg.Parsemap[test.Key]; ok {
			inputdata,expectdata,err := f(test.Input,test.Expect)
			if err != nil {
				fmt.Println("Unmarshal error:",err)
			} else {
				// Append parsed result
				results = append(results, ParsedTestResult{
					TestID: test.TestID,
					Key:    test.Key,
					Input:  inputdata,
					Expect: expectdata,
					Error:  test.Error,
				})
			}
		} else {
			fmt.Println("Notfound key")
		}
	}
	return results, nil
}

func (sg *ScenarioGen) readJSONFileByte(path string) ([]byte,error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return bytes,nil
}
// readJSONFile is a method of ScenarioGen
func (sg *ScenarioGen) readJSONFile(path string) (map[string]interface{}, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var data map[string]interface{}
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return data, nil
}

