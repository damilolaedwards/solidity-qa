package llm

import (
	"fmt"
	"strings"
)

func GenerateInvariantsPrompt(numInvariants int, contractSourceCode string, testContractsSourceCode string, unitTestsSourceCode string, coverageReport string) string {
	var prompt strings.Builder

	prompt.WriteString(fmt.Sprintf("Analyze the provided contract source code, test contract source code, unit tests, and coverage report. Your task is to generate %d comprehensive system-level and function-level invariants for the main contracts. These invariants should aim to improve code coverage and identify potential vulnerabilities.\n\n", numInvariants))

	prompt.WriteString("Instructions:\n")
	prompt.WriteString("1. Generate a detailed list of invariants. Be exhaustive in exploring all possible invariants and code execution paths.\n")
	prompt.WriteString("2. Provide code samples for each invariant to aid in writing the tests.\n")
	prompt.WriteString("3. Do not skip any main contract. Cover all main contracts comprehensively.\n")
	prompt.WriteString("4. Ensure the invariants cover all smart contract vulnerability categories: Access Controls, Auditing and Logging, Authentication, Configuration, Cryptography, Data Exposure, Data Validation, Denial of Service, Error Reporting, Patching, Session Management, Testing, Timing, Undefined Behavior.\n")
	prompt.WriteString("5. Be specific and detailed in your response. Each invariant should have a clear and concise description along with the conditions that must hold true.\n\n")

	prompt.WriteString("Contract Source Code:\n")
	prompt.WriteString(fmt.Sprintf("'''%s'''\n\n", contractSourceCode))

	if testContractsSourceCode != "" {
		prompt.WriteString("Test Contracts Source Code:\n")
		prompt.WriteString(fmt.Sprintf("'''%s'''\n\n", testContractsSourceCode))
	}

	if unitTestsSourceCode != "" {
		prompt.WriteString("Unit Tests Source Code:\n")
		prompt.WriteString(fmt.Sprintf("'''%s'''\n\n", unitTestsSourceCode))
	}

	if coverageReport != "" {
		prompt.WriteString("Coverage Report:\n")
		prompt.WriteString(fmt.Sprintf("'''%s'''\n\n", coverageReport))
	}

	return prompt.String()
}

func TrainingPrompts() []Message {
	return []Message{
		{
			Role: "system",
			Content: "Traditional fuzz testing generally explores a binary by providing random inputs to identify new system states or crash the program. However, this does not translate to the smart contract ecosystem since smart contracts cannot 'crash' in the same way. A transaction that reverts is not equivalent to a binary crashing or panicking.\n\n" +
				"With smart contracts, the fuzzing paradigm shifts to validating the **invariants** of the program.\n\n" +
				"**Definition**: An invariant is a property that remains unchanged after one or more operations are applied to it.\n\n" +
				"Invariants are truths about a system. For smart contracts, these can be:\n" +
				"1. **Mathematical invariants**: For example, `a + b = b + a`. The commutative property must hold in any Solidity math library.\n" +
				"2. **ERC20 tokens**: The sum of all user balances should never exceed the total supply of the token.\n" +
				"3. **Automated market makers (e.g., Uniswap)**: `xy = k`, the constant-product formula, which maintains the economic guarantees of AMMs like Uniswap.\n\n" +
				"**Definition**: Smart contract fuzzing uses random sequences of transactions to test the invariants of the smart contract system.\n\n" +
				"Understanding smart contract fuzzing involves identifying, writing, and testing invariants.\n\n" +
				"## Types of Invariants\n\n" +
				"Defining and testing invariants is crucial for assessing the expected system behavior.\n\n" +
				"Invariants are generally divided into two categories: function-level invariants and system-level invariants.\n\n" +
				"### Function-level invariants\n\n" +
				"A function-level invariant arises from the execution of a specific function.\n\n" +
				"Example:\n\n```solidity\nfunction deposit() public payable {\n    // Ensure total deposited amount does not exceed the limit\n    uint256 amount = msg.value;\n    require(totalDeposited + amount <= MAX_DEPOSIT_AMOUNT);\n\n    // Update user balance and total deposited\n    balances[msg.sender] += amount;\n    totalDeposited += amount;\n\n    emit Deposit(msg.sender, amount, totalDeposited);\n}\n```\n\n" +
				"The `deposit` function has these invariants:\n" +
				"1. The ETH balance of `msg.sender` must decrease by `amount`.\n" +
				"2. The ETH balance of `address(this)` must increase by `amount`.\n" +
				"3. `balances[msg.sender]` should increase by `amount`.\n" +
				"4. `totalDeposited` should increase by `amount`.\n\n" +
				"Function-level invariants can be identified by assessing what must be true before and after the execution of a function.\n\n" +
				"### System-level invariants\n\n" +
				"A system-level invariant holds true across the entire execution of a system.\n\n" +
				"Examples:\n" +
				"1. The `xy=k` constant product formula should always hold for Uniswap pools.\n" +
				"2. No user's balance should ever exceed the total supply for an ERC20 token.\n\n" +
				"In the `deposit` function, a system-level invariant is:\n\n**The `totalDeposited` amount should always be less than or equal to the `MAX_DEPOSIT_AMOUNT`.**\n\n" +
				"Since `totalDeposited` can be affected by other functions, it is best tested at the system level.",
		},
		{
			Role: "system",
			Content: "You are a coverage-guided fuzzing assistant. You will be provided with contracts (referred to as main contracts) to examine for possible vulnerabilities/invariants to be tested for.\n" +
				"Use the provided unit tests to better understand the main contracts.\n" +
				"Use the provided test contracts to see existing fuzz tests for the main contracts.\n" +
				"Use the provided coverage report to identify parts of the main contracts not covered by the invariants in the fuzz test contracts.\n" +
				"Generate invariants to improve coverage and explore every execution path.\n" +
				"NOTE: Generate only fuzzing (system-level or function-level) invariants, not unit tests.\n" +
				"NOTE: Your responses should not be in markdown format or be surrounded with triple backticks.\n" +
				"NOTE: Do not add any other text at the beginning or end of your response other than the generated invariants and code samples.",
		},
		{
			Role: "system",
			Content: "The invariants you generate should cover the following smart contract vulnerability categories:\n" +
				"1. Access Controls: Insufficient authorization or assessment of rights\n" +
				"2. Auditing and Logging: Insufficient auditing of actions or logging of problems\n" +
				"3. Authentication: Improper identification of users\n" +
				"4. Configuration: Misconfigured servers, devices, or software components\n" +
				"5. Cryptography: A breach of system confidentiality or integrity\n" +
				"6. Data Exposure: Exposure of sensitive information\n" +
				"7. Data Validation: Improper reliance on the structure or values of data\n" +
				"8. Denial of Service: A system failure with an availability impact\n" +
				"9. Error Reporting: Insecure or insufficient reporting of error conditions\n" +
				"10. Patching: Use of an outdated software package or library\n" +
				"11. Session Management: Improper identification of authenticated users\n" +
				"12. Testing: Insufficient test methodology or test coverage\n" +
				"13. Timing: Race conditions or other order-of-operations flaws\n" +
				"14. Undefined Behavior: Undefined behavior triggered within the system",
		},
	}
}
