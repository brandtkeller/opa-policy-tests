# opa-policy-tests
Repository for testing some json parsing and OPA policy logic

## Testing
`go run main.go`

### Output
2x matches (expected result)

## Logic credit

https://dzone.com/articles/building-with-open-policy-agent-opa-for-better-pol
https://github.com/open-policy-agent/opa/blob/main/rego/example_test.go
https://itnext.io/generically-working-with-kubernetes-resources-in-go-53bce678f887

## Current Thoughts

- What does rego/OPA offer us as a solution?
    - Rego
    > Rego queries, put simply, ask questions of data stored in OPA. Rego policies then evaluate whether data complies with, or violates, the expected state of a system 
    - OPA
    > OPA is essentially a flexible engine that not only enables you to enumerate and write policies but also has powerful search capabilities that don’t require the learning of any new custom syntax (as with other databases, for example) that can be applied to any JSON data set.

- What does a match indicate?
    - Maybe this doesn't  matter as much as users can write rego independent of this auditor
- How will I retrieve data?
    - Domains
        - Kubernetes
        - AWS
    - Syntax
        - Custom?
        - How do I select which resources to inspect or exclude?
