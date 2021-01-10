package middleware

import (
	"encoding/json"
	"errors"
	"net/http"

	configs "card-keeper-api/internal/configs"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
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

// GetJWTMiddleware returns the configured middleware to handle JWT auth
func GetJWTMiddleware(config configs.AuthConfiguration) *jwtmiddleware.JWTMiddleware {
	return jwtmiddleware.New(
		jwtmiddleware.Options{
			ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
				aud := config.Audience
				checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)

				if !checkAud {
					return token, errors.New("Invalid Audience")
				}

				iss := config.Issuer
				checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)

				if !checkIss {
					return token, errors.New("invalid Issuer")
				}

				cert, err := getPemCert(token, config.JWKS)

				if err != nil {
					panic(err)
				}

				result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))

				return result, nil
			},
			SigningMethod: jwt.SigningMethodRS256,
		})
}

// CorsMiddleware sets the required CORs header
func CorsMiddleware(allowOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

func getPemCert(token *jwt.Token, jwksURI string) (string, error) {
	cert := ""
	resp, err := http.Get(jwksURI)

	if err != nil {
		return cert, errors.New(resp.Status)
	}

	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, errors.New(resp.Status)
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
