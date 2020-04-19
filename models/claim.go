package models

import (
	jwt "github.com/dgrijalva/jwt-go"
)

//save what will be in me payload and allows create a token
type Claim struct {
	User `json:"user"`
	//standar claim
	jwt.StandardClaims
}
