package opa

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
)

func GetMatchedAssets(ctx context.Context, regoPolicy string, dataset []map[string]interface{}) (err error) {
	var wg sync.WaitGroup
	compiler, err := ast.CompileModules(map[string]string{
		"match.rego": regoPolicy,
	})
	if err != nil {
		log.Fatal(err)
		return fmt.Errorf("failed to compile rego policy: %w", err)
	}

	for _, asset := range dataset {
		wg.Add(1)
		fmt.Printf("asset: %s", asset)
		fmt.Println()
		go func(asset map[string]interface{}) {
			defer wg.Done()

			regoCalc := rego.New(
				rego.Query("data.match"),
				rego.Compiler(compiler),
				rego.Input(asset),
			)

			// resultset is empty - something wrong here

			resultSet, err := regoCalc.Eval(ctx)
			fmt.Printf("resultset: %s", resultSet)
			fmt.Println()
			if err != nil || resultSet == nil || len(resultSet) == 0 {
				fmt.Println("calling wg.Done()")
				wg.Done()
			}

			for _, result := range resultSet {
				fmt.Printf("result: %v", result)
				fmt.Println()
				for _, expression := range result.Expressions {
					expressionBytes, err := json.Marshal(expression.Value)
					if err != nil {
						wg.Done()
					}

					var expressionMap map[string]interface{}
					err = json.Unmarshal(expressionBytes, &expressionMap)
					if err != nil {
						wg.Done()
					}

					if matched, ok := expressionMap["match"]; ok && matched.(bool) {
						fmt.Printf("Asset matched policy: %s", result)
					}
				}
			}
		}(asset)
	}

	wg.Wait()

	return nil
}
