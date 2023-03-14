package opa

import (
	"context"
	"encoding/json"
	"fmt"
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
		return fmt.Errorf("failed to complie rego policy: %w", err)
	}

	for _, asset := range dataset {
		wg.Add(1)

		go func(assetMap map[string]interface{}) {
			defer wg.Done()

			regoCalc := rego.New(
				rego.Query("data.match"),
				rego.Compiler(compiler),
				rego.Input(assetMap),
			)
			resultSet, err := regoCalc.Eval(ctx)
			if err != nil || resultSet == nil || len(resultSet) == 0 {
				wg.Done()
			}

			for _, result := range resultSet {
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
