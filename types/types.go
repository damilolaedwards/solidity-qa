package types

type Parameter struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	IsConstant bool   `json:"is_constant"`
	IsStorage  bool   `json:"is_storage"`
}

type Function struct {
	ID         int         `json:"id"`
	Name       string      `json:"name"`
	Visibility string      `json:"visibility"`
	View       bool        `json:"view"`
	Pure       bool        `json:"pure"`
	Returns    []string    `json:"returns"`
	Modifiers  []string    `json:"modifiers"`
	Parameters []Parameter `json:"parameters"`
}

type Contract struct {
	ID                 int        `json:"id"`
	Name               string     `json:"name"`
	Functions          []Function `json:"functions"`
	InheritedContracts []Contract `json:"inherited_contracts"`
}

type SlitherOutput struct {
	Contracts []Contract `json:"contracts"`
}
