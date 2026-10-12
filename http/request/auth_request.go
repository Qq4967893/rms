package request

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	CodeId   string `json:"code_id" validate:"required"`
	Code     string `json:"code" validate:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required,refresh_token_len"`
}
