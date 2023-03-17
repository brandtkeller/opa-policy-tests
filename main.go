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
	Domain    string   `yaml:"domain"`
	ApiGroup  string   `yaml:"apiGroups"`
	Kinds     []string `yaml:"kinds"`
	Namespace string   `yaml:"namespaces"`
	Rego      string   `yaml:"rego"`
}

func main() {

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

	fmt.Printf("policy data: %s\n", poldata)

	ctx := context.Background()
	config := ctrl.GetConfigOrDie()
	dynamic := dynamic.NewForConfigOrDie(config)

	for _, target := range poldata.Targets {

		// for each target
		for _, kind := range target.Kinds {

			// TODO: split apigroup into group and version
			// TODO: create a slice of items per target to run validation against

			fmt.Println("Calling GetResourcesDynamically")
			// TODO: make items be a slice we can append to for later use
			items, err := kube.GetResourcesDynamically(dynamic, ctx,
				"", "v1", kind, target.Namespace)
			if err != nil {
				fmt.Println(err)
			}
		}
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

	// Call GetMatchedAssets()
	opa.GetMatchedAssets(ctx, string(poldata.Targets[0].Rego), data)
}
