package models

type TokenBody struct {
	TokenString string `json:"refresh_token"`
}

type RefreshToken struct {
	UUID      string
	GUID      string
	TokenHash []byte
}
