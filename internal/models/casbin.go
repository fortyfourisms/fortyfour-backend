package models

type CasbinPolicy struct {
	Role     string `json:"role"`
	Resource string `json:"resource"`
	Action   string `json:"action"`
}

type CasbinPolicyResponse struct {
	Policies []CasbinPolicy `json:"policies"`
	Count    int            `json:"count"`
}
