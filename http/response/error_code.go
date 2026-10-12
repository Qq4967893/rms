package response

type ErrorCode int

const (
	CodeSuccess ErrorCode = 102000 // 业务成功
	CodeFailure ErrorCode = 102001 // 业务失败

	CodeBadRequest ErrorCode = 104000

	CodeUnauthorized        ErrorCode = 104010
	CodeAccessTokenExpired  ErrorCode = 104011
	CodeRefreshTokenExpired ErrorCode = 104012

	CodeForbidden ErrorCode = 104030

	CodeNotFound ErrorCode = 104040

	CodeServerError ErrorCode = 105000
)
