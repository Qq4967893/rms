package jwt

import (
	"backend-rms/core/config"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	onpkg "github.com/golang-jwt/jwt/v4"
)

var (
	ErrTokenInvalid           error = errors.New("token invalid")
	ErrTokenExpired           error = errors.New("token has expired")
	ErrTokenMalformed         error = errors.New("token malformed")
	ErrTokenExpiredMaxRefresh error = errors.New("token has passed the maximum refresh time")
	ErrHeaderEmpty            error = errors.New("token is required to access")
	ErrHeaderMalformed        error = errors.New("authorization malformed in request header")
)

type JWT struct {
	SecretKey  []byte
	ExpireTime time.Duration
	MaxRefresh time.Duration
}

type Claims struct {
	LoginManagerID int64 `json:"login_manager_id"`
	onpkg.RegisteredClaims
}

func NewJWT(conf config.JWTConf) *JWT {
	return &JWT{
		SecretKey:  []byte(conf.SecretKey),
		ExpireTime: time.Duration(conf.ExpireTime) * time.Minute,
		MaxRefresh: time.Duration(conf.MaxRefreshTime) * time.Minute,
	}
}

func (jwt *JWT) GetTokenFromHeader(req *http.Request) (string, error) {
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		return "", ErrHeaderEmpty
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", ErrHeaderMalformed
	}

	return parts[1], nil
}

func (jwt *JWT) CreateToken(loginManagerID int64) (tokenString string, expire time.Time, err error) {
	claims := &Claims{
		loginManagerID,
		onpkg.RegisteredClaims{
			ID:        uuid.New().String(),
			IssuedAt:  onpkg.NewNumericDate(time.Now()),
			ExpiresAt: onpkg.NewNumericDate(jwt.expireTime()),
		},
	}

	tokenString, err = jwt.newTokenString(claims)
	expire = claims.ExpiresAt.Time
	return
}

func (jwt *JWT) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.parseTokenString(tokenString)
	if err != nil {
		validationErr, ok := err.(*onpkg.ValidationError)
		if ok {
			if validationErr.Errors == onpkg.ValidationErrorMalformed {
				return nil, ErrTokenMalformed
			} else if validationErr.Errors == onpkg.ValidationErrorExpired {
				return nil, ErrTokenExpired
			}
		}
		return nil, ErrTokenInvalid
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}

func (jwt *JWT) RefreshToken(req *http.Request) (tokenString string, expire time.Time, claims *Claims, err error) {
	var oldToken string
	oldToken, err = jwt.GetTokenFromHeader(req)
	if err != nil {
		return
	}
	var token *onpkg.Token
	token, err = jwt.parseTokenString(oldToken)
	if err != nil {
		validationErr, ok := err.(*onpkg.ValidationError)
		if !ok || validationErr.Errors != onpkg.ValidationErrorExpired {
			return
		}
	}

	claims = token.Claims.(*Claims)
	if claims.IssuedAt.Add(jwt.MaxRefresh).Unix() < time.Now().Unix() {
		err = ErrTokenExpiredMaxRefresh
		return
	}

	claims.ID = uuid.New().String()
	claims.ExpiresAt = onpkg.NewNumericDate(jwt.expireTime())
	tokenString, err = jwt.newTokenString(claims)
	expire = claims.ExpiresAt.Time
	return
}

func (jwt *JWT) newTokenString(claims *Claims) (tokenString string, err error) {
	token := onpkg.NewWithClaims(onpkg.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwt.SecretKey)
	return
}

func (jwt *JWT) parseTokenString(tokenString string) (*onpkg.Token, error) {
	return onpkg.ParseWithClaims(tokenString, &Claims{}, func(token *onpkg.Token) (interface{}, error) {
		return jwt.SecretKey, nil
	})
}

func (jwt *JWT) expireTime() time.Time {
	return time.Now().Add(jwt.ExpireTime)
}
