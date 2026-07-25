package ether_test

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fullNodePackages are packages that only geth's full node pulls in. They ride
// along with ethclient/simulated and roughly quadruple an SDK user's binary
// (~12MB -> ~49MB), so no package but ether/simulated may reach them.
// See #494.
var fullNodePackages = []string{
	"github.com/ethereum/go-ethereum/ethclient/simulated",
	"github.com/ethereum/go-ethereum/core/vm",
	"github.com/ethereum/go-ethereum/eth",
	"github.com/ethereum/go-ethereum/node",
	"github.com/cockroachdb/pebble",
	"github.com/urfave/cli/v2",
}

// simulatedOnlyPackages are the packages allowed to link the full node:
// ether/simulated is where simulated support deliberately lives, and playground
// is a local sandbox binary that drives a simulated chain.
var simulatedOnlyPackages = map[string]bool{
	"github.com/poteto-go/go-alchemy-sdk/ether/simulated": true,
	"github.com/poteto-go/go-alchemy-sdk/playground":      true,
}

// TestCoreImportPath_DoesNotLinkFullNode walks every package of the module -
// so packages added later are guarded without touching this test - and asserts
// the full node closure stays confined to simulatedOnlyPackages.
func TestCoreImportPath_DoesNotLinkFullNode(t *testing.T) {
	// Act: one `go list` reports the transitive deps of every package at once.
	// The pattern is absolute: `go test` runs this binary inside ether/.
	out, err := exec.Command(
		"go", "list", "-f", `{{.ImportPath}} {{join .Deps " "}}`,
		"github.com/poteto-go/go-alchemy-sdk/...",
	).Output()
	require.NoError(t, err)

	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		pkg, depList, _ := strings.Cut(line, " ")
		if simulatedOnlyPackages[pkg] {
			continue
		}

		t.Run(pkg, func(t *testing.T) {
			deps := make(map[string]bool)
			for _, dep := range strings.Fields(depList) {
				deps[dep] = true
			}

			// Assert
			for _, forbidden := range fullNodePackages {
				assert.Falsef(
					t,
					deps[forbidden],
					"%s must not depend on %s: it links a full geth node into every SDK user (#494)",
					pkg, forbidden,
				)
			}
		})
	}
}
