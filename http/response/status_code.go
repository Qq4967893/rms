package response

type StatusCode int

const (
	StatusOK           StatusCode = 200
	StatusBadRequest   StatusCode = 400
	StatusUnauthorized StatusCode = 401
	StatusForbidden    StatusCode = 403
	StatusNotFound     StatusCode = 404
	StatusServerError  StatusCode = 500
)
