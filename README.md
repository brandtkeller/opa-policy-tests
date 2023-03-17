# opa-policy-tests
Repository for testing some json parsing and rego policy logic

## Hypothesis
Through the use of OPA/Rego - would could open up the possibility of validating any payload that we can receive or transform to JSON. Initial test cases would be kubernetes API adapter and AWS adapter - but this doesn't mean a generic REST/http adapter would be off the table either. This means that others could potentially setup transformers of data for validation purposes as long as it can produce a JSON payload. IE a daemonset with access to the OS that exposes an API for validating OS configurations. 

## Match
What does a match mean in this context?

If we're writing our policy to meet parity with what we have today - we need the ability to:
- Look at all pods
- process exclusions
- process rego to define that all pods must have 1 container that includes `istio/proxy`

## Parity with Lula Kyverno engine
- Ability to target a cluster - Done
- Define target for execution - Apiversion - kind - namespace
- Define exclusions - namespaces or namespaces and resources
- Write rego for processing container in pod - use rego playground
- Ability to target static manifest

## Testing
- Create/connect to a cluster
- Deploy the demo workloads
    - `kubectl apply -f demo/`
- Run the main.go
    - `go run main.go`
- Should result in 2 policy matches and 1 policy no match

## Logic credit

https://dzone.com/articles/building-with-open-policy-agent-opa-for-better-pol
https://github.com/open-policy-agent/opa/blob/main/rego/example_test.go
https://itnext.io/generically-working-with-kubernetes-resources-in-go-53bce678f887
https://snyk.io/blog/opa-rego-usage-for-policy-as-code/
https://www.openpolicyagent.org/docs/latest/policy-reference/#strings
https://play.openpolicyagent.org/

## Current Thoughts

- What does rego/OPA offer us as a solution?
    - Rego
    > Rego queries, put simply, ask questions of data stored in OPA. Rego policies then evaluate whether data complies with, or violates, the expected state of a system 
    - OPA
    > OPA is essentially a flexible engine that not only enables you to enumerate and write policies but also has powerful search capabilities that don’t require the learning of any new custom syntax (as with other databases, for example) that can be applied to any JSON data set.
    - Together: A domain agnostic engine for processing JSON data and matching against 

- What does a match indicate?
    - Maybe this doesn't  matter as much as users can write rego independent of this auditor
- How will I retrieve data?
    - Domains
        - Kubernetes
        - AWS
    - Syntax
        - Custom?
        - How do I select which resources to inspect or exclude?


## Execution

- Lock down namespaces and apigroups to a single entry for now
- Per target (loop)
    - per kind (loop)
        - Call getresourcesdynamically()
        - append each iteration to slice of items