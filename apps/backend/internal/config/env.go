package config

import "strings"

// NormalizeEnv maps common environment aliases to canonical values.
// Canonical values are: dev, staging, production.
func NormalizeEnv(env string) string {
	e := strings.ToLower(strings.TrimSpace(env))
	switch e {
	case "":
		return "dev"
	case "prod":
		return "production"
	case "stage":
		return "staging"
	default:
		return e
	}
}

// IsProduction returns true when running in production.
func IsProduction(env string) bool {
	return NormalizeEnv(env) == "production"
}

// IsProductionLike returns true for environments that should behave like production
// from a security/config perspective (e.g. require secrets).
func IsProductionLike(env string) bool {
	switch NormalizeEnv(env) {
	case "production", "staging":
		return true
	default:
		return false
	}
}
