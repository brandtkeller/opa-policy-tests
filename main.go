package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/brandtkeller/opa-policy-test/pkg/k8s"
	"github.com/brandtkeller/opa-policy-test/pkg/opa"
	"gopkg.in/yaml.v2"
	"k8s.io/client-go/dynamic"
	ctrl "sigs.k8s.io/controller-runtime"
)

type customRegoPolicy struct {
	Targets []regoTargets `yaml:"targets"`
}

type regoTargets struct {
	Target string `yaml:"target"`
	Rego   string `yaml:"rego"`
}

func main() {

	ctx := context.Background()
	config := ctrl.GetConfigOrDie()
	dynamic := dynamic.NewForConfigOrDie(config)

	namespace := "default"
	fmt.Println("Calling GetResourcesDynamically")
	items, err := kube.GetResourcesDynamically(dynamic, ctx,
		"", "v1", "pods", namespace)
	if err != nil {
		fmt.Println(err)
	}

	// Maybe silly? marshall to json and unmarshall to []map[string]interface{}
	jsonData, err := json.Marshal(items)
	if err != nil {
		fmt.Println(err)
	}
	var data []map[string]interface{}
	err = json.Unmarshal(jsonData, &data)

	if err != nil {
		fmt.Println(err)
	}

	// Read in policy
	var poldata customRegoPolicy
	yamlFile, err := ioutil.ReadFile("./policy/latest_tag_pod.yaml")
	if err != nil {
		fmt.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &poldata)
	if err != nil {
		fmt.Printf("Unmarshal: %v", err)
	}

	//fmt.Printf("policy data: %s", string(poldata.Targets[0].Rego))
	// Call GetMatchedAssets()
	opa.GetMatchedAssets(ctx, string(poldata.Targets[0].Rego), data)
}
