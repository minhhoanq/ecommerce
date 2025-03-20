package token

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type JWTMaker struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func LoadPrivateKey(path string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read private key file: %v", err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	privateKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err == nil {
		rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
		if ok {
			return rsaPrivateKey, nil
		}
	}

	rsaPrivateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err == nil {
		return rsaPrivateKey, nil
	}

	return nil, fmt.Errorf("not a valid RSA private key")
}

func LoadPublicKey(path string) (*rsa.PublicKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read public key file: %v", err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	pubKeyGeneric, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err == nil {
		rsaPubKey, ok := pubKeyGeneric.(*rsa.PublicKey)
		if ok {
			return rsaPubKey, nil
		}
	}

	rsaPubKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err == nil {
		return rsaPubKey, nil
	}

	return nil, fmt.Errorf("not a valid RSA public key")
}

func NewJWTMaker(publicKeyPath, privateKeyPath string) (Maker, error) {

	publicKey, err := LoadPublicKey(publicKeyPath)
	if err != nil {
		return nil, err
	}

	privateKey, err := LoadPrivateKey(privateKeyPath)
	if err != nil {
		return nil, err
	}

	return &JWTMaker{
		publicKey:  publicKey,
		privateKey: privateKey,
	}, nil
}

func (maker *JWTMaker) CreateToken(user_id uuid.UUID, role_id int, duration time.Duration) (string, *Payload, error) {
	payload := NewPayload(user_id, role_id, duration)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, payload)
	token, err := jwtToken.SignedString(maker.privateKey)

	return token, payload, err
}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodRSA)
		if !ok {
			return nil, ErrInvalidToken
		}
		return maker.publicKey, nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}

		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
