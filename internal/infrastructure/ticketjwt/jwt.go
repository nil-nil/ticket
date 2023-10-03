package ticketjwt

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/nil-nil/ticket/internal/domain"
)

var (
	ErrGettingToken   = errors.New("unable to get token from token string")
	ErrTokenInvalid   = errors.New("token failed validity check")
	ErrGettingClaims  = errors.New("unable to get claims from token")
	ErrGettingSubject = errors.New("token claims does not have a subject")
	ErrInvalidSubject = errors.New("token subject is not valid")
	ErrGettingUser    = errors.New("error getting user for subject")
)

type GetUserFunc func(ctx context.Context, userID uint64) (user domain.User, err error)

type jwtAuthProvider struct {
	getUserFunc   GetUserFunc
	publicKey     interface{}
	privateKey    interface{}
	signingMethod Protocol
	tokenLifetime uint64
}

func (p jwtAuthProvider) GetUser(ctx context.Context, tokenString string) (user domain.User, err error) {
	token, err := p.getToken(tokenString)
	if err != nil {
		return domain.User{}, errors.Join(ErrGettingToken, err)
	}
	if token != nil && !token.Valid {
		return domain.User{}, ErrTokenInvalid
	}

	var userID uint64

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return domain.User{}, ErrGettingClaims
	}

	// We're using the "sub" claim for the user ID
	sub, ok := claims["sub"]
	if !ok {
		return domain.User{}, ErrGettingSubject
	}

	// The JSON decoder trats the number as a float64
	floatSub, ok := sub.(float64)
	if !ok {
		return domain.User{}, ErrInvalidSubject
	}
	userID = uint64(floatSub)

	u, err := p.getUserFunc(ctx, userID)
	if err != nil {
		return domain.User{}, errors.Join(ErrGettingSubject, err)
	}

	return u, nil
}

func (p jwtAuthProvider) NewToken(user domain.User) (string, error) {
	var method jwt.SigningMethod
	switch p.signingMethod {
	case RS512:
		method = jwt.SigningMethodRS512
	}

	if user.ID == 0 {
		return "", fmt.Errorf("invalid jwt subject for user %+v", user)
	}

	token := jwt.NewWithClaims(method, jwt.MapClaims{
		"sub": user.ID,
		"nbf": time.Now().Unix(),
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Second * time.Duration(p.tokenLifetime)).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString(p.privateKey)
}

func (p jwtAuthProvider) ValidateToken(tokenString string) (err error) {
	token, err := p.getToken(tokenString)
	if err != nil {
		return errors.Join(ErrGettingToken, err)
	}
	if !token.Valid {
		return ErrTokenInvalid
	}
	return nil
}

func (p jwtAuthProvider) getToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// Check the signing method
		if t.Method.Alg() != p.signingMethod.String() {
			return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
		}

		return p.publicKey, nil
	})
}

func NewJwtAuthProvider(
	getUserFunc GetUserFunc,
	publicKeyBytes []byte,
	privateKeyBytes []byte,
	signingMethod Protocol,
	tokenLifetime uint64,
) (jwtAuthProvider, error) {
	var (
		publicKey  interface{}
		privateKey interface{}
		err        error
	)

	switch signingMethod {
	case RS512:
		publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
		if err != nil {
			return jwtAuthProvider{}, err
		}
		privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
		if err != nil {
			return jwtAuthProvider{}, err
		}
	default:
		return jwtAuthProvider{}, fmt.Errorf("invalid jwt signingMethod %s", signingMethod)
	}

	return jwtAuthProvider{
		getUserFunc:   getUserFunc,
		publicKey:     publicKey,
		privateKey:    privateKey,
		signingMethod: signingMethod,
		tokenLifetime: tokenLifetime,
	}, nil
}

type Protocol int

const (
	InvalidProtocol Protocol = iota
	RS512
	// ECDSA
)

func (p Protocol) String() string {
	switch p {
	case RS512:
		return "RS512"
		// case ECDSA:
		// 	return "ECDSA"
	}
	return "unknown"
}

func GetJWTProtocol(s string) Protocol {
	switch s {
	case RS512.String():
		return RS512
	default:
		return InvalidProtocol
	}
}
