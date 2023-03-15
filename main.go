package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/brandtkeller/opa-policy-test/pkg/opa"
	"gopkg.in/yaml.v2"
)

type customRegoPolicy struct {
	Targets []regoTargets `yaml:"targets"`
}

type regoTargets struct {
	Target string `yaml:"target"`
	Rego   string `yaml:"rego"`
}

func main() {

	// Read in dataset
	dataset, _ := ioutil.ReadFile("./dataset/latest_tag.json")
	var data []map[string]interface{}
	err := json.Unmarshal(dataset, &data)

	if err != nil {
		fmt.Println("error occurred")
	}

	//fmt.Printf("dataset: %s", data)

	// Read in policy
	var poldata customRegoPolicy
	yamlFile, err := ioutil.ReadFile("./policy/latest_tag.yaml")
	if err != nil {
		fmt.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &poldata)
	if err != nil {
		fmt.Printf("Unmarshal: %v", err)
	}

	//fmt.Printf("policy data: %s", string(poldata.Targets[0].Rego))
	// Call GetMatchedAssets()
	ctx := context.TODO()
	opa.GetMatchedAssets(ctx, string(poldata.Targets[0].Rego), data)
}
