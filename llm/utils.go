package llm

import (
	"fmt"
	"github.com/pkoukk/tiktoken-go"
)

func CalculateNumTokens(messages []ApiMessage) (int, error) {
	var numTokens = 0

	// TODO: Calculate token length for specified model
	tkm, err := tiktoken.EncodingForModel("gpt-4o")
	if err != nil {
		return numTokens, fmt.Errorf("unable to get model encoding: %v", err)
	}

	for _, message := range messages {
		token := tkm.Encode(message.Content, nil, nil)
		numTokens += len(token)
	}

	return numTokens, nil
}
