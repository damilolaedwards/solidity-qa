package llm

import "fmt"

func InitialPrompt(targetContracts string) Message {
	return Message{
		Role:    "user",
		Content: fmt.Sprintf("The code in triple quotes is my solidity codebase: '''%s'''", targetContracts),
	}
}
