package main 

import (
	"fmt"
	"domainsim"
)

func main() {
	ec := domainsim.NewEntityControl()
	sg := domainsim.NewScenarioGen() 

	trs,err := sg.ParseYAMLAndJSON("scenario.yaml")
	if err != nil {
		fmt.Println("err",err)
	}
	for _,tv := range trs {
		ex,err := ec.ProcessData(tv.Key,tv.Input)
		var checkerr error = nil
		if tv.Error != "none" {
			checkerr = fmt.Errorf(tv.Error)
		}
		valid := sg.TestCheck(tv.Key,ex,tv.Expect,err,checkerr)
		if valid {
			fmt.Printf("%s is passed.\n",tv.TestID)
		} else {
			fmt.Printf("%s is failed.\n",tv.TestID)
		}
	}
}
