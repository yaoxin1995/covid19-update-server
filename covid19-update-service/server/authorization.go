package server

import (
	"covid19-update-service/model"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
)

type jwkCollection struct {
	Keys []jsonWebKeys `json:"keys"`
}

type jsonWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

// Handles authorization via Auth0
type authorizationHandler struct {
	JWKS       jwkCollection
	ISS        string
	AUD        string
	Middleware *jwtmiddleware.JWTMiddleware
}

const tokenContext = "tokenContext"

// Iss is the issuer address (Auth0 authorization server URI) and aud is the API identifier used in Auth0.
func newAuthenticationHandler(iss, aud string) (authorizationHandler, error) {
	jwks, err := getJwks(iss)
	if err != nil {
		return authorizationHandler{}, err
	}
	handler := authorizationHandler{
		JWKS: jwks,
		ISS:  iss,
		AUD:  aud,
	}
	handler.createJWTMiddleWare()
	return handler, nil
}

func getJwks(iss string) (jwkCollection, error) {
	var jwks = jwkCollection{}
	resp, err := http.Get(fmt.Sprintf("%s.well-known/jwks.json", iss))

	if err != nil {
		return jwks, fmt.Errorf("could not load jwkCollection: %v", err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return jwkCollection{}, fmt.Errorf("could not decode jwkCollection: %v", err)
	}
	return jwks, nil
}

func (a *authorizationHandler) getPemCert(token *jwt.Token) (string, error) {
	cert := ""
	for k := range a.JWKS.Keys {
		if token.Header["kid"] == a.JWKS.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + a.JWKS.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("unable to find appropriate key")
		return cert, err
	}

	return cert, nil
}

func (a *authorizationHandler) createJWTMiddleWare() {
	a.Middleware = jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			// Verify 'aud' claim
			checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(a.AUD, false)
			if !checkAud {
				return token, errors.New("invalid audience")
			}
			// Verify 'iss' claim
			checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(a.ISS, false)
			if !checkIss {
				return token, errors.New("invalid issuer")
			}

			cert, err := a.getPemCert(token)
			if err != nil {
				return nil, fmt.Errorf("could not get certificate: %v", err)
			}

			result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
			return result, nil
		},
		SigningMethod: jwt.SigningMethodRS256,
		UserProperty:  tokenContext,
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err string) {
			w.Header().Set("WWW-Authenticate", fmt.Sprintf("Bearer realm=\"%soauth/token\"", a.ISS))
			writeHTTPResponse(model.NewError(fmt.Sprintf("could not perform authentcation: %v", err)), http.StatusUnauthorized, w, r)
		},
		EnableAuthOnOptions: false,
	})
}

func (a *authorizationHandler) getSubject(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		cert, err := a.getPemCert(token)
		if err != nil {
			return nil, err
		}
		result, err := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
		if err != nil {
			return nil, err
		}
		return result, nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		return "", errors.New("could not get claims")
	}

	return claims.Subject, nil
}
