package frontend

type ExchangePublicTokenInput struct {
	PublicToken     string `json:"public_token" binding:"required"`
	InstitutionName string `json:"institution_name" binding:"required"`
}
