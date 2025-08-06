package main

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ohler55/ojg/jp"
)

func applyRules(data map[string]any, rules []rule, claims jwt.MapClaims) error {
	for _, rule := range rules {
		switch rule.Rule {
		case "included-for-claim":
			val, err := hasClaim(claims, rule.Value)
			if err != nil {
				return err
			}
			if val {
				continue
			}

			// remove if claim is *not* present
			path, err := jp.ParseString(rule.Path)
			if err != nil {
				return err
			}
			path.RemoveOne(data)

		case "excluded-for-claim":
			val, err := hasClaim(claims, rule.Value)
			if err != nil {
				return err
			}
			if !val {
				continue
			}

			// remove if claim *is* present
			path, err := jp.ParseString(rule.Path)
			if err != nil {
				return err
			}
			path.RemoveOne(data)

		default:
			return fmt.Errorf("Unknown rule: %s\n", rule.Rule)
		}

	}

	return nil
}
