package types

type Credentials struct {
	PlanetId string `json:"planetId" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CredentialsWithCode struct {
	Credentials
	Code string `json:"code" binding:"required"`
}
