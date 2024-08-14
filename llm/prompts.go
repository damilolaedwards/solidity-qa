package llm

import (
	"assistant/types"
	"fmt"
)

func InitialPrompts(targetContracts string) []types.Message {
	return []types.Message{{
		Role:    "user",
		Content: fmt.Sprintf("The code in triple quotes is my solidity codebase: '''%s'''", targetContracts),
		Type:    "text",
		Model:   DefaultModel,
	}, {
		Role:    "assistant",
		Content: "Okay!!! What do you need help with?",
		Type:    "text",
		Model:   DefaultModel,
	}}
}
