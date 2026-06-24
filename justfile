
alias ut := unit-test
[group("ci")]
unit-test path="./..." *args="":
    @go test $(go list {{path}} | grep -v /e2e) {{args}}

alias ut-race := unit-test-race
[group("ci")]
unit-test-race path="./..." *args="":
    @go test $(go list {{path}} | grep -v /e2e) {{args}} -race

alias ut-cov := unit-test-coverage
[group("ci")]
unit-test-coverage path="./..." *args="":
    @go test $(go list {{path}} | grep -v /e2e) {{args}} -cover -gcflags=all=-l -coverprofile=coverage.out

alias e2e := e2e-test
[group("ci")]
e2e-test rpcPort="8545" *args="":
    @RPC_PORT={{rpcPort}} go test ./e2e/ {{args}}

[group("ci")]
lint:
    @golangci-lint run -c .golangci.yaml

[group("ci")]
fmt:
    @go fmt ./...

alias a-up := anvil-up
[group("develop")]
anvil-up port="8545":
    anvil --port {{port}}
