//
// Copyright (c) 2015 The heketi Authors
// Copyright (c) 2018 Portworx
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package auth

import (
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
)

var (
	requiredClaims = []string{"sub", "iat", "exp"}
	optionalClaims = []string{"name", "email"}
)

type JwtAuth struct {
	adminKey []byte
	userKey  []byte
}

type JwtAuthConfig struct {
	AdminKey []byte
	UserKey  []byte
}

func NewJwtAuth(config *JwtAuthConfig) *JwtAuth {

	j := &JwtAuth{}
	j.adminKey = config.AdminKey
	j.userKey = config.UserKey

	return j
}

func (j *JwtAuth) Type() string {
	return "jwt"
}

func (j *JwtAuth) AuthenticateToken(rawtoken string) (*Token, error) {

	// Parse token
	var claims jwt.MapClaims
	tokenInfo := &Token{
		Role: RoleUnknown,
	}
	token, err := jwt.Parse(rawtoken, func(token *jwt.Token) (interface{}, error) {

		// Verify Method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		var ok bool
		claims, ok = token.Claims.(jwt.MapClaims)
		if claims == nil || !ok {
			return nil, fmt.Errorf("No claims found in token")
		}

		// Get claims
		if sub, ok := claims["sub"]; ok {
			switch sub {
			case string(RoleAdministrator):
				tokenInfo.Role = RoleAdministrator
				return j.adminKey, nil
			case string(RoleUser):
				tokenInfo.Role = RoleUser
				return j.userKey, nil
			default:
				return nil, fmt.Errorf("Unknown role: %s", sub)
			}
		}

		return nil, fmt.Errorf("Token missing iss claim")
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("Invalid token")
	}

	// Check for required claims
	for _, requiredClaim := range requiredClaims {
		if _, ok := claims[requiredClaim]; !ok {
			// Claim missing
			return nil, fmt.Errorf("Required claim %v missing from token", requiredClaim)
		}
	}

	tokenInfo.Email, _ = claims["email"].(string)
	tokenInfo.User, _ = claims["name"].(string)

	return tokenInfo, nil
}
