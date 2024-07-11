package llm

import (
	"fmt"
	"strings"
)

var trainingPrompts = TrainingPrompts()

func GenerateInvariants(contractsSourceCode string, unitTestsSourceCode string) (string, error) {
	var invariants strings.Builder
	fmt.Println("Generating invariants...")

	return invariants.String(), nil
}

func ImproveCoverage(contractsSourceCode string, testContractsSourceCode string, unitTestsSourceCode string, coverageReport []byte) (string, error) {
	var improvedInvariants strings.Builder

	return improvedInvariants.String(), nil
}
