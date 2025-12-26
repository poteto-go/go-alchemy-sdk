
alias ut := unit-test
[group("ci")]
unit-test path="./..." *args="":
    @go test $(go list {{path}} | grep -v /e2e) {{args}}

alias ut-cov := unit-test-coverage
[group("ci")]
unit-test-coverage path="./..." *args="":
    @go test $(go list {{path}} | grep -v /e2e) {{args}} -cover -gcflags=all=-l -coverprofile=coverage.out

alias ci-e2e := ci-e2e-test
[group("ci")]
ci-e2e-test +rpcPort:
    @RPC_PORT={{rpcPort}} go test ./e2e/

[group("ci")]
lint:
    @golangci-lint run -c .golangci.yaml
