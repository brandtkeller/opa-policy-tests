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
	Exclude   []string `yaml:"exclude"`
	Rego      string   `yaml:"rego"`
}

func main() {
	ctx := context.Background()
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

	// for each target
	for _, target := range poldata.Targets {
		var resources []unstructured.Unstructured
		// TODO - Per target - process domain and execute query accordingly
		switch domain := target.Domain; domain {
		case "kubernetes":
			resources, err = queryKube(ctx, target)
			if err != nil {
				fmt.Println(err)
			}
		default:
			fmt.Printf("No domain connector available for %s", domain)
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

		var includedData []map[string]interface{}
		for _, value := range data {
			resourceNamespace := value["metadata"].(map[string]interface{})["namespace"]

			exclude := false
			for _, exns := range target.Exclude {
				if exns == resourceNamespace {
					exclude = true
				}
			}
			if !exclude {
				includedData = append(includedData, value)
			}
		}

		// Call GetMatchedAssets()
		results, err := opa.GetMatchedAssets(ctx, string(poldata.Targets[0].Rego), includedData)
		if err != nil {
			fmt.Println(err)
		}

		// Now let's do something with this
		fmt.Println(results.Match)
	}

}

func queryKube(ctx context.Context, target regoTargets) (resources []unstructured.Unstructured, err error) {

	config := ctrl.GetConfigOrDie()
	dynamic := dynamic.NewForConfigOrDie(config)

	for _, kind := range target.Kinds {
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
	return resources, nil
}
