
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

alias ci-e2e := ci-e2e-test
# wait about 10 minutes for kurtosis syncing
[group("ci")]
ci-e2e-test rpcPort *args="":
    @RPC_PORT={{rpcPort}} go test ./e2e/ {{args}}

[group("ci")]
lint:
    @golangci-lint run -c .golangci.yaml

[group("ci")]
fmt:
    @go fmt ./...

alias k-up := kurtosis-up
# local kurtosis network up for e2e testing, or test actual behavior
[group("develop")]
kurtosis-up:
    @bash ./_fixture/kurtosis/scripts/run.sh

alias k-port := kurtosis-port
# detect port of local kurtosis network
[group("develop")]
kurtosis-port:
    @bash ./_fixture/kurtosis/scripts/detect-rpc-port.sh

alias k-wait := kurtosis-wait
# wait for kurtosis syncing
[group("develop")]
kurtosis-wait:
    @timeout 600 bash -c 'until kurtosis service logs gas-testnet cl-1-lighthouse-geth --match Synced 2>/dev/null | grep -q "Synced"; do echo "Still syncing..."; sleep 5; done' || (echo "Sync timeout after 10 minutes" && exit 1)
    @echo "Synchronization completed!"
