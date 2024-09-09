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

type SlitherContract struct {
	ID                 int               `json:"id"`
	Name               string            `json:"name"`
	Code               string            `json:"code"`
	IsAbstract         bool              `json:"is_abstract"`
	IsInterface        bool              `json:"in_interface"`
	IsLibrary          bool              `json:"is_library"`
	Functions          []Function        `json:"functions"`
	InheritedContracts []SlitherContract `json:"inherited_contracts"`
}

type SlitherOutput struct {
	Contracts []SlitherContract `json:"contracts"`
}

type Contract struct {
	ID                 int        `json:"id"`
	Name               string     `json:"name"`
	IsAbstract         bool       `json:"is_abstract"`
	IsInterface        bool       `json:"in_interface"`
	IsLibrary          bool       `json:"is_library"`
	Functions          []Function `json:"functions"`
	InheritedContracts []Contract `json:"inherited_contracts"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	Type    string `json:"type"`
	Model   string `json:"model"`
	Hidden  bool   `json:"display"`
}
