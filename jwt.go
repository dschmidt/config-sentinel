package main

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ohler55/ojg/jp"
)

func extractClaimsWithoutValidation(tokenStr string) (jwt.MapClaims, error) {
	parser := jwt.NewParser(jwt.WithoutClaimsValidation()) // skip expiration/audience checks

	token, _, err := parser.ParseUnverified(tokenStr, jwt.MapClaims{})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}
	return nil, fmt.Errorf("unexpected claims type")
}

func hasClaim(claims jwt.MapClaims, jsonPath string) (bool, error) {
	claimPath, err := jp.ParseString(jsonPath)
	if err != nil {
		return false, fmt.Errorf("Error parsing claim path %s: %w", jsonPath, err)
	}

	return len(claimPath.Get(claims)) > 0, nil
}
