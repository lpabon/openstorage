/*
Package auth is the gRPC implementation of the SDK gRPC server
Copyright 2018 Portworx

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package auth

import (
	"fmt"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

var (
	requiredClaims = []string{"sub", "iat", "exp"}
	optionalClaims = []string{"name", "email"}
)

type JwtAuthConfig struct {
	SharedSecret []byte
	RsaPublicPem []byte
	ECDSPublicPem []byte
}

type JwtAuthenticator struct {
	config *JwtAuthConfig
	rsaKey interface{}
	ecdsKey interface{}
	sharedSecretKey interface{}
}

func New(config *JwtAuthConfig) (*JwtAuthenticator, error) {

	authenticator := &JwtAuthenticator{
		config: *config,
	}

	if len(config.Secret) != 0 {
		authenticator.sharedSecretKey = config.SharedSecret
	}
	if len(config.RsaPublicPem) != 0 {
		authenticator.rsaKey, err := jwt.ParseRSAPublicKeyFromPEM(config.RsaPublicPem)
		if err != nil {
			return nil, fmt.Errorf("Unable to parse rsa public key: %v", err)
		}
	}
	if len(config.ECDSPublicPem) !0 {
		authenticator.ecdsKey, err := jwt.ParseECPublicKeyFromPEM(config.ECDSPublicPem)
		if err != nil {
			return nil, fmt.Errorf("Unable to parse ecds public key: %v", err)
		}
	}

	return authenticator, nil
}

func (j *JwtAuthenticator) AuthenticateToken(rawtoken string) (*Token, error) {

	// Parse token
	token, err := jwt.Parse(rawtoken, func(token *jwt.Token) (interface{}, error) {

		// Verify Method
		var key interface{}
		if strings.HasPrefix(token.Method.Alg(), "RS") {
			// RS256, RS384, or RS512
			return j.rsaKey, nil
		} else if strings.HasPrefix(token.Method.Alg(), "ES") {
			// ES256, ES384, or ES512
			return j.ecdsKey, nil
		} else if strings.HasPrefix(token.Method.Alg(), "HS") {
			// HS256, HS384, or HS512
			return j.sharedSecretKey, nil
		} else {
			return nil, fmt.Errorf( "Unknown token algorithm: %s", token.Method.Alg())
		}
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("Token failed validation")
	}

	// Get claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if claims == nil || !ok {
		return nil, fmt.Errorf("No claims found in token")
	}

	// Check for required claims
	for _, requiredClaim := range requiredClaims {
		if _, ok := claims[requiredClaim]; !ok {
			// Claim missing
			return nil, fmt.Errorf("Required claim %v missing from token", requiredClaim)
		}
	}

	// Create token information
	tokenInfo := &Token{
		Role: claims["sub"],
	}
	tokenInfo.Email, _ = claims["email"].(string)
	tokenInfo.User, _ = claims["name"].(string)

	return tokenInfo, nil
}
