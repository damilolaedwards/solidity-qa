package llm

import (
	"assistant/types"
	"fmt"
	"strings"
)

// InitialPrompts generates the initial messages for a conversation with an LLM.
// It creates two messages:
// 1. A user message that provides the Solidity codebase.
// 2. An assistant message that acknowledges the codebase and asks for further instructions (this is required because Claude requires messages to come in "user", "assistant" pairs).
//
// Parameters:
//   - targetContracts: A string containing the Solidity codebase.
//
// Returns:
//   - A slice of types.Message representing the initial conversation.
func InitialPrompts(targetContracts string) []types.Message {
	return []types.Message{{
		Role:    "user",
		Content: fmt.Sprintf("The code in triple quotes is my solidity codebase: '''%s'''", targetContracts),
		Type:    "text",
		Model:   Models[DefaultModelIdentifier].Name,
	}, {
		Role:    "assistant",
		Content: "Okay!!! What do you need help with?",
		Type:    "text",
		Model:   Models[DefaultModelIdentifier].Name,
	}}
}

// AskAboutFunctionPrompt creates a prompt string for asking about a specific function in a contract.
// The prompt instructs the LLM to provide a detailed explanation of the given function.
//
// Parameters:
//   - functionName: A string containing the name of the function to be explained.
//   - contractName: A string containing the name of the contract where the function is located.
//
// Returns:
//   - A string containing the complete prompt for asking about the function.
func AskAboutFunctionPrompt(functionName string, contractName string) string {
	var prompt strings.Builder

	prompt.WriteString(fmt.Sprintf("Give me a detailed explanation of function **%s** in contract **%s**", functionName, contractName))

	return prompt.String()
}

// GenerateReportPrompt creates a prompt string for generating an audit report.
// The prompt instructs the LLM to create a report focused on specific types of bugs
// and provides a sample format for the report.
//
// Parameters:
//   - reportType: A string specifying the type of report to be generated.
//   - reportSample: A string containing a sample of how the report should look.
//   - additionalMessage: An optional string with additional instructions for report creation.
//
// Returns:
//   - A string containing the complete prompt for generating the audit report.
func GenerateReportPrompt(reportType string, reportSample string, additionalMessage string) string {
	var prompt strings.Builder

	prompt.WriteString(fmt.Sprintf("Generate an audit report that contains details on %s related bugs in the provided codebase.\n", reportType))
	prompt.WriteString(fmt.Sprintf("Here's some samples in triple quotes on how the report should look: '''%s'''", reportSample))

	if additionalMessage != "" {
		prompt.WriteString(fmt.Sprintf("Take the instructions in triple quotes into consideration in creating the report: '''%s'''", additionalMessage))
	}

	return prompt.String()
}
