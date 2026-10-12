package response

type StatusCode int

const (
	StatusOK StatusCode = iota
	StatusError
	StatusInvalidParams
	StatusUnauthorized
)
