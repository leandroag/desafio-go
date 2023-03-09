package entities

// Login representa as credenciais de login de uma conta bancária
type Login struct {
    CPF    string `json:"cpf"`
    Secret string `json:"secret"`
}