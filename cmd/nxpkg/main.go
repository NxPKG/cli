package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/nxpkg/cli/internal/cmd"
	"github.com/nxpkg/cli/internal/nxpkgstate"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("go-nxpkg is expected to be invoked via nxpkg")
		os.Exit(1)
	}

	executionStateString := os.Args[1]
	var executionState nxpkgstate.ExecutionState
	decoder := json.NewDecoder(strings.NewReader(executionStateString))
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&executionState)
	if err != nil {
		fmt.Printf("Error unmarshalling execution state: %v\n Execution state string: %v\n", err, executionStateString)
		os.Exit(1)
	}

	exitCode := cmd.RunWithExecutionState(&executionState, nxpkgVersion)
	os.Exit(exitCode)
}
