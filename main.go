package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	kube "github.com/brandtkeller/opa-policy-test/pkg/k8s"
	"github.com/brandtkeller/opa-policy-test/pkg/opa"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
	ctrl "sigs.k8s.io/controller-runtime"
)

type customRegoPolicy struct {
	Targets []regoTargets `yaml:"targets"`
}

type regoTargets struct {
	Domain    string   `yaml:"domain"`
	ApiGroup  string   `yaml:"apiGroup"`
	Kinds     []string `yaml:"kinds"`
	Namespace string   `yaml:"namespace"`
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

	// fmt.Printf("policy data: %s\n", poldata)

	ctx := context.Background()
	config := ctrl.GetConfigOrDie()
	dynamic := dynamic.NewForConfigOrDie(config)

	var resources []unstructured.Unstructured
	// for each target
	for _, target := range poldata.Targets {

		// TODO - process exclusions?

		// for each target kind
		for _, kind := range target.Kinds {

			// TODO: split apigroup into group and version
			fmt.Println(kind)
			// check for group/version combo - there is only ever one `/` right?
			var group, version string
			if strings.Contains(target.ApiGroup, "/") {
				split := strings.Split(target.ApiGroup, "/")
				group = split[0]
				version = split[1]
			} else {
				group = ""
				version = target.ApiGroup
			}

			// TODO - Better way to get proper lowercase + plural of a resource
			// Pod and Pods and pod are not acceptable inputs
			resource := strings.ToLower(kind) + "s"

			items, err := kube.GetResourcesDynamically(dynamic, ctx,
				group, version, resource, target.Namespace)
			if err != nil {
				fmt.Println(err)
			}
			resources = append(resources, items...)
		}
		// Maybe silly? marshall to json and unmarshall to []map[string]interface{}
		jsonData, err := json.Marshal(resources)
		if err != nil {
			fmt.Println(err)
		}
		var data []map[string]interface{}
		err = json.Unmarshal(jsonData, &data)

		if err != nil {
			fmt.Println(err)
		}

		// Call GetMatchedAssets()
		results, err := opa.GetMatchedAssets(ctx, string(poldata.Targets[0].Rego), data)
		if err != nil {
			fmt.Println(err)
		}

		// Now let's do something with this
		fmt.Println(results.Match)
	}

}
