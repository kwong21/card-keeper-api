package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
)

// Jwks struct from Auth0
type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

// JSONWebKeys struct from Auth0
type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

func getJWTMiddleware() *jwtmiddleware.JWTMiddleware {

	return jwtmiddleware.New(
		jwtmiddleware.Options{
			ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
				aud := os.Getenv("AUTH0_AUDIENCE")
				checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)

				if !checkAud {
					return token, errors.New("Invalid Audience")
				}

				iss := os.Getenv("AUTH0_ISSUER")
				checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)

				if !checkIss {
					return token, errors.New("invalid Issuer")
				}

				cert, err := getPemCert(token)

				if err != nil {
					panic(err.Error())
				}

				result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))

				return result, nil
			},
			SigningMethod: jwt.SigningMethodRS256,
		})
}

func getPemCert(token *jwt.Token) (string, error) {
	cert := ""
	resp, err := http.Get("https://dev-shibatek.us.auth0.com/.well-known/jwks.json")

	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	for k := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("unable to find appropriate key")
		return cert, err
	}

	return cert, nil
}
