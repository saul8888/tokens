package authentication

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"

	"github.com/tokens/models"
)

//se firman los token con una  llave privada
//create -> openssl genrsa -out private.rsa 1024
//se verifican con una llave publica
//create -> openssl rsa -in private.rsa -pubout > public.rsa.pub

//create variables, type punters for keys
var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
)

func init() {
	//the read archive in format bytes for save  the keys private anad public
	privateBytes, err := ioutil.ReadFile("./private.rsa")
	if err != nil {
		log.Fatal("private key was not read")
	}

	publicBytes, err := ioutil.ReadFile("./public.rsa.pub")
	if err != nil {
		log.Fatal("public key was not read")
	}

	//for load in the form of a key private and public
	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateBytes)
	if err != nil {
		log.Fatal("could not do the parse of private")
	}
	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicBytes)
	if err != nil {
		log.Fatal("could not do the parse of public")
	}
}

func GenerateJWT(user models.User) string {
	//create a struct of my Claim
	claims := models.Claim{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			Issuer:    "token test", //object of token
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	//encode to base64
	result, err := token.SignedString(privateKey)
	if err != nil {
		log.Fatal("could not sign private token")
	}
	return result

}

func Login(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		//fmt.Println("error reading user %s", err)
		fmt.Println("error reading user")
		return
	}
	if user.Name == "prueba" && user.Password == "prueba" {
		user.Password = ""
		user.Role = "admin"
		token := GenerateJWT(user)
		result := models.Responsetoken{Token: token}
		jsonResult, err := json.Marshal(result)
		if err != nil {
			fmt.Println("error generating the json")
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		//responseJSON(w, newUser)
		w.Write(jsonResult)
	} else {
		w.WriteHeader(http.StatusForbidden)
		fmt.Println(w, "usser or password invalid")
	}
}

//func hola(token *jwt.Token) (interface{}, error) {
//	return publicKey, nil
//}

func ValidateToken(w http.ResponseWriter, r *http.Request) {
	//token, err := request.ParseFromRequest(r, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
	//token, err := jwt.ParseFromRequest(r, func(token *jwt.Token) (interface{}, error) {
	token, err := request.ParseFromRequestWithClaims(r, request.OAuth2Extractor, &models.Claim{}, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	//token, err := request.ParseFromRequestWithClaims(r, request.OAuth2Extractor, &models.Claim{}, hola)
	if err != nil {
		switch err.(type) {
		case *jwt.ValidationError:
			vErr := err.(*jwt.ValidationError)
			switch vErr.Errors {
			case jwt.ValidationErrorExpired:
				fmt.Fprintln(w, "your token expired")
				return
			case jwt.ValidationErrorSignatureInvalid:
				fmt.Fprintln(w, "the signature does not match")
				return
			default:
				fmt.Fprintln(w, "the signature does not match")
				return
			}
		default:
			fmt.Fprintln(w, "your token is not valid")
			return
		}
	}
	if token.Valid {
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprintln(w, "welcome to the system")
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "your token is not valid")
	}
}
